package testy

import (
	"os"
	"strings"
	"testing"
)

const (
	TagUnit        = "unit"
	TagIntegration = "integration"
	TagPostgres    = "postgres"
	TagRabbit      = "rabbit"
)

// TaggedTestsEnvVar defines the name of the environment variable for the tagged tests.
// Example:
//
//	TaggedTestsEnvVar="TEST_TAGS"
//	env TEST_TAGS="unit" go test ./...
var TaggedTestsEnvVar = "TEST_TAGS"

// RequireTestTag runs the test if the provided tag matches at least one runtime tag.
// Example:
//
//			func TestSomething(t *testing.T) {
//	   		RequireTestTag(t, "unit")
//	   		...
//			}
//
// Run with:
//
//	env TEST_TAGS="unit,integration" go test ./...
func RequireTestTag(t *testing.T, testTag string) {
	if !getRuntimeTags().contains(testTag) {
		t.Skipf("skipping test '%s', requires '%s' tag", t.Name(), testTag)
	}
}

// RequireOneOfTestTags runs the test if any of the provided test tags matches one of the runtime tags.
func RequireOneOfTestTags(t *testing.T, testTags ...string) {
	if !getRuntimeTags().containsAny(testTags...) {
		t.Skipf("skipping test '%s', requires at least one of the following tags: '%s'",
			t.Name(), strings.Join(testTags, ", "))
	}
}

// RequireAllTestTags runs the test if all the provided test tags appear in runtime tags.
func RequireAllTestTags(t *testing.T, testTags ...string) {
	if !getRuntimeTags().containsAll(testTags...) {
		t.Skipf("skipping test '%s', requires all of the following tags: '%s'",
			t.Name(), strings.Join(testTags, ", "))
	}
}

type runtimeTags []string

func getRuntimeTags() runtimeTags {
	return parseTags(os.Getenv(TaggedTestsEnvVar))
}

func parseTags(rawTags string) runtimeTags {
	rawTags = strings.ReplaceAll(rawTags, " ", "")
	if rawTags == "" {
		return nil
	}
	return strings.Split(rawTags, ",")
}

func (rt runtimeTags) contains(targetTag string) bool {
	for _, tag := range rt {
		if tag == targetTag {
			return true
		}
	}
	return false
}

func (rt runtimeTags) containsAny(targetTags ...string) bool {
	for _, targetTag := range targetTags {
		if rt.contains(targetTag) {
			return true
		}
	}
	return false
}

func (rt runtimeTags) containsAll(targetTags ...string) bool {
	for _, targetTag := range targetTags {
		if !rt.contains(targetTag) {
			return false
		}
	}
	return true
}
