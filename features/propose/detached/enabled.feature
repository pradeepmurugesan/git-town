@skipWindows
Feature: GitHub support

  Background:
    Given a Git repo with origin
    And tool "open" is installed
    And a proposal for this branch does not exist
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --tags                                           |
      | (none)  | Looking for proposal online ... ok                                 |
      | feature | git merge --no-edit --ff main                                      |
      |         | git merge --no-edit --ff origin/feature                            |
      | (none)  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
