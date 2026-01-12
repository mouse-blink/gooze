package model

// Report represents the result of testing a mutation.
type Report struct {
	MutationID string
	Killed     bool   // true if test detected the mutation (test failed)
	Output     string // test output/error message
	Error      error  // error executing test (not test failure)
}
