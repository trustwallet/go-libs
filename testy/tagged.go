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
// 		TaggedTestsEnvVar="TEST_TAGS"
//		env TEST_TAGS="unit" go test ./...
var TaggedTestsEnvVar = "TEST_TAGS"

// TaggedTest runs the test if the provided tag matches at least one runtime tag.
// Example:
// 		func TestSomething(t *testing.T) {
//    		TaggedTest(t, "unit")
//    		...
// 		}
// Run with:
// 		env TEST_TAGS="unit,integration" go test ./...
func TaggedTest(t *testing.T, testTag string) {
	tags := getRuntimeTags()
	if shouldRun := tags.empty() || tags.contains(testTag); !shouldRun {
		t.SkipNow()
	}
}

// TaggedOrTest runs the test if any of the provided test tags matches one of the runtime tags.
func TaggedOrTest(t *testing.T, testTags ...string) {
	tags := getRuntimeTags()
	if shouldRun := tags.empty() || tags.containsAny(testTags...); !shouldRun {
		t.SkipNow()
	}
}

// TaggedAndTest runs the test if all the provided test tags appear in runtime tags.
func TaggedAndTest(t *testing.T, testTags ...string) {
	tags := getRuntimeTags()
	if shouldRun := tags.empty() || tags.containsAll(testTags...); !shouldRun {
		t.SkipNow()
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

func (rt runtimeTags) empty() bool {
	return len(rt) == 0
}
