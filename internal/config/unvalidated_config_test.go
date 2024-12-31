package config_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestUnvalidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("Reload", func(t *testing.T) {
		t.Parallel()
		t.Run("lineage changed", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			repo.CreateFeatureBranch("branch", "main")
			repo.Config.Reload()
			want := configdomain.NewLineageWith(configdomain.LineageData{
				"branch": "main",
			})
			must.Eq(t, want, repo.Config.NormalConfig.Lineage)
		})
		t.Run("contribution branches changed", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			repo.CreateBranch("branch", "main")
			err := repo.Config.NormalConfig.SetBranchTypeOverride(configdomain.BranchTypeContributionBranch, "branch")
			must.NoError(t, err)
			repo.Config.Reload()
			must.Eq(t, configdomain.BranchTypeContributionBranch, repo.Config.BranchType("branch"))
		})
	})
}
