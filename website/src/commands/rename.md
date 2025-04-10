# git town rename

```command-summary
git town rename [<old-name>] <new-name> [-f | --force] [--dry-run] [-v | --verbose]
```

The _rename_ command renames the current branch and its tracking branch. The
branch to rename must be fully [synced](sync.md). Updates all affected
[proposals](propose.md).

Please be aware that most forges are unable to update the head branch (aka
source branch) of proposals. If you rename a branch that already has a proposal,
the existing proposal will most likely end up closed and you have to create a
new proposal that supersedes the old one. If that happens, Git Town will notify
you. Updating proposals of child branches usually works.

## Positional arguments

When called with only one argument, the _rename_ command renames the current
branch to the given name.

When called with two arguments, it renames the branch with the given name to the
given name.

## Options

#### `-f`<br>`--force`

Renaming perennial branches requires confirmation with the `--force` aka `-f`
flag.

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
