Feature: does not compress empty branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town compress"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "feature" has no commits
      """
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
