package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRequest_Get(t *testing.T) {
	const aBaseURL = "http://www.example.com"

	var responses = []string{
		`{"status": "success"}`,
		`{"status": "success with data"}`,
	}

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Data(http.StatusOK, gin.MIMEJSON, []byte(responses[0]))
	})

	router.GET("/path/with/query", func(c *gin.Context) {
		queryData := c.Query("data")
		if queryData != "testdata" {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("ooops"))
			return
		}
		c.Data(http.StatusOK, gin.MIMEJSON, []byte(responses[1]))
	})

	httpClient := httpClientFromGinEngine(t, router, aBaseURL)
	c := InitClient(aBaseURL, nil, WithHttpClient(httpClient))

	tests := []struct {
		name         string
		path         string
		query        url.Values
		expectedResp string
		assertError  require.ErrorAssertionFunc
	}{
		{
			name:         "happy path simple",
			path:         "/test",
			query:        nil,
			expectedResp: responses[0],
			assertError:  require.NoError,
		},
		{
			name: "happy path with query string",
			path: "/path/with/query",
			query: func() url.Values {
				v := url.Values{}
				v.Set("data", "testdata")
				return v
			}(),
			expectedResp: responses[1],
			assertError:  require.NoError,
		},
		{
			name: "error path",
			path: "/path/with/query",
			query: func() url.Values {
				v := url.Values{}
				v.Set("data", "wrong_value")
				return v
			}(),
			expectedResp: "{}",
			assertError:  require.Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("GetWithContext", func(t *testing.T) {
				resObj := map[string]string{}
				err := c.GetWithContext(context.Background(), &resObj, test.path, test.query)
				test.assertError(t, err)

				actualRespStr, err := json.Marshal(resObj)
				require.NoError(t, err)
				require.JSONEq(t, test.expectedResp, string(actualRespStr))
			})

			t.Run("GetRaw", func(t *testing.T) {
				bytes, err := c.GetRaw(test.path, test.query)
				test.assertError(t, err)

				if string(bytes) == "" {
					bytes = []byte("{}")
				}
				require.JSONEq(t, test.expectedResp, string(bytes))
			})
		})
	}
}

func TestRequest_Post(t *testing.T) {
	const aBaseURL = "http://www.example.com"

	type reqStruct struct {
		Data string `json:"data"`
	}
	var responses = []string{
		`{"status": "success"}`,
		`{"status": "success with request"}`,
	}

	router := gin.New()
	router.POST("/test", func(c *gin.Context) {
		c.Data(http.StatusOK, gin.MIMEJSON, []byte(responses[0]))
	})

	router.POST("/a/very/long/path", func(c *gin.Context) {
		var req reqStruct
		_ = c.Bind(&req)
		if req.Data != "testdata" {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("ooops"))
			return
		}
		c.Data(http.StatusOK, gin.MIMEJSON, []byte(responses[1]))
	})

	httpClient := httpClientFromGinEngine(t, router, aBaseURL)
	c := InitJSONClient(aBaseURL, nil, WithHttpClient(httpClient))

	tests := []struct {
		name         string
		path         string
		body         any
		expectedResp string
		assertError  require.ErrorAssertionFunc
	}{
		{
			name:         "happy path no request",
			path:         "/test",
			expectedResp: responses[0],
			body:         nil,
			assertError:  require.NoError,
		},
		{
			name:         "happy path - long path, with request",
			path:         "/a/very/long/path",
			expectedResp: responses[1],
			body:         reqStruct{Data: "testdata"},
			assertError:  require.NoError,
		},
		{
			name:         "error path",
			path:         "/path/with/query",
			body:         reqStruct{Data: "wrong_data"},
			expectedResp: "{}",
			assertError:  require.Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("PostWithContext", func(t *testing.T) {
				resObj := map[string]string{}
				err := c.PostWithContext(context.Background(), &resObj, test.path, test.body)
				test.assertError(t, err)

				actualRespStr, err := json.Marshal(resObj)
				require.NoError(t, err)
				require.JSONEq(t, test.expectedResp, string(actualRespStr))
			})

			t.Run("PostRaw", func(t *testing.T) {
				bytes, err := c.PostRaw(test.path, test.body)
				test.assertError(t, err)

				if string(bytes) == "" {
					bytes = []byte("{}")
				}
				require.JSONEq(t, test.expectedResp, string(bytes))
			})
		})
	}
}

func httpClientFromGinEngine(t *testing.T, engine *gin.Engine, baseURL string) *http.Client {
	return &http.Client{
		Transport: RoundTripperFunc(func(request *http.Request) (*http.Response, error) {
			require.Equal(t, baseURL, fmt.Sprintf("%s://%s", request.URL.Scheme, request.URL.Host))

			w := httptest.NewRecorder()
			engine.ServeHTTP(w, request)
			res := w.Result()
			res.Request = request
			return res, nil
		}),
	}
}
