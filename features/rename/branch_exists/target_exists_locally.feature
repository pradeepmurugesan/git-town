Feature: destination branch exists

  Background:
    Given a Git repo with origin

  Scenario: destination branch exists locally
    Given the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
    And the current branch is "alpha"
    When I run "git-town rename alpha beta"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
    And Git Town prints the error:
      """
      there is already a branch "beta"
      """
    And the initial branches and lineage exist now

  Scenario: destination branch exists in origin
    Given the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | origin        |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | origin        | beta commit  |
    And the current branch is "alpha"
    When I run "git-town rename alpha beta"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
    And Git Town prints the error:
      """
      there is already a branch "beta" at the "origin" remote
      """
    And the initial branches and lineage exist now
