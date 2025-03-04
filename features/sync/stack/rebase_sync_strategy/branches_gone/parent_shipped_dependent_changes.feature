Feature: syncing a branch whose parent with dependent changes was shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME | FILE CONTENT   |
      | parent | local, origin | parent commit | file      | parent content |
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | child | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | child  | local, origin | child commit | file      | child content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "parent" branch using the "squash-merge" ship-strategy
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | child  | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git checkout child                      |
      | child  | git pull                                |
      |        | git rebase --onto main parent           |
      |        | git push --force-with-lease             |
      |        | git branch -D parent                    |
    And Git Town prints:
      """
      deleted branch "parent"
      """
    And the current branch is still "child"
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME | FILE CONTENT   |
      | main   | local, origin | parent commit | file      | parent content |
      | child  | local, origin | child commit  | file      | child content  |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                |
      | child  | git reset --hard {{ sha 'child commit' }}              |
      |        | git push --force-with-lease --force-if-includes        |
      |        | git checkout main                                      |
      | main   | git reset --hard {{ sha 'initial commit' }}            |
      |        | git branch parent {{ sha-before-run 'parent commit' }} |
      |        | git checkout child                                     |
    And the current branch is still "child"
    And the initial branches and lineage exist now
