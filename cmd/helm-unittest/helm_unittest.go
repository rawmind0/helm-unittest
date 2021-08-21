package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lrills/helm-unittest/internal/printer"
	"github.com/lrills/helm-unittest/pkg/unittest"
	"github.com/lrills/helm-unittest/pkg/unittest/formatter"
	"github.com/spf13/cobra"
)

// testOptions stores options setup by user in command line
type testOptions struct {
	useHelmV3      bool
	useFailfast    bool
	useStrict      bool
	colored        bool
	updateSnapshot bool
	withSubChart   bool
	testFiles      []string
	valuesFiles    []string
	outputFile     string
	outputType     string
}

var testConfig = testOptions{}

var cmd = &cobra.Command{
	Use:   "unittest [flags] CHART [...]",
	Short: "unittest for helm charts",
	Long: `Running chart unittest written in YAML.

This renders your charts locally (without tiller) and
validates the rendered output with the tests defined in
test suite files. Simplest test suite file looks like
below:

---
# CHART_PATH/tests/deployment_test.yaml
suite: test my deployment
templates:
  - deployment.yaml
tests:
  - it: should be a Deployment
    asserts:
      - isKind:
          of: Deployment
---

Put the test files in "tests" directory under your chart
with suffix "_test.yaml", and run:

$ helm unittest my-chart

Or specify the suite files glob path pattern:

$ helm unittest -f 'my-tests/*.yaml' my-chart

Check https://github.com/quintush/helm-unittest for more
details about how to write tests.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, chartPaths []string) {
		var colored *bool
		if cmd.PersistentFlags().Changed("color") {
			colored = &testConfig.colored
		}

		formatter := formatter.NewFormatter(testConfig.outputFile, testConfig.outputType)
		printer := printer.NewPrinter(os.Stdout, colored)
		runner := unittest.TestRunner{
			Printer:        printer,
			Formatter:      formatter,
			UpdateSnapshot: testConfig.updateSnapshot,
			WithSubChart:   testConfig.withSubChart,
			Strict:         testConfig.useStrict,
			Failfast:       testConfig.useFailfast,
			TestFiles:      testConfig.testFiles,
			ValuesFiles:    testConfig.valuesFiles,
			OutputFile:     testConfig.outputFile,
		}
		var passed bool

		if !testConfig.useHelmV3 {
			passed = runner.RunV2(chartPaths)
		} else {
			passed = runner.RunV3(chartPaths)
		}

		if !passed {
			os.Exit(1)
		}
	},
}

// main to execute execute unittest command
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cmd.PersistentFlags().BoolVar(
		&testConfig.colored, "color", false,
		"enforce printing colored output even stdout is not a tty. Set to false to disable color",
	)

	cmd.PersistentFlags().BoolVar(
		&testConfig.useStrict, "strict", false,
		"strict parse the testsuites",
	)

	defaultFilePattern := filepath.Join("tests", "*_test.yaml")
	cmd.PersistentFlags().StringArrayVarP(
		&testConfig.testFiles, "file", "f", []string{defaultFilePattern},
		"glob paths of test files location, default to "+defaultFilePattern,
	)

	cmd.PersistentFlags().StringArrayVarP(
		&testConfig.valuesFiles, "values", "v", []string{},
		"absolute or glob paths of values files location, default no values files",
	)

	cmd.PersistentFlags().BoolVarP(
		&testConfig.updateSnapshot, "update-snapshot", "u", false,
		"update the snapshot cached if needed, make sure you review the change before update",
	)

	cmd.PersistentFlags().BoolVarP(
		&testConfig.withSubChart, "with-subchart", "s", true,
		"include tests of the subcharts within `charts` folder",
	)

	cmd.PersistentFlags().StringVarP(
		&testConfig.outputFile, "output-file", "o", "",
		"output-file the file where testresults are written in JUnit format, defaults no output is written to file",
	)

	cmd.PersistentFlags().StringVarP(
		&testConfig.outputType, "output-type", "t", "XUnit",
		"output-type the file-format where testresults are written in, accepted types are (JUnit, NUnit, XUnit)",
	)

	cmd.PersistentFlags().BoolVarP(
		&testConfig.useHelmV3, "helm3", "3", false,
		"parse helm charts as helm3 charts",
	)

	cmd.PersistentFlags().BoolVarP(
		&testConfig.useFailfast, "failfast", "q", false,
		"direct quit testing, when a test is failed",
	)
}
