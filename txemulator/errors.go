package txemulator

// ErrorWithExitCode is an error returned when emulation failed with exit code.
type ErrorWithExitCode struct {
	Message   string
	Iteration int
	ExitCode  int
}

func (e ErrorWithExitCode) Error() string {
	return e.Message
}

func (e ErrorWithExitCode) Is(err error) bool {
	_, ok := err.(ErrorWithExitCode)
	return ok
}
