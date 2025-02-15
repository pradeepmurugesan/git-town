package configdomain

import (
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/gohacks"
	"github.com/git-town/git-town/v18/internal/gohacks/mapstools"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// PartialConfig contains configuration data as it is stored in the local or global Git configuration.
type PartialConfig struct {
	Aliases                  Aliases
	BitbucketAppPassword     Option[BitbucketAppPassword]
	BitbucketUsername        Option[BitbucketUsername]
	BranchTypeOverrides      BranchTypeOverrides
	ContributionBranches     gitdomain.LocalBranchNames
	ContributionRegex        Option[ContributionRegex]
	DefaultBranchType        Option[BranchType]
	DevRemote                Option[gitdomain.Remote]
	FeatureRegex             Option[FeatureRegex]
	ForgeType                Option[ForgeType]
	GitHubToken              Option[GitHubToken]
	GitLabToken              Option[GitLabToken]
	GitUserEmail             Option[GitUserEmail]
	GitUserName              Option[GitUserName]
	GiteaToken               Option[GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	Lineage                  Lineage
	MainBranch               Option[gitdomain.LocalBranchName]
	NewBranchType            Option[BranchType]
	ObservedBranches         gitdomain.LocalBranchNames
	ObservedRegex            Option[ObservedRegex]
	Offline                  Option[Offline]
	ParkedBranches           gitdomain.LocalBranchNames
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	PrototypeBranches        gitdomain.LocalBranchNames
	PushHook                 Option[PushHook]
	PushNewBranches          Option[PushNewBranches]
	ShipDeleteTrackingBranch Option[ShipDeleteTrackingBranch]
	ShipStrategy             Option[ShipStrategy]
	SyncFeatureStrategy      Option[SyncFeatureStrategy]
	SyncPerennialStrategy    Option[SyncPerennialStrategy]
	SyncPrototypeStrategy    Option[SyncPrototypeStrategy]
	SyncTags                 Option[SyncTags]
	SyncUpstream             Option[SyncUpstream]
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{
		Aliases: Aliases{},
	} //exhaustruct:ignore
}

func NewPartialConfigFromSnapshot(snapshot SingleSnapshot, updateOutdated bool, removeLocalConfigValue removeLocalConfigValueFunc) (PartialConfig, error) {
	ec := gohacks.ErrorCollector{}
	aliases := snapshot.Aliases()
	branchTypeOverrides, err := NewBranchTypeOverridesInSnapshot(snapshot, removeLocalConfigValue)
	ec.Check(err)
	contributionRegex, err := ParseContributionRegex(snapshot[KeyContributionRegex])
	ec.Check(err)
	defaultBranchType, err := ParseBranchType(snapshot[KeyDefaultBranchType])
	ec.Check(err)
	featureRegex, err := ParseFeatureRegex(snapshot[KeyFeatureRegex])
	ec.Check(err)
	forgeType, err := ParseForgeType(snapshot[KeyForgeType])
	ec.Check(err)
	lineage, err := NewLineageFromSnapshot(snapshot, updateOutdated, removeLocalConfigValue)
	ec.Check(err)
	newBranchType, err := ParseBranchType(snapshot[KeyNewBranchType])
	ec.Check(err)
	observedRegex, err := ParseObservedRegex(snapshot[KeyObservedRegex])
	ec.Check(err)
	offline, err := ParseOffline(snapshot[KeyOffline], KeyOffline)
	ec.Check(err)
	perennialRegex, err := ParsePerennialRegex(snapshot[KeyPerennialRegex])
	ec.Check(err)
	pushHook, err := ParsePushHook(snapshot[KeyPushHook], KeyPushHook)
	ec.Check(err)
	pushNewBranches, err := ParsePushNewBranches(snapshot[KeyPushNewBranches], KeyPushNewBranches)
	ec.Check(err)
	shipDeleteTrackingBranch, err := ParseShipDeleteTrackingBranch(snapshot[KeyShipDeleteTrackingBranch], KeyShipDeleteTrackingBranch)
	ec.Check(err)
	shipStrategy, err := ParseShipStrategy(snapshot[KeyShipStrategy])
	ec.Check(err)
	syncFeatureStrategy, err := ParseSyncFeatureStrategy(snapshot[KeySyncFeatureStrategy])
	ec.Check(err)
	syncPerennialStrategy, err := ParseSyncPerennialStrategy(snapshot[KeySyncPerennialStrategy])
	ec.Check(err)
	syncPrototypeStrategy, err := ParseSyncPrototypeStrategy(snapshot[KeySyncPrototypeStrategy])
	ec.Check(err)
	syncTags, err := ParseSyncTags(snapshot[KeySyncTags], KeySyncTags)
	ec.Check(err)
	syncUpstream, err := ParseSyncUpstream(snapshot[KeySyncUpstream], KeySyncUpstream)
	ec.Check(err)
	return PartialConfig{
		Aliases:                  aliases,
		BitbucketAppPassword:     ParseBitbucketAppPassword(snapshot[KeyBitbucketAppPassword]),
		BitbucketUsername:        ParseBitbucketUsername(snapshot[KeyBitbucketUsername]),
		BranchTypeOverrides:      branchTypeOverrides,
		ContributionBranches:     gitdomain.ParseLocalBranchNames(snapshot[KeyContributionBranches]),
		ContributionRegex:        contributionRegex,
		DefaultBranchType:        defaultBranchType,
		DevRemote:                gitdomain.NewRemote(snapshot[KeyDevRemote]),
		FeatureRegex:             featureRegex,
		ForgeType:                forgeType,
		GitHubToken:              ParseGitHubToken(snapshot[KeyGithubToken]),
		GitLabToken:              ParseGitLabToken(snapshot[KeyGitlabToken]),
		GitUserEmail:             ParseGitUserEmail(snapshot[KeyGitUserEmail]),
		GitUserName:              ParseGitUserName(snapshot[KeyGitUserName]),
		GiteaToken:               ParseGiteaToken(snapshot[KeyGiteaToken]),
		HostingOriginHostname:    ParseHostingOriginHostname(snapshot[KeyHostingOriginHostname]),
		Lineage:                  lineage,
		MainBranch:               gitdomain.NewLocalBranchNameOption(snapshot[KeyMainBranch]),
		NewBranchType:            newBranchType,
		ObservedBranches:         gitdomain.ParseLocalBranchNames(snapshot[KeyObservedBranches]),
		ObservedRegex:            observedRegex,
		Offline:                  offline,
		ParkedBranches:           gitdomain.ParseLocalBranchNames(snapshot[KeyParkedBranches]),
		PerennialBranches:        gitdomain.ParseLocalBranchNames(snapshot[KeyPerennialBranches]),
		PerennialRegex:           perennialRegex,
		PrototypeBranches:        gitdomain.ParseLocalBranchNames(snapshot[KeyPrototypeBranches]),
		PushHook:                 pushHook,
		PushNewBranches:          pushNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
	}, ec.Err
}

// a function that deletes the local Git configuration value with the given key
type removeLocalConfigValueFunc func(Key) error

// Merges the given PartialConfig into this configuration object.
func (self PartialConfig) Merge(other PartialConfig) PartialConfig {
	return PartialConfig{
		Aliases:                  mapstools.Merge(other.Aliases, self.Aliases),
		BitbucketAppPassword:     other.BitbucketAppPassword.Or(self.BitbucketAppPassword),
		BitbucketUsername:        other.BitbucketUsername.Or(self.BitbucketUsername),
		BranchTypeOverrides:      other.BranchTypeOverrides.Concat(self.BranchTypeOverrides),
		ContributionBranches:     append(other.ContributionBranches, self.ContributionBranches...),
		ContributionRegex:        other.ContributionRegex.Or(self.ContributionRegex),
		DefaultBranchType:        other.DefaultBranchType.Or(self.DefaultBranchType),
		DevRemote:                other.DevRemote.Or(self.DevRemote),
		FeatureRegex:             other.FeatureRegex.Or(self.FeatureRegex),
		ForgeType:                other.ForgeType.Or(self.ForgeType),
		GitHubToken:              other.GitHubToken.Or(self.GitHubToken),
		GitLabToken:              other.GitLabToken.Or(self.GitLabToken),
		GitUserEmail:             other.GitUserEmail.Or(self.GitUserEmail),
		GitUserName:              other.GitUserName.Or(self.GitUserName),
		GiteaToken:               other.GiteaToken.Or(self.GiteaToken),
		HostingOriginHostname:    other.HostingOriginHostname.Or(self.HostingOriginHostname),
		Lineage:                  other.Lineage.Merge(self.Lineage),
		MainBranch:               other.MainBranch.Or(self.MainBranch),
		NewBranchType:            other.NewBranchType.Or(self.NewBranchType),
		ObservedBranches:         append(other.ObservedBranches, self.ObservedBranches...),
		ObservedRegex:            other.ObservedRegex.Or(self.ObservedRegex),
		Offline:                  other.Offline.Or(self.Offline),
		ParkedBranches:           append(other.ParkedBranches, self.ParkedBranches...),
		PerennialBranches:        append(other.PerennialBranches, self.PerennialBranches...),
		PerennialRegex:           other.PerennialRegex.Or(self.PerennialRegex),
		PrototypeBranches:        append(other.PrototypeBranches, self.PrototypeBranches...),
		PushHook:                 other.PushHook.Or(self.PushHook),
		PushNewBranches:          other.PushNewBranches.Or(self.PushNewBranches),
		ShipDeleteTrackingBranch: other.ShipDeleteTrackingBranch.Or(self.ShipDeleteTrackingBranch),
		ShipStrategy:             other.ShipStrategy.Or(self.ShipStrategy),
		SyncFeatureStrategy:      other.SyncFeatureStrategy.Or(self.SyncFeatureStrategy),
		SyncPerennialStrategy:    other.SyncPerennialStrategy.Or(self.SyncPerennialStrategy),
		SyncPrototypeStrategy:    other.SyncPrototypeStrategy.Or(self.SyncPrototypeStrategy),
		SyncTags:                 other.SyncTags.Or(self.SyncTags),
		SyncUpstream:             other.SyncUpstream.Or(self.SyncUpstream),
	}
}

func (self PartialConfig) ToNormalConfig(defaults NormalConfigData) NormalConfigData {
	syncFeatureStrategy := self.SyncFeatureStrategy.GetOrElse(defaults.SyncFeatureStrategy)
	return NormalConfigData{
		Aliases:                  self.Aliases,
		BitbucketAppPassword:     self.BitbucketAppPassword,
		BitbucketUsername:        self.BitbucketUsername,
		BranchTypeOverrides:      self.BranchTypeOverrides,
		ContributionBranches:     self.ContributionBranches,
		ContributionRegex:        self.ContributionRegex,
		DefaultBranchType:        self.DefaultBranchType.GetOrElse(BranchTypeFeatureBranch),
		DevRemote:                self.DevRemote.GetOrElse(defaults.DevRemote),
		FeatureRegex:             self.FeatureRegex,
		ForgeType:                self.ForgeType,
		GitHubToken:              self.GitHubToken,
		GitLabToken:              self.GitLabToken,
		GiteaToken:               self.GiteaToken,
		HostingOriginHostname:    self.HostingOriginHostname,
		Lineage:                  self.Lineage,
		NewBranchType:            self.NewBranchType.Or(defaults.NewBranchType),
		ObservedBranches:         self.ObservedBranches,
		ObservedRegex:            self.ObservedRegex,
		Offline:                  self.Offline.GetOrElse(defaults.Offline),
		ParkedBranches:           self.ParkedBranches,
		PerennialBranches:        self.PerennialBranches,
		PerennialRegex:           self.PerennialRegex,
		PrototypeBranches:        self.PrototypeBranches,
		PushHook:                 self.PushHook.GetOrElse(defaults.PushHook),
		PushNewBranches:          self.PushNewBranches.GetOrElse(defaults.PushNewBranches),
		ShipDeleteTrackingBranch: self.ShipDeleteTrackingBranch.GetOrElse(defaults.ShipDeleteTrackingBranch),
		ShipStrategy:             self.ShipStrategy.GetOrElse(defaults.ShipStrategy),
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    self.SyncPerennialStrategy.GetOrElse(defaults.SyncPerennialStrategy),
		SyncPrototypeStrategy:    self.SyncPrototypeStrategy.GetOrElse(NewSyncPrototypeStrategyFromSyncFeatureStrategy(syncFeatureStrategy)),
		SyncTags:                 self.SyncTags.GetOrElse(defaults.SyncTags),
		SyncUpstream:             self.SyncUpstream.GetOrElse(defaults.SyncUpstream),
	}
}

func (self PartialConfig) ToUnvalidatedConfig() UnvalidatedConfigData {
	return UnvalidatedConfigData{
		GitUserEmail: self.GitUserEmail,
		GitUserName:  self.GitUserName,
		MainBranch:   self.MainBranch,
	}
}
