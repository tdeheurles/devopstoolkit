Feature: Configuration is able to parse common devopsrunner parameters

    Scenario: we ask for a version switch from cli
        When NewConfiguration is executed with
            |--spec.devopsrunner.version=0.0.2|--config=./foo|
        Then Configuration Data should be
            """
            spec:
                devopsrunner:
                    version: 0.0.2
            """

    Scenario: we ask for a version switch from configuration
        When NewConfiguration is executed with
            |--config=./tests/test0.config.yaml|
        Then Configuration Data should be
            """
            spec:
                devopsrunner:
                    version: 0.0.2
            """

    Scenario: we ask for a command with CLI
        When NewConfiguration is executed with
            |--spec.command.path=foo,bar|
        Then Configuration Data should be
            """
            spec:
                devopsrunner:
                    version: local
                command:
                    path:
                        - foo
                        - bar
            """

    Scenario: we ask for a command with configuration
        When NewConfiguration is executed with
            |--config=./tests/test1.config.yaml|
        Then Configuration Data should be
            """
            spec:
                devopsrunner:
                    version: local
                command:
                    path:
                        - foo
                        - bar
            """

    Scenario: we ask for a command with CLI subcommand
        When NewConfiguration is executed with
            |foo|bar|
        Then Configuration Data should be
            """
            spec:
                devopsrunner:
                    version: local
                command:
                    path:
                        - foo
                        - bar
            """
