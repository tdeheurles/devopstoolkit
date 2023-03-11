package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"gopkg.in/yaml.v2"
)

func TestVersion(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			// GIVEN
			sc.Step(`^we have a "([^"]*)"$`, func(ctx context.Context, objectType string, data string) (context.Context, error) {
				switch objectType {
				case "Config":
					var config Config
					ctx, config = getConfig(ctx)
					err := yaml.Unmarshal([]byte(data), &config)
					ctx = context.WithValue(ctx, configKey{}, config)
					if err != nil {
						return ctx, err
					}
				}
				return ctx, nil
			})

			sc.Step(`^we expect execution of$`, func(ctx context.Context, table *godog.Table) (context.Context, error) {
				for _, row := range table.Rows[1:] {
					expectedCalls := row.Cells[0].Value
					mockType := row.Cells[1].Value
					methodName := row.Cells[2].Value
					parameters, err := partialTableToSlice("Arg", table.Rows[0], row, ctx)
					if err != nil {
						return ctx, err
					}
					returns, err := partialTableToSlice("Return", table.Rows[0], row, ctx)
					if err != nil {
						return ctx, err
					}

					switch mockType {
					case "binaryExecutor":
						var mock *MockBinaryExecutorer
						ctx, mock = getBinaryExecutor(ctx)
						if expectedCalls != "0" {
							mock.On(methodName, parameters...).Return(returns...)
						}

					case "command":
						var mock *MockCommander
						ctx, mock = getCommand(ctx)
						if expectedCalls != "0" {
							mock.On(methodName, parameters...).Return(returns...)
						}

					case "commandFactory":
						var mock *MockCommandFactorier
						ctx, mock = getCommandFactory(ctx)
						if expectedCalls != "0" {
							mock.On(methodName, parameters...).Return(returns...)
						}

					case "component":
						var mock *MockComponenter
						ctx, mock = getComponent(ctx)
						if expectedCalls != "0" {
							mock.On(methodName, parameters...).Return(returns...)
						}

					case "configuration":
						var mock *MockConfigurationer
						ctx, mock = getConfiguration(ctx)
						if expectedCalls != "0" {
							mock.On(methodName, parameters...).Return(returns...)
						}

					default:
						return ctx, fmt.Errorf("unknown mock type %s", mockType)
					}
				}

				return ctx, nil
			})

			// WHEN
			sc.Step(`^"([^"]*)" "([^"]*)" "([^"]*)" is executed$`, func(ctx context.Context, objectName, version, methodName string) (context.Context, error) {
				ctx, binaryExecutor := getBinaryExecutor(ctx)
				ctx, commandFactory := getCommandFactory(ctx)
				ctx, configuration := getConfiguration(ctx)
				devopsRunner := NewDevopsRunner(version, binaryExecutor, commandFactory, configuration)
				switch objectName {
				case "devopsrunner":
					switch methodName {
					case "Run":
						devopsRunner.Run()
					}
				}
				return ctx, nil
			})

			// THEN
			sc.Step(`^the expectation are met$`, func(ctx context.Context) (context.Context, error) {
				expected := getExpected(ctx)

				for _, mockType := range expected {
					switch mockType {
					case "command":
						command := ctx.Value(commandKey{}).(*MockCommander)
						success := command.AssertExpectations(t)
						if !success {
							return ctx, fmt.Errorf("command expectations not met")
						}
					case "commandFactory":
						commandFactory := ctx.Value(commandFactoryKey{}).(*MockCommandFactorier)
						success := commandFactory.AssertExpectations(t)
						if !success {
							return ctx, fmt.Errorf("commandFactory expectations not met")
						}

					case "component":
						component := ctx.Value(componentKey{}).(*MockComponenter)
						success := component.AssertExpectations(t)
						if !success {
							return ctx, fmt.Errorf("component expectations not met")
						}

					case "binaryExecutor":
						binaryExecutor := ctx.Value(binaryExecutorKey{}).(*MockBinaryExecutorer)
						success := binaryExecutor.AssertExpectations(t)
						if !success {
							return ctx, fmt.Errorf("binaryExecutor expectations not met")
						}
					}
				}

				return ctx, nil
			})
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/version.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
