package kubetest

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type testResult struct {
	invalid, passed, failed, total int
	successRatio, threshold        float64
	testRuns                       [][]string
}

func (res *testResult) print() {

	testList := tablewriter.NewWriter(os.Stdout)
	testList.SetHeader([]string{"Test Type", "Spec", "Passed"})
	testList.SetBorder(false)
	testList.AppendBulk(res.testRuns)
	fmt.Println()
	testList.Render()
	fmt.Println()

	res.successRatio = res.calculateSuccessRatio()

	data := [][]string{
		{"Total", fmt.Sprintf("%d", res.total)},
		{"Passed", fmt.Sprintf("%d", res.passed)},
		{"Failed", fmt.Sprintf("%d", res.failed)},
		{"Invalid", fmt.Sprintf("%d", res.invalid)},
		{"Expected Coverage", fmt.Sprintf("%.2f", res.threshold)},
		{"Actual Coverage", fmt.Sprintf("%.2f", res.successRatio)},
	}
	testSumamry := tablewriter.NewWriter(os.Stdout)
	testSumamry.SetHeader([]string{"Tests", "Number"})
	testSumamry.SetBorder(false)
	testSumamry.AppendBulk(data)
	fmt.Println()
	testSumamry.Render()
	fmt.Println()
}

func (res *testResult) addResultToRow(row int, add string) {
	if res.testRuns == nil {
		res.testRuns = make([][]string, 10)
	}
	res.testRuns[row] = append(res.testRuns[row], add)
}

func (res *testResult) validate() {
	if res.successRatio < res.threshold {
		log.Fatal(fmt.Errorf("Expected %.2f, Got %.2f", res.threshold, res.successRatio), "Failed Test Threshold")
	}
}

func (res *testResult) calculateSuccessRatio() float64 {
	if res.total < 1 {
		log.Fatal(fmt.Errorf("Total number of tests is less than 0"), "Exiting")
	}
	return (float64(res.passed) / float64(res.total)) * 100
}
