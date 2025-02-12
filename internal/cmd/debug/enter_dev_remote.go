package debug

import (
	"os"

	"github.com/git-town/git-town/v18/internal/cli/dialog"
	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/spf13/cobra"
)

func enterDevRemote() *cobra.Command {
	return &cobra.Command{
		Use: "dev-remote",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.DevRemote(gitdomain.RemoteOrigin, gitdomain.Remotes{gitdomain.RemoteOrigin, "fork"}, dialogInputs.Next())
			return err
		},
	}
}
