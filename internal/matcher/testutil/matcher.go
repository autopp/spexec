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

type ExampleStreamMatcher = exampleMatcher[[]byte]

func NewExampleStreamMatcher(matched bool, message string, err error) *ExampleStreamMatcher {
	return newExampleMatcher[[]byte](matched, message, err)
}

func GenParseExampleStatusMatcher(matched bool, message string, err error) matcher.StatusMatcherParser {
	return func(env *model.Env, v *spec.Validator, r *matcher.StatusMatcherRegistry, x any) model.StatusMatcher {
		return NewExampleStatusMatcher(matched, message, err)
	}
}

func GenFailedParseStatusMatcher(violationMessage string) matcher.StatusMatcherParser {
	return func(env *model.Env, v *spec.Validator, r *matcher.StatusMatcherRegistry, x any) model.StatusMatcher {
		v.AddViolation(violationMessage)
		return nil
	}
}

func GenParseExampleStreamMatcher(matched bool, message string, err error) matcher.StreamMatcherParser {
	return func(env *model.Env, v *spec.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
		return NewExampleStreamMatcher(matched, message, err)
	}
}

func GenFailedParseStreamMatcher(violationMessage string) matcher.StreamMatcherParser {
	return func(env *model.Env, v *spec.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
		v.AddViolation(violationMessage)
		return nil
	}
}
