Feature: compress the commits on an entire stack when at the stack root

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | alpha  | local, origin | alpha 1 | alpha_1   | alpha 1      |
      |        |               | alpha 2 | alpha_2   | alpha 2      |
      |        |               | alpha 3 | alpha_3   | alpha 3      |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta 1  | beta_1    | beta 1       |
      |        |               | beta 2  | beta_2    | beta 2       |
      |        |               | beta 3  | beta_3    | beta 3       |
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | gamma | feature | beta   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | gamma  | local, origin | gamma 1 | gamma_1   | gamma 1      |
      |        |               | gamma 2 | gamma_2   | gamma 2      |
      |        |               | gamma 3 | gamma_3   | gamma 3      |
    And the current branch is "alpha"
    When I run "git-town compress --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git fetch --prune --tags                        |
      |        | git reset --soft main                           |
      |        | git commit -m "alpha 1"                         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git reset --soft alpha                          |
      |        | git commit -m "beta 1"                          |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout gamma                              |
      | gamma  | git reset --soft beta                           |
      |        | git commit -m "gamma 1"                         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout alpha                              |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE |
      | alpha  | local, origin | alpha 1 |
      | beta   | local, origin | beta 1  |
      | gamma  | local, origin | gamma 1 |
    And file "alpha_1" still has content "alpha 1"
    And file "alpha_2" still has content "alpha 2"
    And file "alpha_3" still has content "alpha 3"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git reset --hard {{ sha 'alpha 3' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git reset --hard {{ sha 'beta 3' }}             |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout gamma                              |
      | gamma  | git reset --hard {{ sha 'gamma 3' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout alpha                              |
    And the initial commits exist now
    And the initial branches and lineage exist now
