Feature: park a branch verbosely

  Background:
    Given the current branch is a feature branch "branch"
    And an uncommitted file
    When I run "git-town park --verbose"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                    |
      |        | git version                                |
      |        | git config -lz --global                    |
      |        | git config -lz --local                     |
      |        | git rev-parse --show-toplevel              |
      |        | git status --long --ignore-submodules      |
      |        | git rev-parse --abbrev-ref HEAD            |
      |        | git config git-town.parked-branches branch |
      |        | git config -lz --global                    |
      |        | git config -lz --local                     |
    And it prints:
      """
      Ran 9 shell commands
      """
    And the current branch is still "branch"
    And the uncommitted file still exists
    And branch "branch" is now parked

  Scenario: undo
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      |        | git version                                 |
      |        | git config -lz --global                     |
      |        | git config -lz --local                      |
      |        | git rev-parse --show-toplevel               |
      |        | git stash list                              |
      |        | git status --long --ignore-submodules       |
      |        | git branch -vva                             |
      |        | git rev-parse --verify --abbrev-ref @{-1}   |
      |        | git remote get-url origin                   |
      | branch | git add -A                                  |
      |        | git stash                                   |
      | <none> | git config --unset git-town.parked-branches |
      |        | git show-ref --verify --quiet refs/heads/   |
      |        | git stash list                              |
      | branch | git stash pop                               |
    And it prints:
      """
      Ran 15 shell commands
      """
    And the current branch is still "branch"
    And the uncommitted file still exists
    And branch "branch" is now a feature branch
