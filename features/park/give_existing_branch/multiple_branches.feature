Feature: parking multiple other branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS |
      | feature      | feature      | main   | local     |
      | contribution | contribution |        | local     |
      | observed     | observed     | main   | local     |
      | prototype    | prototype    | main   | local     |
    When I run "git-town park feature contribution observed prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "feature" is now parked
      """
    And branch "feature" now has type "parked"
    And branch "contribution" now has type "parked"
    And there are now no contribution branches
    And branch "observed" now has type "parked"
    And there are now no observed branches
    And branch "prototype" now has type "parked"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And there are now no parked branches
