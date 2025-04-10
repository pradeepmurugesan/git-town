Feature: already existing remote branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | origin    |
    And an uncommitted file
    When I run "git-town hack existing"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | main     | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout -b existing    |
      | existing | git stash pop               |
      |          | git restore --staged .      |
    And the uncommitted file still exists
    And no commits exist now
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | existing | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout main           |
      | main     | git branch -D existing      |
      |          | git stash pop               |
      |          | git restore --staged .      |
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now
