package mock

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trustwallet/go-libs/client"
)

type response struct {
	Status bool
}

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
	cli := client.InitClient(server.URL, nil)

	var resp response
	_, err := cli.Execute(context.TODO(), client.NewReqBuilder().
		Method(http.MethodGet).
		PathStatic("1").
		WriteTo(&resp).
		Build())

	assert.Nil(t, err)
	assert.True(t, resp.Status)
}

func TestParseJsonFromFilePath(t *testing.T) {
	var s response
	err := JsonModelFromFilePath("test.json", &s)

	assert.Nil(t, err)
	assert.True(t, s.Status)
}

func TestJsonStringFromFilePath(t *testing.T) {
	data, err := JsonStringFromFilePath("test.json")
	assert.Nil(t, err)
	assert.Equal(t, `{
  "status": true
}`, data)
}
