package mock

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func JsonModelFromFilePath(file string, intoStruct interface{}) error {
	jsonFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, &intoStruct)
	if err != nil {
		return err
	}
	return nil
}

func JsonStringFromFilePath(file string) (string, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	return string(byteValue), nil
}

func CreateMockedAPI(funcsMap map[string]func(http.ResponseWriter, *http.Request)) http.Handler {
	r := http.NewServeMux()
	for pattern, f := range funcsMap {
		r.HandleFunc(pattern, f)
	}
	return r
}
