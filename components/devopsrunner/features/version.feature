Feature: devopsrunner should be able to update itself to the version in the configuration

    Scenario: no version is defined in configuration
        And we expect execution of
            | Calls | Mock            | Method    | Arg                       | Return       | Return |
            | 1     | configuration   | GetString | spec.devopsrunner.version | string:      | -      |
            | 1     | command         | Execute   | -                         | int:0        | -      |
            | 1     | commandFactory  | Parse     | -                         | mock:command | -      |
            | 0     | binaryExecutor  | Execute   | -                         | bool:false   | int:0  |
        When "devopsrunner" "0.0.2" "Run" is executed
        Then the expectation are met

    Scenario: devopsrunner is not in correct version
        And we expect execution of
            | Calls | Mock            | Method    | Arg                       | Return    | Return |
            | 1     | configuration   | GetString | spec.devopsrunner.version | 0.0.2     | -      |
            | 0     | command         | Execute   | -                         | -         | -      |
            | 0     | commandFactory  | Parse     | -                         | -         | -      |
            | 1     | binaryExecutor  | Execute   | 0.0.2                     | bool:true | int:0  |
        When "devopsrunner" "0.0.1" "Run" is executed
        Then the expectation are met
    
    Scenario: devopsrunner is in the correct version
        And we expect execution of
            | Calls | Mock            | Method    | Arg                       | Return       | Return |
            | 1     | configuration   | GetString | spec.devopsrunner.version | 0.0.2        | -      |
            | 1     | command         | Execute   | -                         | int:0        | -      |
            | 1     | commandFactory  | Parse     | -                         | mock:command | -      |
            | 0     | binaryExecutor  | Execute   | -                         | bool:false   | int:0  |
        When "devopsrunner" "0.0.2" "Run" is executed
        Then the expectation are met

    # Scenario: no configuration exists
    #     Given devopsrunner is in version "0.0.2"
    #     And the configuration file does not exists
    #     And we have a "command" mock
    #     When devopsrunner "Run" is executed

    Scenario: configuration is local
        And we expect execution of
            | Calls | Mock            | Method    | Arg                       | Return       | Return |
            | 1     | configuration   | GetString | spec.devopsrunner.version | local        | -      |
            | 1     | command         | Execute   | -                         | int:0        | -      |
            | 1     | commandFactory  | Parse     | -                         | mock:command | -      |
            | 0     | binaryExecutor  | Execute   | -                         | bool:false   | int:0  |
        When "devopsrunner" "0.0.2" "Run" is executed
        Then the expectation are met
