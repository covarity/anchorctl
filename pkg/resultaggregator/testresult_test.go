package resultaggregator

import (
	"anchorctl/pkg/logging"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"os/exec"
	"testing"
)

func TestNewTestResult(t *testing.T) {
	result := NewTestResult(2, logging.Log)

	assert.Equal(t, 2, result.total, "Expected total of TestResult to be 2.")
}

func TestNewTestRun(t *testing.T) {
	result := NewTestRun("SomeTest", "Some test step")

	assert.Equal(t, "SomeTest", result.TestType, "Expected test type as SomeTest.")
	assert.Equal(t, false, result.Passed, "Expected test passed to be false.")
	assert.Equal(t, false, result.Invalid, "Expected test invalid to be false.")
}

func TestAddInvalidTestRun(t *testing.T) {
	errorMessage := "Test Error"
	log.SetVerbosity(0)
	testResult := NewTestResult(2, logging.Log)

	testResult.AddInvalidTestRun("SomeTest", fmt.Errorf(errorMessage))

	assert.Equal(t, "Test Error", testResult.testRuns[0].Desc, "Expected test run desc to be error string")
	assert.Equal(t, "SomeTest", testResult.testRuns[0].TestType, "Expected test type to be SomeTest")
	assert.Equal(t, true, testResult.testRuns[0].Invalid, "Expected test to be Invalid")
}

func TestAddRun(t *testing.T) {
	result := NewTestResult(1, logging.Log)

	result.AddRun(&TestRun{
		TestType: "SomeTest",
		Desc:     "Some test type",
		Passed:   true,
	})

	assert.Equal(t, "SomeTest", result.testRuns[0].TestType, "Expected test type to be SomeTest")
	assert.Equal(t, "Some test type", result.testRuns[0].Desc, "Expected test type to be Some test type")
}

func TestSetThreshold(t *testing.T) {
	result := NewTestResult(1, logging.Log)
	result.SetThreshold(70)

	assert.Equal(t, 70.00, result.threshold, "Expected threshold to be 70")
}

func TestCheckThresholdPass(t *testing.T) {

	log.SetVerbosity(0)

	var success, threshold = 1.0, 2.0

	failThreshold := &TestResult{
		successRatio: success,
		threshold:    threshold,
	}

	if os.Getenv("FLAG") == "1" {
		failThreshold.checkThresholdPass()
		return
	}

	OSExitTest(t, "TestCheckThresholdPass")

}

func TestPrintSummary(t *testing.T) {
	result := NewTestResult(3, logging.Log)

	result.AddRun(&TestRun{
		TestType: "PassedTest",
		Desc:     "Passed test",
		Passed:   true,
	})

	result.AddRun(&TestRun{
		TestType: "FailedTest",
		Desc:     "Failed Test",
		Passed:   false,
	})

	result.AddRun(&TestRun{
		TestType: "InvalidTest",
		Desc:     "Invalid Test",
		Invalid:  true,
	})

	result.printSummary()

	assert.Equal(t, 1, result.passed, "Expected passed to be one.")
	assert.Equal(t, 1, result.failed, "Expected failed to be one.")
	assert.Equal(t, 1, result.invalid, "Expected invalid to be one.")
	assert.Equal(t, math.Floor(100/3.0), result.successRatio, "Expected success threshold to be 1/3.")
}

func OSExitTest(t *testing.T, testName string) {

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
