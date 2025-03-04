package step

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-xcode-test-shard-calculation/mocks"
)

func TestRun(t *testing.T) {
	step, testingMocks := createStepAndMocks(t)
	testingMocks.command.On("PrintableCommandArgs").Return("", nil)
	testingMocks.command.On("Run").Return(nil)

	call := testingMocks.commandFactory.On("Create", "xcodebuild", mock.Anything, mock.Anything)
	call.RunFn = func(arguments mock.Arguments) {
		saveTestData(t, arguments)

		call.ReturnArguments = mock.Arguments{testingMocks.command, nil}
	}

	result, err := step.Run(Config{
		ProductPath: "test.xctestrun",
		ShardCount:  2,
		Destination: "test-device",
	})
	require.NoError(t, err)

	files, err := contents(t, result.TestShardsDir)
	require.NoError(t, err)

	expectedFiles := map[string]string{
		"0": "Target/Class/test1\nTarget/Class/test2\nTarget/Class/test3\n",
		"1": "Target/Class/test4\nTarget/Class/test5\nTarget/Class/test6\n",
	}
	assert.Equal(t, expectedFiles, files)
}

func contents(t *testing.T, path string) (map[string]string, error) {
	items, err := os.ReadDir(path)
	require.NoError(t, err)

	files := make(map[string]string)

	for _, item := range items {
		if item.IsDir() {
			t.Fail()
		}

		bytes, err := os.ReadFile(filepath.Join(filepath.Clean(path), item.Name()))
		if err != nil {
			return nil, err
		}

		files[item.Name()] = string(bytes)
	}

	return files, nil
}

func TestSupportedProductTypes(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectedErr string
	}{
		{
			name: "xctestrun supported",
			path: "/path/to/test.xctestrun",
		},
		{
			name: "xctestproducts supported",
			path: "/path/to/test.xctestproducts",
		},
		{
			name:        "xcodeproj is not supported",
			path:        "/path/to/test.xcodeproj",
			expectedErr: "unsupported test products file extension: /path/to/test.xcodeproj",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step, testingMocks := createStepAndMocks(t)
			testingMocks.command.On("PrintableCommandArgs").Return("", nil)
			testingMocks.command.On("Run").Return(nil)

			call := testingMocks.commandFactory.On("Create", "xcodebuild", mock.Anything, mock.Anything)
			call.RunFn = func(arguments mock.Arguments) {
				saveTestData(t, arguments)

				call.ReturnArguments = mock.Arguments{testingMocks.command, nil}
			}

			_, err := step.Run(Config{
				ProductPath: tt.path,
				ShardCount:  2,
				Destination: "test-device",
			})

			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func saveTestData(t *testing.T, arguments mock.Arguments) {
	path := ""
	params := arguments[1].([]string)

	for _, param := range params {
		if strings.Contains(param, "result.txt") {
			path = param
			break
		}
	}

	if path == "" {
		t.Fail()
	}

	data := `
{
  "values" : [
    {
      "enabledTests" : [
        {
          "identifier" : "Target/Class/test1"
        },
        {
          "identifier" : "Target/Class/test2"
        },
        {
          "identifier" : "Target/Class/test3"
        },
        {
          "identifier" : "Target/Class/test4"
        },
        {
          "identifier" : "Target/Class/test5"
        },
        {
          "identifier" : "Target/Class/test6"
        }
      ],
      "testPlan" : "Tests"
    }
  ]
}
`
	err := os.WriteFile(path, []byte(data), 0644)
	require.NoError(t, err)
}

func TestExport(t *testing.T) {
	step, testingMocks := createStepAndMocks(t)
	result := Result{TestShardsDir: "path/to/test/shards"}

	testingMocks.command.On("RunAndReturnTrimmedCombinedOutput").Return("", nil)
	testingMocks.commandFactory.On("Create", "envman", []string{"add", "--key", "BITRISE_TEST_SHARDS_PATH", "--value", result.TestShardsDir}, mock.Anything).Return(testingMocks.command)

	err := step.Export(result)
	require.NoError(t, err)

	testingMocks.commandFactory.AssertExpectations(t)
}

type testingMocks struct {
	envRepository  *mocks.Repository
	commandFactory *mocks.Factory
	command        *mocks.Command
}

func createStepAndMocks(t *testing.T) (Step, testingMocks) {
	envRepository := new(mocks.Repository)
	inputParser := stepconf.NewInputParser(envRepository)
	commandFactory := new(mocks.Factory)
	command := new(mocks.Command)
	exporter := export.NewExporter(commandFactory)
	step := NewStep(inputParser, commandFactory, exporter, log.NewLogger())

	m := testingMocks{
		envRepository:  envRepository,
		commandFactory: commandFactory,
		command:        command,
	}

	return step, m
}
