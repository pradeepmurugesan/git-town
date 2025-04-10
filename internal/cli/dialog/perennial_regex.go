package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/messages"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

const (
	perennialRegexTitle = `Regular expression for perennial branches`
	PerennialRegexHelp  = `
All branches whose name matches this regular expression
are also considered perennial branches.

If you are not sure, leave this empty.

`
)

func PerennialRegex(oldValue Option[configdomain.PerennialRegex], inputs components.TestInput) (Option[configdomain.PerennialRegex], bool, error) {
	value, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          PerennialRegexHelp,
		Prompt:        "Perennial regex: ",
		TestInput:     inputs,
		Title:         perennialRegexTitle,
	})
	if err != nil {
		return None[configdomain.PerennialRegex](), false, err
	}
	fmt.Printf(messages.PerennialRegex, components.FormattedSelection(value, aborted))
	perennialRegex, err := configdomain.ParsePerennialRegex(value)
	return perennialRegex, aborted, err
}
