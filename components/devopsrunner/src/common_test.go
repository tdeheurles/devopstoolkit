package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/messages-go/v16"
)

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
				default:
					return nil, fmt.Errorf("unknown mock type %s", parameterValue)
				}
			}
			if strings.HasPrefix(parameterValue, "bool:") {
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
			} else {
				// default is string without the string: prefix
				returns = append(returns, parameterValue)
			}
		}
	}
	return returns, nil
}

type expectedKey struct{}
type componentKey struct{}
type commandKey struct{}
type commandFactoryKey struct{}
type configurationKey struct{}
type binaryExecutorKey struct{}
