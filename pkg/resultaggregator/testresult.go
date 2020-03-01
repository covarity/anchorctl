package resultaggregator

import (
	"anchorctl/pkg/logging"
	"fmt"
	"os"

	"github.com/kataras/tablewriter"
	"github.com/landoop/tableprinter"
)

var log *logging.Logger

func (res *TestResult) Render() {
	res.printTests()
	res.printSummary()
	res.checkThresholdPass()
}

func NewTestResult(len int, logger *logging.Logger) *TestResult {
	log = logger
	return &TestResult{
		total:    len,
		testRuns: []TestRun{},
	}
}

func NewTestRun(testType string, desc string) *TestRun {
	return &TestRun{
		TestType: testType,
		Desc:     desc,
		Passed:   false,
		Invalid:  false,
	}
}

func (res *TestResult) AddInvalidTestRun(testType string, err error) *TestRun {
	log.Warn("Error", err.Error(), "Decoding "+testType+" returned error")
	return &TestRun{
		Index:    len(res.testRuns) + 1,
		TestType: testType,
		Desc:     err.Error(),
		Passed:   false,
		Invalid:  true,
	}
}

func (res *TestResult) AddRun(testRun *TestRun) {
	testRun.Index = len(res.testRuns) + 1
	res.testRuns = append(res.testRuns, *testRun)
}

func (res *TestResult) SetThreshold(threshold float64) {
	res.threshold = threshold
}

func (res *TestResult) checkThresholdPass() {
	if res.successRatio < res.threshold {
		log.Fatal(fmt.Errorf("expected %.2f, Got %.2f", res.threshold, res.successRatio), "Failed test Threshold")
	}
}

func (res *TestResult) printTests() {
	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	printer.HeaderBgColor = tablewriter.BgBlackColor
	printer.HeaderFgColor = tablewriter.FgGreenColor

	fmt.Println("=========================================================================================================")
	fmt.Println()
	fmt.Println("Test Runs")
	tableprinter.Print(os.Stdout, res.testRuns)
	fmt.Println()
}

func (res *TestResult) printSummary() {
	for _, i := range res.testRuns {
		if i.Invalid {
			res.invalid++
		}

		if i.Passed {
			res.passed++
		} else {
			res.failed++
		}
	}

	if res.total < 1 {
		log.Fatal(fmt.Errorf("total number of tests is less than 0"), "Exiting")
	}

	res.successRatio = (float64(res.passed) / float64(res.total)) * 100

	data := [][]string{
		{"Total", fmt.Sprintf("%d", res.total)},
		{"passed", fmt.Sprintf("%d", res.passed)},
		{"Failed", fmt.Sprintf("%d", res.failed)},
		{"invalid", fmt.Sprintf("%d", res.invalid)},
		{"Expected Coverage", fmt.Sprintf("%.2f", res.threshold)},
		{"Actual Coverage", fmt.Sprintf("%.2f", res.successRatio)},
	}

	testSumamry := tablewriter.NewWriter(os.Stdout)
	testSumamry.SetHeader([]string{"Tests", "Number"})
	testSumamry.SetBorder(false)
	testSumamry.AppendBulk(data)

	fmt.Println("Test Summary")
	fmt.Println()
	testSumamry.Render()
	fmt.Println("=========================================================================================================")
}
