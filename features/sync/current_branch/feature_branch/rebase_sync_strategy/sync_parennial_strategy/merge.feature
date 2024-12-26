Feature: with sync-perennial-strategy set to "merge"

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And Git setting "git-town.sync-perennial-strategy" is "merge"
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git checkout main                               |
      | main    | git merge --no-edit --ff origin/main            |
      |         | git push                                        |
      |         | git checkout feature                            |
      | feature | git rebase main --no-update-refs                |
      |         | git push --force-with-lease --force-if-includes |
      |         | git rebase origin/feature --no-update-refs      |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                    |
      | main    | local, origin | local main commit                          |
      |         |               | origin main commit                         |
      |         |               | Merge remote-tracking branch 'origin/main' |
      | feature | local, origin | origin feature commit                      |
      |         |               | local main commit                          |
      |         |               | origin main commit                         |
      |         |               | local feature commit                       |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                           |
      | feature | git reset --hard {{ sha-before-run 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin feature commit' }}:feature |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                    |
      | main    | local, origin | local main commit                          |
      |         |               | origin main commit                         |
      |         |               | Merge remote-tracking branch 'origin/main' |
      | feature | local         | local feature commit                       |
      |         | origin        | origin feature commit                      |
    And the initial branches and lineage exist now
