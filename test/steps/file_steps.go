package steps

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

// FileSteps defines Cucumber step implementations around files.
func FileSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I don't have any uncommitted files$`, func() error {
		files, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.UncommittedFiles()
		if err != nil {
			return fmt.Errorf("cannot determine uncommitted files: %w", err)
		}
		if len(files) > 0 {
			return fmt.Errorf("unexpected uncommitted files: %s", files)
		}
		return nil
	})

	suite.Step(`^my uncommitted file is stashed$`, func() error {
		uncommittedFiles, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.UncommittedFiles()
		if err != nil {
			return err
		}
		for ucf := range uncommittedFiles {
			if uncommittedFiles[ucf] == fs.activeScenarioState.uncommittedFileName {
				return fmt.Errorf("expected file %q to be stashed but it is still uncommitted", fs.activeScenarioState.uncommittedFileName)
			}
		}
		stashSize, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.StashSize()
		if err != nil {
			return err
		}
		if stashSize != 1 {
			return fmt.Errorf("expected 1 stash but found %d", stashSize)
		}
		return nil
	})

	suite.Step(`^my workspace still contains the file "([^"]*)" with content "([^"]*)"$`, func(file, expectedContent string) error {
		actualContent, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.FileContent(file)
		if err != nil {
			return err
		}
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	suite.Step(`^my repository (?:now|still) has the following committed files$`, func(table *messages.PickleStepArgument_PickleTable) error {
		fileTable, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.FilesInBranches()
		if err != nil {
			return fmt.Errorf("cannot determine files in branches in the developer repo: %w", err)
		}
		diff, errorCount := fileTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing files\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching files found, see diff above")
		}
		return nil
	})
}
