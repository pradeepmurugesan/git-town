Feature: override the configured ship-strategy

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And the current branch is "feature"
    And Git setting "git-town.ship-strategy" is "api"
    When I run "git-town ship --strategy=fast-forward"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git checkout main           |
      | main    | git merge --ff-only feature |
      |         | git branch -D feature       |
    And the branches are now
      | REPOSITORY | BRANCHES |
      | local      | main     |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE        |
      | main   | local    | feature commit |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git reset --hard {{ sha 'initial commit' }}   |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout feature                          |
    And the initial commits exist now
    And the initial branches and lineage exist now
