package kubetest

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

type testResult struct {
	invalid, passed, failed, total int
	successRatio, threshold        float64
}

func (res *testResult) print() {
	data := [][]string{
		{"Total", fmt.Sprintf("%d", res.total)},
		{"Passed", fmt.Sprintf("%d", res.passed)},
		{"Failed", fmt.Sprintf("%d", res.failed)},
		{"Invalid", fmt.Sprintf("%d", res.invalid)},
		{"Expected Coverage", fmt.Sprintf("%.2f", res.threshold)},
		{"Actual Coverage", fmt.Sprintf("%.2f", res.successRatio)},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tests", "Number"})
	table.SetBorder(false)
	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

func (res *testResult) validate() {
	if res.total < 1{
		log.Fatal(fmt.Errorf("Total number of tests is less than 0"), "Exiting")
	}

	res.successRatio = (float64(res.passed) / float64(res.total)) * 100

	if res.successRatio < res.threshold {
		log.Fatal(fmt.Errorf("Expected %.2f, Got %.2f", res.threshold, res.successRatio), "Failed Test Threshold")
	}
}
