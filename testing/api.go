package testing

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

const baseAPI = "http://127.0.0.1:4242"

type (
	// HandlerForTest implement the function signature used to check the req/resp
	HandlerForTest = func(t *testing.T, resp *http.Response)
)

// PushAPI is used to push a request to a local API
func PushAPI(t *testing.T, path string, content []byte) *http.Response {
	t.Helper()

	var url = baseAPI + path

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		t.Fatalf("can't post the new request : %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return resp
}

// RequestAPI is used to request the local API
func RequestAPI(t *testing.T, path string) (resp *http.Response) {
	resp, err := http.Get(baseAPI + path)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return resp
}

func PushAndTestAPI(t *testing.T, path string, content []byte, handler HandlerForTest) {
	resp := PushAPI(t, path, content)
	defer resp.Body.Close()
	handler(t, resp)
}

func RequestAndTestAPI(t *testing.T, path string, handler HandlerForTest) {
	resp := RequestAPI(t, path)
	defer resp.Body.Close()
	handler(t, resp)
}

func FetchBody(t *testing.T, resp *http.Response) (body string, err error) {
	var bbody []byte

	if bbody, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Fatalf("error fetching the body response : %s", err.Error())
	}
	defer resp.Body.Close()

	body = string(bbody)

	return
}

func AssertBody(t *testing.T, resp *http.Response, expected string) {
	t.Helper()
	body, err := FetchBody(t, resp)
	AssertNil(t, err)
	AssertStringEqual(t, body, expected)
}

func AssertBodyDiffere(t *testing.T, resp *http.Response, expected string) {
	t.Helper()
	body, err := FetchBody(t, resp)
	AssertNil(t, err)
	AssertStringNotEqual(t, body, expected)
}

func AssertStatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	AssertIntEqual(t, resp.StatusCode, expected)
}

func AssertHeader(t *testing.T, resp *http.Response, key, val string) bool {
	// test existence
	if out, ok := resp.Header[key]; !ok || len(out) == 0 || out[0] != val {
		t.Errorf("Invalid response header [%s] expected: [%s]", out[0], val)
		return false
	}

	return true
}
