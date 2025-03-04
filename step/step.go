package step

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
)

const (
	testShardsDirKey = "BITRISE_TEST_SHARDS_PATH"
)

var ErrXcodebuildFailed = errors.New("unexpected xcodebuild failure")

type Input struct {
	ProductPath      string `env:"product_path,required"`
	ShardCount       int    `env:"shard_count,required"`
	ShardCalculation string `env:"shard_calculation,opt[alphabetically]"`
	Destination      string `env:"destination,required"`
	Verbose          bool   `env:"verbose,opt[true,false]"`
}

type Config struct {
	ProductPath      string
	ShardCount       int
	ShardCalculation string
	Destination      string
}

type Result struct {
	TestShardsDir string
}

type Step struct {
	inputParser    stepconf.InputParser
	commandFactory command.Factory
	exporter       export.Exporter
	logger         log.Logger
}

func NewStep(
	inputParser stepconf.InputParser,
	commandFactory command.Factory,
	exporter export.Exporter,
	logger log.Logger,
) Step {
	return Step{
		inputParser:    inputParser,
		commandFactory: commandFactory,
		exporter:       exporter,
		logger:         logger,
	}
}

func (s *Step) ProcessConfig() (*Config, error) {
	var input Input
	err := s.inputParser.Parse(&input)
	if err != nil {
		return &Config{}, err
	}

	stepconf.Print(input)
	s.logger.EnableDebugLog(input.Verbose)

	return &Config{
		ProductPath:      input.ProductPath,
		ShardCount:       input.ShardCount,
		ShardCalculation: input.ShardCalculation,
		Destination:      input.Destination,
	}, nil
}

func (s *Step) Run(config Config) (Result, error) {
	s.logger.Println()
	s.logger.Infof("Collecting tests:")

	var tests []string

	err := retry.Times(3).Wait(10).TryWithAbort(func(attempt uint) (error, bool) {
		s.logger.Infof("%d. attempt", attempt+1)

		testList, err := s.collectTests(config.ProductPath, config.Destination)
		if err != nil {
			if errors.Is(err, ErrXcodebuildFailed) {
				s.logger.Warnf("Test collection failed: %s", err)
				s.logger.Warnf("Retrying...")

				return err, false
			}
			return err, true
		}

		tests = testList

		return nil, true
	})
	if err != nil {
		return Result{}, err
	}
	if len(tests) == 0 {
		return Result{}, fmt.Errorf("no tests found in %s", config.ProductPath)
	}

	s.logger.Printf("Found %d tests in %s", len(tests), config.ProductPath)

	shards := shardAlphabetically(tests, config.ShardCount)

	shardFolder, err := createTempFolder()
	if err != nil {
		return Result{}, err
	}

	for i, shard := range shards {
		shardPath := filepath.Join(shardFolder, fmt.Sprintf("%d", i))

		content := ""
		for _, test := range shard {
			content += fmt.Sprintf("%s\n", test)
		}

		if err := os.WriteFile(shardPath, []byte(content), 0644); err != nil {
			return Result{}, err
		}
	}

	return Result{
		TestShardsDir: shardFolder,
	}, nil
}

func (s *Step) Export(result Result) error {
	s.logger.Println()
	s.logger.Infof("Exporting outputs:")

	if err := s.exporter.ExportOutput(testShardsDirKey, result.TestShardsDir); err != nil {
		s.logger.Warnf("Failed to export: %s: %s", testShardsDirKey, err)
	} else {
		s.logger.Donef("%s: %s", testShardsDirKey, result.TestShardsDir)
	}

	return nil
}

func (s *Step) collectTests(testProductsPath, destination string) ([]string, error) {
	tmpDir, err := createTempFolder()
	if err != nil {
		return nil, err
	}

	testOutput := filepath.Join(tmpDir, "result.txt")
	options := []string{
		"test-without-building",
		"-enumerate-tests",
		"-test-enumeration-format", "json",
		"-test-enumeration-style", "flat",
		"-test-enumeration-output-path", testOutput,
		"-destination", destination,
	}

	if filepath.Ext(testProductsPath) == ".xctestproducts" {
		options = append(options, "-testProductsPath", testProductsPath)
	} else if filepath.Ext(testProductsPath) == ".xctestrun" {
		options = append(options, "-xctestrun", testProductsPath)
	} else {
		return nil, fmt.Errorf("unsupported test products file extension: %s", testProductsPath)
	}

	cmd := s.commandFactory.Create("xcodebuild", options, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stdout,
		Env:    []string{"NSUnbufferedIO=YES"},
	})

	s.logger.TDonef(cmd.PrintableCommandArgs())

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	defer os.Remove(testOutput)

	bytes, err := os.ReadFile(testOutput)
	if err != nil {
		return nil, err
	}

	return processTestResults(bytes)
}

func processTestResults(result []byte) ([]string, error) {
	type testData struct {
		Errors []string `json:"errors"`
		Values []struct {
			Tests []struct {
				Identifier string `json:"identifier"`
			} `json:"enabledTests"`
		} `json:"values"`
	}

	var data testData
	if err := json.Unmarshal(result, &data); err != nil {
		return nil, err
	}

	if len(data.Errors) != 0 {
		return nil, fmt.Errorf("%w: %s", ErrXcodebuildFailed, data.Errors)
	}

	var tests []string
	for _, value := range data.Values {
		for _, test := range value.Tests {
			tests = append(tests, test.Identifier)
		}
	}

	return tests, nil
}

func shardAlphabetically(tests []string, shards int) [][]string {
	slices.Sort(tests)

	buckets := make([][]string, shards)
	bucketSize := (len(tests) + shards - 1) / shards

	for i, test := range tests {
		bucketIndex := i / bucketSize
		buckets[bucketIndex] = append(buckets[bucketIndex], test)
	}

	return buckets
}

func createTempFolder() (string, error) {
	return os.MkdirTemp("", "")
}
