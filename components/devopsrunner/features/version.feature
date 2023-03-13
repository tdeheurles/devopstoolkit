Feature: devopsrunner should be able to update itself to the version in the configuration

    Scenario: no version is defined in configuration
        Given we have a "Config"
            """
            """
        And we expect execution of
            | Calls | Mock            | Method  | Arg | Return        | Return |
            | 1     | configuration   | Data    | -   | struct:Config | -      |
            | 1     | command         | Execute | -   | int:0         | -      |
            | 1     | commandFactory  | Parse   | -   | mock:command  | nil    |
            | 0     | binaryExecutor  | Execute | -   | bool:false    | int:0  |
        When "devopsrunner" "0.0.2" "Run" is executed
        Then the expectation are met

    Scenario: devopsrunner is not in correct version
        Given we have a "Config"
            """
            spec:
                devopsrunner:
                    version: "0.0.2"
            """
        And we expect execution of
            | Calls | Mock            | Method  | Arg    | Return        | Return |
            | 1     | configuration   | Data    | -      | struct:Config | -      |
            | 0     | command         | Execute | -      | -             | -      |
            | 0     | commandFactory  | Parse   | -      | -             | -      |
            | 1     | binaryExecutor  | Execute | 0.0.2  | bool:true     | int:0  |
        When "devopsrunner" "0.0.1" "Run" is executed
        Then the expectation are met
    
    Scenario: devopsrunner is in the correct version
        Given we have a "Config"
            """
            spec:
                devopsrunner:
                    version: "0.0.2"
            """
        And we expect execution of
            | Calls | Mock            | Method  | Arg | Return        | Return |
            | 1     | configuration   | Data    | -   | struct:Config | -      |
            | 1     | command         | Execute | -   | int:0         | -      |
            | 1     | commandFactory  | Parse   | -   | mock:command  | nil    |
            | 0     | binaryExecutor  | Execute | -   | bool:false    | int:0  |
        When "devopsrunner" "0.0.2" "Run" is executed
        Then the expectation are met

    # Scenario: no configuration exists
    #     Given devopsrunner is in version "0.0.2"
    #     And the configuration file does not exists
    #     And we have a "command" mock
    #     When devopsrunner "Run" is executed

    Scenario: configuration is local
        Given we have a "Config"
            """
            spec:
                devopsrunner:
                    version: "local"
            """
        And we expect execution of
            | Calls | Mock            | Method  | Arg | Return        | Return |
            | 1     | configuration   | Data    | -   | struct:Config | -      |
            | 1     | command         | Execute | -   | int:0         | -      |
            | 1     | commandFactory  | Parse   | -   | mock:command  | nil    |
            | 0     | binaryExecutor  | Execute | -   | bool:false    | int:0  |
        When "devopsrunner" "0.0.2" "Run" is executed
        Then the expectation are met
