Feature: commit message with double-quotes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    When I run "git-town ship -m 'with "double quotes"'"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                              |
      | feature | git fetch --prune --tags             |
      |         | git checkout main                    |
      | main    | git merge --squash --ff feature      |
      |         | git commit -m "with "double quotes"" |
      |         | git push                             |
      |         | git push origin :feature             |
      |         | git branch -D feature                |
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no uncommitted files exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | with "double quotes" |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git revert {{ sha 'with "double quotes"' }}   |
      |        | git push                                      |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u origin feature                    |
      |        | git checkout feature                          |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                       |
      | main    | local, origin | with "double quotes"          |
      |         |               | Revert "with "double quotes"" |
      | feature | local, origin | feature commit                |
    And the initial branches and lineage exist now
