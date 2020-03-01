package resultaggregator

type TestRun struct {
	Index    int    `header:"index"`
	TestType string `header:"type"`
	Desc     string `header:"desc"`
	Passed   bool   `header:"passed"`
	Invalid  bool   `header:"invalid"`
	Details  string `header:"details"`
}

type TestResult struct {
	invalid, passed, failed, total int
	successRatio, threshold        float64
	testRuns                       []TestRun
}
