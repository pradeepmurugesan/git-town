Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    And the current branch is "old"
    When I run "git-town rename new --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git remote                                    |
      |        | backend  | git branch --show-current                     |
      | old    | frontend | git fetch --prune --tags                      |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git remote get-url origin                     |
      | old    | frontend | git branch --move old new                     |
      |        | frontend | git checkout new                              |
      |        | backend  | git config git-town-branch.new.parent main    |
      |        | backend  | git config --unset git-town-branch.old.parent |
      | new    | frontend | git push -u origin new                        |
      |        | frontend | git push origin :old                          |
      |        | backend  | git show-ref --verify --quiet refs/heads/old  |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git stash list                                |
    And Git Town prints:
      """
      Ran 23 shell commands.
      """

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git remote get-url origin                     |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git remote get-url origin                     |
      | new    | frontend | git branch old {{ sha 'old commit' }}         |
      |        | frontend | git push -u origin old                        |
      |        | frontend | git checkout old                              |
      | old    | frontend | git branch -D new                             |
      |        | frontend | git push origin :new                          |
      |        | backend  | git config --unset git-town-branch.new.parent |
      |        | backend  | git config git-town-branch.old.parent main    |
    And Git Town prints:
      """
      Ran 17 shell commands.
      """
