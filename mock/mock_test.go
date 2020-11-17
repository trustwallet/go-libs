package mock

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateMockedAPI(t *testing.T) {
	data := make(map[string]func(http.ResponseWriter, *http.Request))
	data["/1"] = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, `{"status": true}`); err != nil {
			panic(err)
		}
	}

	server := httptest.NewServer(CreateMockedAPI(data))
	defer server.Close()

	resp, err := req.Get(server.URL + "/1")
	assert.Nil(t, err)
	type S struct {
		Status bool
	}
	var s S

	assert.Nil(t, resp.ToJSON(&s))
	assert.True(t, s.Status)
}

func TestParseJsonFromFilePath(t *testing.T) {
	type S struct {
		Status bool
	}
	var s S
	err := ParseJsonFromFilePath("test.json", &s)
	assert.Nil(t, err)
	assert.True(t, s.Status)
}

func TestJsonFromFilePathToString(t *testing.T) {
	data, err := JsonFromFilePathToString("test.json")
	assert.Nil(t, err)
	assert.Equal(t, `{
  "status": true
}`, data)
}
