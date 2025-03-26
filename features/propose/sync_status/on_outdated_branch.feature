@skipWindows
Feature: sync before proposing

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And a proposal for this branch does not exist
    And the current branch is "child"
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                   |
      | child  | git fetch --prune --tags                                                  |
      | (none) | Looking for proposal online ... ok                                        |
      | child  | git checkout main                                                         |
      | main   | git rebase origin/main --no-update-refs                                   |
      |        | git push                                                                  |
      |        | git checkout parent                                                       |
      | parent | git merge --no-edit --ff main                                             |
      |        | git merge --no-edit --ff origin/parent                                    |
      |        | git push                                                                  |
      |        | git checkout child                                                        |
      | child  | git merge --no-edit --ff parent                                           |
      |        | git merge --no-edit --ff origin/child                                     |
      |        | git push                                                                  |
      | (none) | open https://github.com/git-town/git-town/compare/parent...child?expand=1 |
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parent...child?expand=1
      """
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                  |
      | main   | local, origin | origin main commit                                       |
      |        |               | local main commit                                        |
      | child  | local, origin | local child commit                                       |
      |        |               | Merge branch 'parent' into child                         |
      |        |               | origin child commit                                      |
      |        |               | Merge remote-tracking branch 'origin/child' into child   |
      | parent | local, origin | local parent commit                                      |
      |        |               | Merge branch 'main' into parent                          |
      |        |               | origin parent commit                                     |
      |        |               | Merge remote-tracking branch 'origin/parent' into parent |
