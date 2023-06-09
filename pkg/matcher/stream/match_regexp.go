package stream

import (
	"fmt"
	"regexp"

	"github.com/autopp/spexec/pkg/matcher"
	"github.com/autopp/spexec/pkg/model"
)

type MatchRegexpMatcher struct {
	expected *regexp.Regexp
}

func (m *MatchRegexpMatcher) Match(actual []byte) (bool, string, error) {
	if m.expected.Match(actual) {
		return true, fmt.Sprintf("should not match to %q, but match", m.expected.String()), nil
	}

	return false, fmt.Sprintf("should match to %q, but got %q", m.expected.String(), string(actual)), nil
}

func ParseMatchRegexpMatcher(v *model.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
	pattern, ok := v.MustBeString(x)
	if !ok {
		return nil
	}

	expected, err := regexp.Compile(pattern)
	if err != nil {
		v.AddViolation("cannot parse regexp %q: %s", pattern, err)
		return nil
	}

	return &MatchRegexpMatcher{expected: expected}
}
