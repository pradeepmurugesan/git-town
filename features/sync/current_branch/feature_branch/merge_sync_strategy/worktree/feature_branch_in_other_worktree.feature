Feature: Sync a feature branch that is in another worktree than the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the perennial branches are "main"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And the current branch is "main"
    And branch "feature" is active in another worktree
    When I run "git-town sync" in the other worktree

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff origin/main    |
      |         | git merge --no-edit --ff origin/feature |
      |         | git push                                |
    And the current branch in the other worktree is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE                                                    |
      | main    | local            | local main commit                                          |
      |         | origin           | origin main commit                                         |
      | feature | origin, worktree | local feature commit                                       |
      |         |                  | Merge remote-tracking branch 'origin/main' into feature    |
      |         |                  | origin feature commit                                      |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |
      |         | worktree         | origin main commit                                         |

  Scenario: undo
    When I run "git-town undo" in the other worktree
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                |
      | feature | git reset --hard {{ sha 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin 'origin feature commit' }}:feature |
    And the current branch in the other worktree is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | origin   | origin feature commit |
      |         | worktree | local feature commit  |
    And the initial branches and lineage exist now
