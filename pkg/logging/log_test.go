package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestSetVerbosity(t *testing.T) {

	Log.SetVerbosity(3)

	assert.Equal(t, 3, Log.Verbosity, "Expected verbosity to  be 3")
}

func TestFatal(t *testing.T) {
	log := Logger{}
	log.SetVerbosity(0)

	if os.Getenv("FLAG") == "1" {
		log.Fatal(fmt.Errorf("Some fatal error"), "Testing os exit on fatal error")
		return
	}

	OSExitTest(t, "TestFatal")
}

func TestConvertLevel(t *testing.T) {
	result := convertLevel(5)

	assert.Equal(t, logrus.InfoLevel, result, "Expected verbosity to be info")
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