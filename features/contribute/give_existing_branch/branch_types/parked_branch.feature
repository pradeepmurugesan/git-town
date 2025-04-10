Feature: make another parked branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    When I run "git-town contribute parked"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "parked" is now a contribution branch
      """
    And the contribution branches are now "parked"
    And there are now no parked branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the parked branches are now "parked"
