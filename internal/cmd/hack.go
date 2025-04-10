package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/cli/flags"
	"github.com/git-town/git-town/v18/internal/cli/print"
	"github.com/git-town/git-town/v18/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v18/internal/cmd/sync"
	"github.com/git-town/git-town/v18/internal/config"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/execute"
	"github.com/git-town/git-town/v18/internal/forge"
	"github.com/git-town/git-town/v18/internal/git"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/gohacks"
	"github.com/git-town/git-town/v18/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v18/internal/messages"
	"github.com/git-town/git-town/v18/internal/undo/undoconfig"
	"github.com/git-town/git-town/v18/internal/validate"
	configInterpreter "github.com/git-town/git-town/v18/internal/vm/interpreter/config"
	fullInterpreter "github.com/git-town/git-town/v18/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v18/internal/vm/program"
	"github.com/git-town/git-town/v18/internal/vm/runstate"
	. "github.com/git-town/git-town/v18/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	hackDesc = "Create a new feature branch off the main branch"
	hackHelp = `
Consider this stack:

main
 \
  branch-1
   \
*   branch-2

We are on the "branch-2" branch. After running "git hack branch-3", our
workspace contains these branches:

main
 \
  branch-1
   \
    branch-2
 \
* branch-3

The new branch "feature-2"
is a child of the main branch.

If there are no uncommitted changes,
it also syncs all affected branches.
`
)

func hackCmd() *cobra.Command {
	addCommitFlag, readCommitFlag := flags.Commit()
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addCommitMessageFlag, readCommitMessageFlag := flags.CommitMessage("the commit message")
	addProposeFlag, readProposeFlag := flags.Propose()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "hack <branch>",
		GroupID: cmdhelpers.GroupIDBasic,
		Args:    cobra.ArbitraryArgs,
		Short:   hackDesc,
		Long:    cmdhelpers.Long(hackDesc, hackHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			commit, err := readCommitFlag(cmd)
			if err != nil {
				return err
			}
			commitMessage, err := readCommitMessageFlag(cmd)
			if err != nil {
				return err
			}
			detached, err := readDetachedFlag(cmd)
			if err != nil {
				return err
			}
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			propose, err := readProposeFlag(cmd)
			if err != nil {
				return err
			}
			prototype, err := readPrototypeFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			if commitMessage.IsSome() || propose.IsTrue() {
				commit = true
			}
			return executeHack(args, commit, commitMessage, detached, dryRun, propose, prototype, verbose)
		},
	}
	addCommitFlag(&cmd)
	addCommitMessageFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addProposeFlag(&cmd)
	addPrototypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeHack(args []string, commit configdomain.Commit, commitMessage Option[gitdomain.CommitMessage], detached configdomain.Detached, dryRun configdomain.DryRun, propose configdomain.Propose, prototype configdomain.Prototype, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineHackData(args, repo, commit, commitMessage, detached, dryRun, propose, prototype, verbose)
	if err != nil || exit {
		return err
	}
	createNewFeatureBranchData, doCreateNewFeatureBranch, convertToFeatureBranchData, doConvertToFeatureBranch := data.Get()
	if doCreateNewFeatureBranch {
		return createFeatureBranch(createFeatureBranchArgs{
			appendData:            createNewFeatureBranchData,
			backend:               repo.Backend,
			beginBranchesSnapshot: createNewFeatureBranchData.branchesSnapshot,
			beginConfigSnapshot:   repo.ConfigSnapshot,
			beginStashSize:        createNewFeatureBranchData.stashSize,
			commandsCounter:       repo.CommandsCounter,
			dryRun:                dryRun,
			finalMessages:         repo.FinalMessages,
			frontend:              repo.Frontend,
			git:                   repo.Git,
			rootDir:               repo.RootDir,
			verbose:               verbose,
		})
	}
	if doConvertToFeatureBranch {
		return convertToFeatureBranch(convertToFeatureBranchArgs{
			beginConfigSnapshot: repo.ConfigSnapshot,
			config:              convertToFeatureBranchData.config,
			makeFeatureData:     convertToFeatureBranchData,
			repo:                repo,
			rootDir:             repo.RootDir,
			verbose:             verbose,
		})
	}
	panic("both config arms were nil")
}

// If set to createNewFeatureData, the user wants to create a new feature branch.
// If set to convertToFeatureData, the user wants to convert an already existing branch into a feature branch.
type hackData = Either[appendFeatureData, convertToFeatureData]

// this configuration is for when "git town hack" is used to make contribution, observed, or parked branches feature branches
type convertToFeatureData struct {
	config         config.ValidatedConfig
	targetBranches configdomain.BranchesAndTypes
}

func createFeatureBranch(args createFeatureBranchArgs) error {
	runProgram := appendProgram(args.appendData, args.finalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: args.beginBranchesSnapshot,
		BeginConfigSnapshot:   args.beginConfigSnapshot,
		BeginStashSize:        args.beginStashSize,
		Command:               "hack",
		DryRun:                args.dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 args.backend,
		CommandsCounter:         args.commandsCounter,
		Config:                  args.appendData.config,
		Connector:               args.appendData.connector,
		DialogTestInputs:        args.appendData.dialogTestInputs,
		FinalMessages:           args.finalMessages,
		Frontend:                args.frontend,
		Git:                     args.git,
		HasOpenChanges:          args.appendData.hasOpenChanges,
		InitialBranch:           args.appendData.initialBranch,
		InitialBranchesSnapshot: args.beginBranchesSnapshot,
		InitialConfigSnapshot:   args.beginConfigSnapshot,
		InitialStashSize:        args.beginStashSize,
		RootDir:                 args.rootDir,
		RunState:                runState,
		Verbose:                 args.verbose,
	})
}

