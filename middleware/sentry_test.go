package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSentryConditionAnd(t *testing.T) {
	tests := []struct {
		name       string
		conditions []SentryCondition
		expected   bool
	}{
		{
			name: "all conditions satisfied",
			conditions: []SentryCondition{
				func(res *http.Response, url string) bool {
					return true
				},
				func(res *http.Response, url string) bool {
					return true
				},
			},
			expected: true,
		},
		{
			name: "all conditions unsatisfied",
			conditions: []SentryCondition{
				func(res *http.Response, url string) bool {
					return false
				},
				func(res *http.Response, url string) bool {
					return false
				},
			},
			expected: false,
		},
		{
			name: "first of two conditions is satisfied",
			conditions: []SentryCondition{
				func(res *http.Response, url string) bool {
					return true
				},
				func(res *http.Response, url string) bool {
					return false
				},
			},
			expected: false,
		},
		{
			name: "second of two conditions is satisfied",
			conditions: []SentryCondition{
				func(res *http.Response, url string) bool {
					return false
				},
				func(res *http.Response, url string) bool {
					return true
				},
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			condition := SentryConditionAnd(tc.conditions...)
			actual := condition(nil, "")

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSentryConditionOr(t *testing.T) {
	tests := []struct {
		name       string
		conditions []SentryCondition
		expected   bool
	}{
		{
			name: "two out of two are satisfied",
			conditions: []SentryCondition{
				func(res *http.Response, url string) bool {
					return true
				},
				func(res *http.Response, url string) bool {
					return true
				},
			},
			expected: true,
		},
		{
			name: "first out of two is satisfied",
			conditions: []SentryCondition{
				func(res *http.Response, url string) bool {
					return true
				},
				func(res *http.Response, url string) bool {
					return false
				},
			},
			expected: true,
		},
		{
			name: "second out of two is satisfied",
			conditions: []SentryCondition{
				func(res *http.Response, url string) bool {
					return false
				},
				func(res *http.Response, url string) bool {
					return true
				},
			},
			expected: true,
		},
		{
			name: "none of conditions is satisfied",
			conditions: []SentryCondition{
				func(res *http.Response, url string) bool {
					return false
				},
				func(res *http.Response, url string) bool {
					return false
				},
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			condition := SentryConditionOr(tc.conditions...)
			actual := condition(nil, "")

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSentryConditionNotStatusOk(t *testing.T) {
	tests := []struct {
		name     string
		resp     *http.Response
		expected bool
	}{
		{
			name:     "response code is below 200",
			resp:     &http.Response{StatusCode: 100},
			expected: true,
		},
		{
			name:     "response code is 200",
			resp:     &http.Response{StatusCode: 200},
			expected: false,
		},
		{
			name:     "response code is between 200 and 300",
			resp:     &http.Response{StatusCode: 201},
			expected: false,
		},
		{
			name:     "response code is between 300",
			resp:     &http.Response{StatusCode: 300},
			expected: true,
		},
		{
			name:     "response code is above 300",
			resp:     &http.Response{StatusCode: 303},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := SentryConditionNotStatusOk(tc.resp, "")
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSentryConditionNotStatusBadRequest(t *testing.T) {
	tests := []struct {
		name     string
		resp     *http.Response
		expected bool
	}{
		{
			name:     "response code is bad request",
			resp:     &http.Response{StatusCode: http.StatusBadRequest},
			expected: false,
		},
		{
			name:     "response code is not bad request",
			resp:     &http.Response{StatusCode: http.StatusOK},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := SentryConditionNotStatusBadRequest(tc.resp, "")
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSentryConditionNotStatusNotFound(t *testing.T) {
	tests := []struct {
		name     string
		resp     *http.Response
		expected bool
	}{
		{
			name:     "response code is not found",
			resp:     &http.Response{StatusCode: http.StatusNotFound},
			expected: false,
		},
		{
			name:     "response code is not \"not found\"",
			resp:     &http.Response{StatusCode: http.StatusBadRequest},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := SentryConditionNotStatusNotFound(tc.resp, "")
			assert.Equal(t, tc.expected, actual)
		})
	}
}
