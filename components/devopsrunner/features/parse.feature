Feature: Parsing a file should result 
    Scenario: we want to test commandFactory parsing of file
        When "commandFactory" "ParseFile" is executed with
            """
            #!/bin/bash

            # do: Value string `yaml:"value" default:"foo" usage:"Value of the thing"`
            export value="<% value %>"

            # do: Path string `yaml:"path" default:"/tmp" usage:"Path to the thing"`
            export path="<% path %>"

            """
        Then a slice of DevopsRunnerParameters should be returned
            | Name  | Type   | Tag                                                     |
            | Value | string | `yaml:"value" default:"foo" usage:"Value of the thing"` |
            | Path  | string | `yaml:"path" default:"/tmp" usage:"Path to the thing"`  |

    Scenario: we want to test commandFactory parsing of file with the different form of prompt
        When "commandFactory" "ParseFile" is executed with
            """
            #!/bin/bash

            #do: Foo0 string `yaml:"foo0"`
            export Foo0="<% foo0 %>"

            #do : Foo1 string `yaml:"foo1"`
            export Foo1="<% foo1 %>"

            #do: Foo0 string `yaml:"foo0"`
            export Foo2="<% foo2 %>"
            """
        Then a slice of DevopsRunnerParameters should be returned
            | Name  | Type   | Tag           |
            | Foo0  | string | `yaml:"foo0"` |
            | Foo1  | string | `yaml:"foo1"` |
