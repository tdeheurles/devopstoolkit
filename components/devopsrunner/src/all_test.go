package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
	"github.com/r3labs/diff/v2"
	"gopkg.in/yaml.v2"
)

func TestAll(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {

			// GIVEN
			sc.Step(`^we have a "([^"]*)"$`,
				func(ctx context.Context, objectType string, data string) (context.Context, error) {
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
			sc.Step(`^we expect execution of$`,
				func(ctx context.Context, table *godog.Table) (context.Context, error) {
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
			sc.Step(`^"([^"]*)" "([^"]*)" "([^"]*)" is executed$`,
				func(ctx context.Context, objectName, version, methodName string) (context.Context, error) {
					switch objectName {
					case "devopsrunner":
						var binaryExecutor *MockBinaryExecutorer
						ctx, binaryExecutor = getBinaryExecutor(ctx)

						var commandFactory *MockCommandFactorier
						ctx, commandFactory = getCommandFactory(ctx)

						var configuration *MockConfigurationer
						ctx, configuration = getConfiguration(ctx)

						devopsRunner := NewDevopsRunner(version, binaryExecutor, commandFactory, configuration)
						switch methodName {
						case "Run":
							devopsRunner.Run()
						}
					}
					return ctx, nil
				})
			sc.Step(`^NewConfiguration is executed with$`,
				func(ctx context.Context, table *godog.Table) (context.Context, error) {
					var args []string
					for _, cell := range table.Rows[0].Cells {
						args = append(args, cell.Value)
					}

					return context.WithValue(ctx, configurationKey{}, NewConfiguration(args)), nil
				})
			sc.Step(`^"([^"]*)" "([^"]*)" is executed with$`,
				func(ctx context.Context, objectName, methodName, content string) (context.Context, error) {
					switch objectName {
					case "commandFactory":
						commandFactory := NewCommandFactory()
						switch methodName {
						case "ParseFile":
							devopsRunnerParameters, err := commandFactory.GetDevopsRunnerParameters(content)
							if err != nil {
								return ctx, err
							}

							ctx = context.WithValue(ctx, devopsRunnerParametersKey{}, devopsRunnerParameters)
						}
					}
					return ctx, nil
				})

			// THEN
			sc.Step(`^the expectation are met$`,
				func(ctx context.Context) (context.Context, error) {
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
			sc.Step(`^Configuration Data should be$`,
				func(ctx context.Context, data string) (context.Context, error) {
					configuration := ctx.Value(configurationKey{}).(*Configuration)
					var expected Config
					err := yaml.Unmarshal([]byte(data), &expected)
					if err != nil {
						return ctx, err
					}

					changelog, err := diff.Diff(configuration.Data(), expected)
					if err != nil {
						return ctx, err
					}

					if len(changelog) != 0 {
						println(changelog)
						return ctx, fmt.Errorf("configuration data not as expected")
					}

					return ctx, nil
				})
			sc.Step(`^a slice of DevopsRunnerParameters should be returned$`,
				func(ctx context.Context, table *godog.Table) (context.Context, error) {
					devopsRunnerParameters := ctx.Value(devopsRunnerParametersKey{}).([]DevopsRunnerParameter)
					for index, row := range table.Rows[1:] {
						if len(devopsRunnerParameters) <= index {
							return ctx, fmt.Errorf("parameter %d \"%s\" is missing", index, row.Cells[0].Value)
						}
						if devopsRunnerParameters[index].Name != row.Cells[0].Value {
							return ctx, fmt.Errorf("parameter %d's name is %s but we were expecting %s", index, devopsRunnerParameters[index].Name, row.Cells[0].Value)
						}
						if devopsRunnerParameters[index].Type != row.Cells[1].Value {
							return ctx, fmt.Errorf("parameter %d's type is %s but we were expecting %s", index, devopsRunnerParameters[index].Type, row.Cells[1].Value)
						}
						if devopsRunnerParameters[index].Tag != row.Cells[2].Value {
							return ctx, fmt.Errorf("parameter %d's tag is %s but we were expecting %s", index, devopsRunnerParameters[index].Tag, row.Cells[2].Value)
						}
					}
					return ctx, nil
				})
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

// TODO: merge all the getXXX into a generic function
func getExpected(ctx context.Context) []string {
	var expected []string
	if ctx.Value(expectedKey{}) != nil {
		expected = ctx.Value(expectedKey{}).([]string)
	} else {
		expected = []string{}
	}
	return expected
}

func getConfig(ctx context.Context) (context.Context, Config) {
	var config Config
	if ctx.Value(configKey{}) != nil {
		config = ctx.Value(configKey{}).(Config)
	} else {
		ctx = context.WithValue(ctx, configKey{}, config)
	}
	return ctx, config
}

func getConfiguration(ctx context.Context) (context.Context, *MockConfigurationer) {
	var Configuration *MockConfigurationer
	if ctx.Value(configurationKey{}) != nil {
		Configuration = ctx.Value(configurationKey{}).(*MockConfigurationer)
	} else {
		Configuration = new(MockConfigurationer)
		ctx = context.WithValue(ctx, configurationKey{}, Configuration)
		expected := getExpected(ctx)
		expected = append(expected, "commandFactory")
		ctx = context.WithValue(ctx, expectedKey{}, expected)
	}
	return ctx, Configuration
}

func getBinaryExecutor(ctx context.Context) (context.Context, *MockBinaryExecutorer) {
	var binaryExecutor *MockBinaryExecutorer
	if ctx.Value(binaryExecutorKey{}) != nil {
		binaryExecutor = ctx.Value(binaryExecutorKey{}).(*MockBinaryExecutorer)
	} else {
		binaryExecutor = new(MockBinaryExecutorer)
		ctx = context.WithValue(ctx, binaryExecutorKey{}, binaryExecutor)
		expected := getExpected(ctx)
		expected = append(expected, "commandFactory")
		ctx = context.WithValue(ctx, expectedKey{}, expected)
	}
	return ctx, binaryExecutor
}

func getCommandFactory(ctx context.Context) (context.Context, *MockCommandFactorier) {
	var commandFactory *MockCommandFactorier
	if ctx.Value(commandFactoryKey{}) != nil {
		commandFactory = ctx.Value(commandFactoryKey{}).(*MockCommandFactorier)
	} else {
		commandFactory = new(MockCommandFactorier)
		ctx = context.WithValue(ctx, commandFactoryKey{}, commandFactory)
		expected := getExpected(ctx)
		expected = append(expected, "commandFactory")
		ctx = context.WithValue(ctx, expectedKey{}, expected)
	}
	return ctx, commandFactory
}

func getCommand(ctx context.Context) (context.Context, *MockCommander) {
	var command *MockCommander
	if ctx.Value(commandKey{}) != nil {
		command = ctx.Value(commandKey{}).(*MockCommander)
	} else {
		command = new(MockCommander)
		ctx = context.WithValue(ctx, commandKey{}, command)
		expected := getExpected(ctx)
		expected = append(expected, "commandFactory")
		ctx = context.WithValue(ctx, expectedKey{}, expected)
	}
	return ctx, command
}

func getComponent(ctx context.Context) (context.Context, *MockComponenter) {
	var component *MockComponenter
	if ctx.Value(componentKey{}) != nil {
		component = ctx.Value(componentKey{}).(*MockComponenter)
	} else {
		component = new(MockComponenter)
		ctx = context.WithValue(ctx, componentKey{}, component)
		expected := getExpected(ctx)
		expected = append(expected, "component")
		ctx = context.WithValue(ctx, expectedKey{}, expected)
	}
	return ctx, component
}

func partialTableToSlice(key string, titles, parameters *messages.PickleTableRow, ctx context.Context) ([]interface{}, error) {
	returns := []interface{}{}

	for index, cell := range parameters.Cells {
		if titles.Cells[index].Value == key && cell.Value != "-" {
			parameterValue := cell.Value
			if strings.HasPrefix(parameterValue, "mock:") {
				mockTypeString := strings.Replace(parameterValue, "mock:", "", -1)
				switch mockTypeString {
				case "command":
					returns = append(returns, ctx.Value(commandKey{}).(*MockCommander))
				case "configuration":
					returns = append(returns, ctx.Value(configurationKey{}).(*MockConfigurationer))
				default:
					return nil, fmt.Errorf("unknown mock type %s", parameterValue)
				}
			} else if strings.HasPrefix(parameterValue, "bool:") {
				boolValue := strings.Replace(parameterValue, "bool:", "", -1)
				returns = append(returns, boolValue == "true")
			} else if strings.HasPrefix(parameterValue, "int:") {
				intValue := strings.Replace(parameterValue, "int:", "", -1)
				intValueAsInt, err := strconv.Atoi(intValue)
				if err != nil {
					return nil, err
				}
				returns = append(returns, intValueAsInt)
			} else if strings.HasPrefix(parameterValue, "stringArray:") {
				stringArrayValue := strings.Replace(parameterValue, "stringArray:", "", -1)
				returns = append(returns, strings.Split(stringArrayValue, ","))
			} else if strings.HasPrefix(parameterValue, "stringMap:") {
				stringMapValue := strings.Replace(parameterValue, "stringMap:", "", -1)
				stringMapValueAsMap := map[string]string{}
				for _, stringMapValueAsPair := range strings.Split(stringMapValue, ",") {
					stringMapValueAsPairAsArray := strings.Split(stringMapValueAsPair, "=")
					stringMapValueAsMap[stringMapValueAsPairAsArray[0]] = stringMapValueAsPairAsArray[1]
				}
				returns = append(returns, stringMapValueAsMap)
			} else if strings.HasPrefix(parameterValue, "string:") {
				stringValue := strings.Replace(parameterValue, "string:", "", -1)
				returns = append(returns, stringValue)
			} else if strings.HasPrefix(parameterValue, "struct:") {
				structName := strings.Replace(parameterValue, "struct:", "", -1)
				switch structName {
				case "Config":
					returns = append(returns, ctx.Value(configKey{}).(Config))
				}
			} else if parameterValue == "nil" {
				returns = append(returns, nil)
			} else {
				// default is string without the string: prefix
				returns = append(returns, parameterValue)
			}
		}
	}
	return returns, nil
}

type expectedKey struct{}
type configKey struct{}
type componentKey struct{}
type commandKey struct{}
type commandFactoryKey struct{}
type configurationKey struct{}
type binaryExecutorKey struct{}
type devopsRunnerParametersKey struct{}

// TODO: WIP: attempt to use generic function to get from context
// func getFromContext[V interface{}](ctx context.Context) (context.Context, V) {
// 	var m V

// 	keyAsInterface := interfaceToKeyAsInterface[V]()

// 	if ctx.Value(keyAsInterface) != nil {
// 		m = ctx.Value(keyAsInterface).(V)
// 	} else {
// 		m = new(V)
// 		ctx = context.WithValue(ctx, keyAsInterface, m)

// 		switch any(m).(type) {
// 		case *mock.Mock:
// 			expected := getExpected(ctx)
// 			expected = append(expected, keyAsString)
// 			ctx = context.WithValue(ctx, expectedKey{}, expected)
// 		}
// 	}
// 	return ctx, m
// }

// func interfaceToKeyAsInterface[V interface{}]() interface{} {
// 	var m V
// 	switch any(m).(type) {
// 	case *MockCommander:
// 		return commandKey{}
// 	case *MockCommandFactorier:
// 		return commandFactoryKey{}
// 	case *MockComponenter:
// 		return componentKey{}
// 	case *MockBinaryExecutorer:
// 		return binaryExecutorKey{}
// 	case []string:
// 		return expectedKey{}
// 	default:
// 		panic("unknown type")
// 	}
// }

// func keyToInterface(key string) interface{} {
// 	switch key {
// 	case "command":
// 		return new(MockCommander)
// 	case "commandFactory":
// 		return new(MockCommandFactorier)
// 	case "component":
// 		return new(MockComponenter)
// 	case "binaryExecutor":
// 		return new(MockBinaryExecutorer)
// 	case "expected":
// 		return []string{}
// 	default:
// 		panic("unknown type")
// 	}
// }
