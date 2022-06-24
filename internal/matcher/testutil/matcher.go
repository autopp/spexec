package testutil

type exampleMatcher[T any] struct {
	matched bool
	message string
	err     error
	calls   int
}

func newExampleMatcher[T any](matched bool, message string, err error) *exampleMatcher[T] {
	return &exampleMatcher[T]{
		matched: matched,
		message: message,
		err:     err,
		calls:   0,
	}
}

func (m *exampleMatcher[T]) Match(actual T) (bool, string, error) {
	m.calls++
	if m.err != nil {
		return false, "", m.err
	}
	return m.matched, m.message, nil
}

func (m *exampleMatcher[T]) Calls() int {
	return m.calls
}

type ExampleStatusMatcher = exampleMatcher[int]

func NewExampleStatusMatcher(matched bool, message string, err error) *ExampleStatusMatcher {
	return newExampleMatcher[int](matched, message, err)
}

type ExampleStreamMatcher = exampleMatcher[string]

func NewExampleStreamMatcher(matched bool, message string, err error) *ExampleStreamMatcher {
	return newExampleMatcher[string](matched, message, err)
}
