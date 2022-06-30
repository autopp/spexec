package testutil

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"
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
	return m.matched, m.message, nil
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
	return func(env *model.Env, v *spec.Validator, r *matcher.StatusMatcherRegistry, x any) model.StatusMatcher {
		calls.Calls = append(calls.Calls, x)
		return NewExampleStatusMatcher(matched, message, err)
	}, calls
}

func GenFailedParseStatusMatcher(violationMessage string) matcher.StatusMatcherParser {
	return func(env *model.Env, v *spec.Validator, r *matcher.StatusMatcherRegistry, x any) model.StatusMatcher {
		v.AddViolation(violationMessage)
		return nil
	}
}

func GenParseExampleStreamMatcher(matched bool, message string, err error) (matcher.StreamMatcherParser, *ParserCalls) {
	calls := &ParserCalls{Calls: make([]any, 0)}
	return func(env *model.Env, v *spec.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
		calls.Calls = append(calls.Calls, x)
		return NewExampleStreamMatcher(matched, message, err)
	}, calls
}

func GenFailedParseStreamMatcher(violationMessage string) matcher.StreamMatcherParser {
	return func(env *model.Env, v *spec.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
		v.AddViolation(violationMessage)
		return nil
	}
}
