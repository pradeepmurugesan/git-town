package configfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v17/internal/messages"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// Decode converts the given config file TOML source into Go data.
func Decode(text string) (*Data, error) {
	var result Data
	_, err := toml.Decode(text, &result)
	return &result, err
}

func Load(rootDir gitdomain.RepoRootDir, fileName string, finalMessages stringslice.Collector) (Option[configdomain.PartialConfig], error) {
	configPath := filepath.Join(rootDir.String(), fileName)
	file, err := os.Open(configPath)
	if err != nil {
		return None[configdomain.PartialConfig](), nil
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return None[configdomain.PartialConfig](), fmt.Errorf(messages.ConfigFileCannotRead, fileName, err)
	}
	configFileData, err := Decode(string(bytes))
	if err != nil {
		return None[configdomain.PartialConfig](), fmt.Errorf(messages.ConfigFileInvalidContent, fileName, err)
	}
	result, err := Validate(*configFileData, finalMessages)
	return Some(result), err
}

// Validate converts the given low-level configfile data into high-level config data.
func Validate(data Data, finalMessages stringslice.Collector) (configdomain.PartialConfig, error) {
	var err error
	var contributionRegex Option[configdomain.ContributionRegex]
	var defaultBranchType Option[configdomain.BranchType]
	var devRemote Option[gitdomain.Remote]
	var featureRegex Option[configdomain.FeatureRegex]
	var hostingOriginHostname Option[configdomain.HostingOriginHostname]
	var hostingPlatform Option[configdomain.HostingPlatform]
	var mainBranch Option[gitdomain.LocalBranchName]
	var newBranchType Option[configdomain.BranchType]
	var observedRegex Option[configdomain.ObservedRegex]
	var perennialBranches gitdomain.LocalBranchNames
	var perennialRegex Option[configdomain.PerennialRegex]
	var pushNewBranches Option[configdomain.PushNewBranches]
	var pushHook Option[configdomain.PushHook]
	var shipDeleteTrackingBranch Option[configdomain.ShipDeleteTrackingBranch]
	var shipStrategy Option[configdomain.ShipStrategy]
	var syncFeatureStrategy Option[configdomain.SyncFeatureStrategy]
	var syncPerennialStrategy Option[configdomain.SyncPerennialStrategy]
	var syncPrototypeStrategy Option[configdomain.SyncPrototypeStrategy]
	var syncTags Option[configdomain.SyncTags]
	var syncUpstream Option[configdomain.SyncUpstream]
	// load legacy definitions first, so that the proper definitions loaded later override them
	if data.CreatePrototypeBranches != nil {
		newBranchType = Some(configdomain.BranchTypePrototypeBranch)
		finalMessages.Add(messages.CreatePrototypeBranchesDeprecation)
	}
	if data.PushNewbranches != nil {
		pushNewBranches = Some(configdomain.PushNewBranches(*data.PushNewbranches))
	}
	if data.PushHook != nil {
		pushHook = Some(configdomain.PushHook(*data.PushHook))
	}
	if data.ShipDeleteTrackingBranch != nil {
		shipDeleteTrackingBranch = Some(configdomain.ShipDeleteTrackingBranch(*data.ShipDeleteTrackingBranch))
	}
	if data.ShipStrategy != nil {
		shipStrategy = Some(configdomain.ShipStrategy(*data.ShipStrategy))
	}
	if data.SyncTags != nil {
		syncTags = Some(configdomain.SyncTags(*data.SyncTags))
	}
	if data.SyncUpstream != nil {
		syncUpstream = Some(configdomain.SyncUpstream(*data.SyncUpstream))
	}
	// load proper definitions, overriding the values from the legacy definitions that were loaded above
	if data.Branches != nil {
		if data.Branches.Main != nil {
			mainBranch = gitdomain.NewLocalBranchNameOption(*data.Branches.Main)
		}
		perennialBranches = gitdomain.NewLocalBranchNames(data.Branches.Perennials...)
		if data.Branches.PerennialRegex != nil {
			perennialRegex, err = configdomain.ParsePerennialRegex(*data.Branches.PerennialRegex)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.Branches.DefaultType != nil {
			defaultBranchType, err = configdomain.ParseBranchType(*data.Branches.DefaultType)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.Branches.FeatureRegex != nil {
			verifiedRegexOpt, err := configdomain.ParseRegex(*data.Branches.FeatureRegex)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				featureRegex = Some(configdomain.FeatureRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.ContributionRegex != nil {
			verifiedRegexOpt, err := configdomain.ParseRegex(*data.Branches.ContributionRegex)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				contributionRegex = Some(configdomain.ContributionRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.ObservedRegex != nil {
			verifiedRegexOpt, err := configdomain.ParseRegex(*data.Branches.ObservedRegex)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				observedRegex = Some(configdomain.ObservedRegex{VerifiedRegex: verifiedRegex})
			}
		}
	}
	if data.Create != nil {
		if data.Create.NewBranchType != nil {
			parsed, err := configdomain.ParseBranchType(*data.Create.NewBranchType)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
			newBranchType = parsed
		}
		if data.Create.PushNewbranches != nil {
			pushNewBranches = Some(configdomain.PushNewBranches(*data.Create.PushNewbranches))
		}
	}
	if data.Hosting != nil {
		if data.Hosting.DevRemote != nil {
			devRemote = gitdomain.NewRemote(*data.Hosting.DevRemote)
		}
		if data.Hosting.Platform != nil {
			hostingPlatform, err = configdomain.ParseHostingPlatform(*data.Hosting.Platform)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.Hosting.OriginHostname != nil {
			hostingOriginHostname = configdomain.ParseHostingOriginHostname(*data.Hosting.OriginHostname)
		}
	}
	if data.Ship != nil {
		if data.Ship.DeleteTrackingBranch != nil {
			shipDeleteTrackingBranch = Some(configdomain.ShipDeleteTrackingBranch(*data.Ship.DeleteTrackingBranch))
		}
		if data.Ship.Strategy != nil {
			shipStrategy = Some(configdomain.ShipStrategy(*data.Ship.Strategy))
		}
	}
	if data.SyncStrategy != nil {
		if data.SyncStrategy.FeatureBranches != nil {
			syncFeatureStrategy, err = configdomain.ParseSyncFeatureStrategy(*data.SyncStrategy.FeatureBranches)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.SyncStrategy.PerennialBranches != nil {
			syncPerennialStrategy, err = configdomain.ParseSyncPerennialStrategy(*data.SyncStrategy.PerennialBranches)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.SyncStrategy.PrototypeBranches != nil {
			syncPrototypeStrategy, err = configdomain.ParseSyncPrototypeStrategy(*data.SyncStrategy.PrototypeBranches)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
	}
	if data.Sync != nil {
		if data.Sync.FeatureStrategy != nil {
			syncFeatureStrategy, err = configdomain.ParseSyncFeatureStrategy(*data.Sync.FeatureStrategy)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.Sync.PerennialStrategy != nil {
			syncPerennialStrategy, err = configdomain.ParseSyncPerennialStrategy(*data.Sync.PerennialStrategy)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.Sync.PrototypeStrategy != nil {
			syncPrototypeStrategy, err = configdomain.ParseSyncPrototypeStrategy(*data.Sync.PrototypeStrategy)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.Sync.PushHook != nil {
			pushHook = Some(configdomain.PushHook(*data.Sync.PushHook))
		}
		if data.Sync.Tags != nil {
			syncTags = Some(configdomain.SyncTags(*data.Sync.Tags))
		}
		if data.Sync.Upstream != nil {
			syncUpstream = Some(configdomain.SyncUpstream(*data.Sync.Upstream))
		}
	}
	return configdomain.PartialConfig{
		Aliases:                  map[configdomain.AliasableCommand]string{},
		BitbucketAppPassword:     None[configdomain.BitbucketAppPassword](),
		BitbucketUsername:        None[configdomain.BitbucketUsername](),
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{},
		ContributionBranches:     gitdomain.LocalBranchNames{},
		ContributionRegex:        contributionRegex,
		DefaultBranchType:        defaultBranchType,
		DevRemote:                devRemote,
		FeatureRegex:             featureRegex,
		GitHubToken:              None[configdomain.GitHubToken](),
		GitLabToken:              None[configdomain.GitLabToken](),
		GitUserEmail:             None[configdomain.GitUserEmail](),
		GitUserName:              None[configdomain.GitUserName](),
		GiteaToken:               None[configdomain.GiteaToken](),
		HostingOriginHostname:    hostingOriginHostname,
		HostingPlatform:          hostingPlatform,
		Lineage:                  configdomain.Lineage{},
		MainBranch:               mainBranch,
		NewBranchType:            newBranchType,
		ObservedBranches:         gitdomain.LocalBranchNames{},
		ObservedRegex:            observedRegex,
		Offline:                  None[configdomain.Offline](),
		ParkedBranches:           gitdomain.LocalBranchNames{},
		PerennialBranches:        perennialBranches,
		PerennialRegex:           perennialRegex,
		PrototypeBranches:        gitdomain.LocalBranchNames{},
		PushHook:                 pushHook,
		PushNewBranches:          pushNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
	}, nil
}
