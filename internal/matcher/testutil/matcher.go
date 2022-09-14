package testutil

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
)

type exampleMatcher[T any] struct {
	matched bool
	message string
	err     error
}

func newExampleMatcher[T any](matched bool, message string, err error) *exampleMatcher[T] {
	return &exampleMatcher[T]{
		matched: matched,
		message: message,
		err:     err,
	}
}

func (m *exampleMatcher[T]) Match(actual T) (bool, string, error) {
	if m.err != nil {
		return false, "", m.err
	}

	var message string
	if m.matched {
		message = m.SuccessMessage()
	} else {
		message = m.FailureMessage()
	}

	return m.matched, message, nil
}

func (m *exampleMatcher[T]) SuccessMessage() string {
	return m.message + " success"
}

func (m *exampleMatcher[T]) FailureMessage() string {
	return m.message + " failure"
}

type ParserCalls struct {
	Calls []any
}

type ExampleStatusMatcher = exampleMatcher[int]

func NewExampleStatusMatcher(matched bool, message string, err error) *ExampleStatusMatcher {
	return newExampleMatcher[int](matched, message, err)
}

type ExampleStreamMatcher = exampleMatcher[[]byte]

func NewExampleStreamMatcher(matched bool, message string, err error) *ExampleStreamMatcher {
	return newExampleMatcher[[]byte](matched, message, err)
}

func GenParseExampleStatusMatcher(matched bool, message string, err error) (matcher.StatusMatcherParser, *ParserCalls) {
	calls := &ParserCalls{Calls: make([]any, 0)}
	return func(v *model.Validator, r *matcher.StatusMatcherRegistry, x any) model.StatusMatcher {
		calls.Calls = append(calls.Calls, x)
		return NewExampleStatusMatcher(matched, message, err)
	}, calls
}

func GenFailedParseStatusMatcher(violationMessage string) matcher.StatusMatcherParser {
	return func(v *model.Validator, r *matcher.StatusMatcherRegistry, x any) model.StatusMatcher {
		v.AddViolation(violationMessage)
		return nil
	}
}

func GenParseExampleStreamMatcher(matched bool, message string, err error) (matcher.StreamMatcherParser, *ParserCalls) {
	calls := &ParserCalls{Calls: make([]any, 0)}
	return func(v *model.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
		calls.Calls = append(calls.Calls, x)
		return NewExampleStreamMatcher(matched, message, err)
	}, calls
}

func GenFailedParseStreamMatcher(violationMessage string) matcher.StreamMatcherParser {
	return func(v *model.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
		v.AddViolation(violationMessage)
		return nil
	}
}
