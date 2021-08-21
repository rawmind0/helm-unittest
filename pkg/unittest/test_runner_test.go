package unittest_test

import (
	"bytes"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/lrills/helm-unittest/internal/printer"
	. "github.com/lrills/helm-unittest/pkg/unittest"
	"github.com/stretchr/testify/assert"
)

var sectionBeginPattern = regexp.MustCompile("( PASS | FAIL |\n*###|\n*Charts:|\n*Snapshot Summary:)")
var timePattern = regexp.MustCompile(`Time:\s+([\d\.]+)(s|ms)`)

func makeOutputSnapshotable(originalOutput string) []interface{} {
	output := strings.ReplaceAll(originalOutput, "\\", "/")
	timeLoc := timePattern.FindStringSubmatchIndex(output)[2:4]
	timeAgnosticOutput := output[:timeLoc[0]] + "XX.XXX" + output[timeLoc[1]:]

	sectionBeggingLocs := sectionBeginPattern.FindAllStringIndex(timeAgnosticOutput, -1)
	sections := make([]string, len(sectionBeggingLocs))

	suiteBeginIdx := -1
	for sectionIdx := 0; sectionIdx < len(sections); sectionIdx++ {
		start := sectionBeggingLocs[sectionIdx][0]
		var end int
		if sectionIdx >= len(sections)-1 {
			end = len(timeAgnosticOutput)
		} else {
			end = sectionBeggingLocs[sectionIdx+1][0]
		}

		sectionContent := timeAgnosticOutput[start:end]
		sectionBegin := sectionContent[:6]
		if sectionBegin == " PASS " || sectionBegin == " FAIL " {
			sections[sectionIdx] = strings.TrimRight(sectionContent, "\n")
			if suiteBeginIdx == -1 {
				suiteBeginIdx = sectionIdx
			}
		} else {
			sections[sectionIdx] = sectionContent
			if suiteBeginIdx != -1 {
				sort.Strings(sections[suiteBeginIdx:sectionIdx])
				suiteBeginIdx = -1
			}
		}
	}

	sectionsToRetrun := make([]interface{}, len(sections))
	for idx, section := range sections {
		sectionsToRetrun[idx] = section
	}
	return sectionsToRetrun
}

func TestV2RunnerInvalidChartDirFailfast(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		Failfast:  true,
		TestFiles: []string{testTestFiles},
	}
	passed := runner.RunV2([]string{testSuiteTests})
	assert.False(t, passed, buffer.String())
}

func TestV2RunnerInvalidTestSuiteFailfast(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		Strict:    true,
		Failfast:  true,
		TestFiles: []string{testTestFiles},
	}
	passed := runner.RunV2([]string{testV2InvalidBasicChart})
	assert.False(t, passed, buffer.String())
}

func TestV2RunnerOkWithPassedTests(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		TestFiles: []string{testTestFiles},
	}
	passed := runner.RunV2([]string{testV2BasicChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV2RunnerOkWithOverrideValuesPassedTests(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:     printer.NewPrinter(buffer, nil),
		TestFiles:   []string{testTestFiles},
		ValuesFiles: []string{testValuesFiles},
	}
	passed := runner.RunV2([]string{testV2BasicChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV2RunnerOkWithAbsoluteOverrideValuesPassedTests(t *testing.T) {
	buffer := new(bytes.Buffer)
	fullPath, _ := filepath.Abs(testValuesFiles)
	runner := TestRunner{
		Printer:     printer.NewPrinter(buffer, nil),
		TestFiles:   []string{testTestFiles},
		ValuesFiles: []string{fullPath},
	}
	passed := runner.RunV2([]string{testV2BasicChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV2RunnerOkWithFailedTests(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		TestFiles: []string{testTestFailedFiles},
	}
	passed := runner.RunV2([]string{testV2BasicChart})
	assert.False(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV2RunnerOkWithSubSubfolder(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		TestFiles: []string{testTestFiles},
	}
	passed := runner.RunV2([]string{testV2WithSubFolderChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV2RunnerWithTestsInSubchart(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:      printer.NewPrinter(buffer, nil),
		WithSubChart: true,
		TestFiles:    []string{testTestFiles},
	}
	passed := runner.RunV2([]string{testV2WithSubChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV2RunnerWithTestsInSubchartButFlagFalse(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:      printer.NewPrinter(buffer, nil),
		WithSubChart: false,
		TestFiles:    []string{testTestFiles},
	}
	passed := runner.RunV2([]string{testV2WithSubChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV3RunnerInvalidChartDirFailfast(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		Failfast:  true,
		TestFiles: []string{testTestFiles},
	}
	passed := runner.RunV3([]string{testTestFiles})
	assert.False(t, passed, buffer.String())
}

func TestV3RunnerInvalidTestSuiteFailfast(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		Strict:    false,
		Failfast:  true,
		TestFiles: []string{testTestFiles},
	}
	passed := runner.RunV3([]string{testV3InvalidBasicChart})
	assert.False(t, passed, buffer.String())
}

func TestV3RunnerOkWithPassedTests(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		TestFiles: []string{testTestFiles},
	}
	passed := runner.RunV3([]string{testV3BasicChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV3RunnerOkWithSubSubChartsPassedTests(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		TestFiles: []string{testTestFiles},
	}
	passed := runner.RunV3([]string{testV3WithSubSubFolderChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV3RunnerOkWithOverrideValuesPassedTests(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:     printer.NewPrinter(buffer, nil),
		TestFiles:   []string{testTestFiles},
		ValuesFiles: []string{testValuesFiles},
	}
	passed := runner.RunV3([]string{testV3BasicChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV3RunnerOkWithAbsoluteOverrideValuesPassedTests(t *testing.T) {
	buffer := new(bytes.Buffer)
	fullPath, _ := filepath.Abs(testValuesFiles)
	runner := TestRunner{
		Printer:     printer.NewPrinter(buffer, nil),
		TestFiles:   []string{testTestFiles},
		ValuesFiles: []string{fullPath},
	}
	passed := runner.RunV3([]string{testV3BasicChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV3RunnerOkWithFailedTests(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		TestFiles: []string{testTestFailedFiles},
	}
	passed := runner.RunV3([]string{testV3BasicChart})
	assert.False(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV3RunnerOkWithSubSubfolder(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(buffer, nil),
		TestFiles: []string{testTestFiles},
	}
	passed := runner.RunV3([]string{testV3WithSubFolderChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV3RunnerWithTestsInSubchart(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:      printer.NewPrinter(buffer, nil),
		WithSubChart: true,
		TestFiles:    []string{testTestFiles},
	}
	passed := runner.RunV3([]string{testV3WithSubChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}

func TestV3RunnerWithTestsInSubchartButFlagFalse(t *testing.T) {
	buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:      printer.NewPrinter(buffer, nil),
		WithSubChart: false,
		TestFiles:    []string{testTestFiles},
	}
	passed := runner.RunV3([]string{testV3WithSubChart})
	assert.True(t, passed, buffer.String())
	cupaloy.SnapshotT(t, makeOutputSnapshotable(buffer.String())...)
}
