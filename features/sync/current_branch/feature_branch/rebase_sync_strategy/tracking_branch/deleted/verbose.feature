Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
      | branch-2 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | branch-1 | local, origin | branch-1 commit |
    And the current branch is "branch-2"
    And origin deletes the "branch-2" branch
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town sync --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | TYPE     | COMMAND                                            |
      |          | backend  | git version                                        |
      |          | backend  | git rev-parse --show-toplevel                      |
      |          | backend  | git config -lz --includes --global                 |
      |          | backend  | git config -lz --includes --local                  |
      |          | backend  | git branch -vva --sort=refname                     |
      |          | backend  | git status --long --ignore-submodules              |
      |          | backend  | git remote                                         |
      | branch-2 | frontend | git fetch --prune --tags                           |
      |          | backend  | git stash list                                     |
      |          | backend  | git branch -vva --sort=refname                     |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}          |
      |          | backend  | git remote get-url origin                          |
      |          | backend  | git log main..branch-2 --format=%s --reverse       |
      | branch-2 | frontend | git checkout main                                  |
      |          | backend  | git config --unset git-town-branch.branch-2.parent |
      | main     | frontend | git rebase --onto main branch-2 --no-update-refs   |
      |          | frontend | git branch -D branch-2                             |
      |          | backend  | git show-ref --verify --quiet refs/heads/branch-2  |
      |          | backend  | git show-ref --verify --quiet refs/heads/branch-1  |
      | main     | frontend | git checkout branch-1                              |
      |          | backend  | git branch -vva --sort=refname                     |
      |          | backend  | git config -lz --includes --global                 |
      |          | backend  | git config -lz --includes --local                  |
      |          | backend  | git stash list                                     |
    And Git Town prints:
      """
      Ran 24 shell commands.
      """
    And the branches are now
      | REPOSITORY    | BRANCHES       |
      | local, origin | main, branch-1 |
    And this lineage exists now
      | BRANCH   | PARENT |
      | branch-1 | main   |

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH   | TYPE     | COMMAND                                         |
      |          | backend  | git version                                     |
      |          | backend  | git rev-parse --show-toplevel                   |
      |          | backend  | git config -lz --includes --global              |
      |          | backend  | git config -lz --includes --local               |
      |          | backend  | git status --long --ignore-submodules           |
      |          | backend  | git stash list                                  |
      |          | backend  | git branch -vva --sort=refname                  |
      |          | backend  | git remote get-url origin                       |
      |          | backend  | git rev-parse --verify --abbrev-ref @{-1}       |
      |          | backend  | git remote get-url origin                       |
      | branch-1 | frontend | git branch branch-2 {{ sha 'initial commit' }}  |
      |          | backend  | git show-ref --quiet refs/heads/branch-2        |
      | branch-1 | frontend | git checkout branch-2                           |
      |          | backend  | git config git-town-branch.branch-2.parent main |
    And Git Town prints:
      """
      Ran 14 shell commands.
      """
    And the initial branches and lineage exist now