type createFeatureBranchArgs struct {
	appendData            appendFeatureData
	backend               gitdomain.RunnerQuerier
	beginBranchesSnapshot gitdomain.BranchesSnapshot
	beginConfigSnapshot   undoconfig.ConfigSnapshot
	beginStashSize        gitdomain.StashSize
	commandsCounter       Mutable[gohacks.Counter]
	dryRun                configdomain.DryRun
	finalMessages         stringslice.Collector
	frontend              gitdomain.Runner
	git                   git.Commands
	rootDir               gitdomain.RepoRootDir
	verbose               configdomain.Verbose
}

func determineHackData(args []string, repo execute.OpenRepoResult, commit configdomain.Commit, commitMessage Option[gitdomain.CommitMessage], detached configdomain.Detached, dryRun configdomain.DryRun, propose configdomain.Propose, prototype configdomain.Prototype, verbose configdomain.Verbose) (data hackData, exit bool, err error) {
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	targetBranches := gitdomain.NewLocalBranchNames(args...)
	var repoStatus gitdomain.RepoStatus
	repoStatus, err = repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 len(args) == 1 && !repoStatus.OpenChanges,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	localBranchNames := branchesSnapshot.Branches.LocalBranches().Names()
	var branchesToValidate gitdomain.LocalBranchNames
	shouldCreateBranch := len(targetBranches) == 1 && !slices.Contains(localBranchNames, targetBranches[0])
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	if shouldCreateBranch {
		branchesToValidate = gitdomain.LocalBranchNames{}
	} else {
		if len(targetBranches) == 0 {
			branchesToValidate = gitdomain.LocalBranchNames{initialBranch}
		} else {
			branchesToValidate = targetBranches
		}
	}
	connector, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(localBranchNames)
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: branchesToValidate,
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranchNames,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	if !shouldCreateBranch {
		data = Right[appendFeatureData](convertToFeatureData{
			config:         validatedConfig,
			targetBranches: validatedConfig.BranchesAndTypes(branchesToValidate),
		})
		return data, false, nil
	}
	if len(targetBranches) > 1 {
		return data, false, errors.New(messages.HackTooManyArguments)
	}
	targetBranch := targetBranches[0]
	var remotes gitdomain.Remotes
	remotes, err = repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch, repo.UnvalidatedConfig.NormalConfig.DevRemote) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	branchNamesToSync := gitdomain.LocalBranchNames{validatedConfig.ValidatedConfigData.MainBranch}
	if detached {
		branchNamesToSync = validatedConfig.RemovePerennials(branchNamesToSync)
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, branchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, false, err
	}
	data = Left[appendFeatureData, convertToFeatureData](appendFeatureData{
		branchInfos:               branchesSnapshot.Branches,
		branchesSnapshot:          branchesSnapshot,
		branchesToSync:            branchesToSync,
		commit:                    commit,
		commitMessage:             commitMessage,
		config:                    validatedConfig,
		connector:                 connector,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             initialBranch,
		newBranchParentCandidates: gitdomain.LocalBranchNames{validatedConfig.ValidatedConfigData.MainBranch},
		nonExistingBranches:       nonExistingBranches,
		preFetchBranchInfos:       preFetchBranchSnapshot.Branches,
		previousBranch:            previousBranch,
		propose:                   propose,
		prototype:                 prototype,
		remotes:                   remotes,
		stashSize:                 stashSize,
		targetBranch:              targetBranch,
	})
	return data, false, err
}

func convertToFeatureBranch(args convertToFeatureBranchArgs) error {
	for branchName, branchType := range args.makeFeatureData.targetBranches {
		switch branchType {
		case
			configdomain.BranchTypeContributionBranch,
			configdomain.BranchTypeObservedBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypePrototypeBranch:
			if err := args.config.NormalConfig.SetBranchTypeOverride(configdomain.BranchTypeFeatureBranch, branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeFeatureBranch:
			return fmt.Errorf(messages.HackBranchIsAlreadyFeature, branchName)
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.HackCannotFeatureMainBranch)
		case configdomain.BranchTypePerennialBranch:
			return fmt.Errorf(messages.HackCannotFeaturePerennialBranch, branchName)
		}
		fmt.Printf(messages.HackBranchIsNowFeature, branchName)
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:               args.repo.Backend,
		BeginBranchesSnapshot: None[gitdomain.BranchesSnapshot](),
		BeginConfigSnapshot:   args.beginConfigSnapshot,
		Command:               "observe",
		CommandsCounter:       args.repo.CommandsCounter,
		FinalMessages:         args.repo.FinalMessages,
		Git:                   args.repo.Git,
		RootDir:               args.rootDir,
		TouchedBranches:       args.makeFeatureData.targetBranches.Keys().BranchNames(),
		Verbose:               args.verbose,
	})
}

type convertToFeatureBranchArgs struct {
	beginConfigSnapshot undoconfig.ConfigSnapshot
	config              config.ValidatedConfig
	makeFeatureData     convertToFeatureData
	repo                execute.OpenRepoResult
	rootDir             gitdomain.RepoRootDir
	verbose             configdomain.Verbose
}
