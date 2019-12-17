package kubetest

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestCalculateSuccessRatio(t *testing.T) {

	log.SetVerbosity(0)

	var passed, total = 1, 3
	expectedRatio := float64(passed) / float64(total) * 100

	successfullyCalculateThreshold := &testResult{
		passed: passed,
		total:  total,
	}

	assert.Equal(t, expectedRatio, successfullyCalculateThreshold.calculateSuccessRatio(), "Calculate success ratio")

	expectError := &testResult{
		total: 0,
	}

	if os.Getenv("FLAG") == "1" {
		expectError.calculateSuccessRatio()
		return
	}

	testOSExit(t, "TestCalculateSuccessRatio")

}

func TestCheckThresholdPass(t *testing.T) {

	log.SetVerbosity(0)

	var success, threshold = 1.0, 2.0

	failThreshold := &testResult{
		successRatio: success,
		threshold:    threshold,
	}

	if os.Getenv("FLAG") == "1" {
		failThreshold.checkThresholdPass()
		return
	}

	testOSExit(t, "TestCheckThresholdPass")

}

func TestAddResultToRow(t *testing.T) {

	log.SetVerbosity(0)

	emptyRuns := &testResult{
		total:    1,
		testRuns: nil,
	}

	testAddingToExistingRow := &testResult{
		total: 1,
		testRuns: [][]string{{
			"hello",
		}},
	}

	testAddingToUninitialisedRow := &testResult{
		total:    3,
		testRuns: nil,
	}

	tables := []struct {
		message  string
		obj      *testResult
		row      int
		addStr   string
		checkRow int
		checkCol int
		expected string
	}{
		{"Add string to uninitialised list", emptyRuns, 0, "hello", 0, 0, "hello"},
		{"Add string to initialised row with existing list", testAddingToExistingRow, 0, "world", 0, 1, "world"},
		{"Add string to initialised row", testAddingToUninitialisedRow, 2, "world", 2, 0, "world"},
	}

	for _, table := range tables {
		table.obj.addResultToRow(table.row, table.addStr)
		assert.Equal(t, table.expected, table.obj.testRuns[table.checkRow][table.checkCol], table.message)
	}

}

func testOSExit(t *testing.T, testName string) {

	// Run the test in a subprocess #nosec
	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()

	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	expectedErrorString := "exit status 1"
	assert.Equal(t, true, ok)
	assert.Equal(t, expectedErrorString, e.Error())
}
