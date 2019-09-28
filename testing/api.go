package testing

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func PushAPI(t *testing.T, path string, content []byte) (resp *http.Response, err error) {
	url := "http://127.0.0.1:4242" + path

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if resp, err = client.Do(req); err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return
}

func RequestAPI(t *testing.T, path string) (resp *http.Response, err error) {
	if resp, err = http.Get("http://127.0.0.1:4242" + path); err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}
	return
}

func PushAndTestAPI(t *testing.T, path string, content []byte, handler func(t *testing.T, resp *http.Response) bool) bool {
	resp, err := PushAPI(t, path, content)
	if err != nil || !handler(t, resp) {
		return false
	}
	return true
}

func RequestAndTestAPI(t *testing.T, path string, handler func(t *testing.T, resp *http.Response) bool) bool {
	resp, err := RequestAPI(t, path)
	if err != nil || !handler(t, resp) {
		return false
	}
	return true
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

func TestBody(t *testing.T, resp *http.Response, expected string) bool {
	if body, err := FetchBody(t, resp); err == nil {
		if body != expected {
			t.Fatalf("error while comparing the body [%s] expected: [%s]", body, expected)
			return false
		}
		return true
	}
	return true
}

func TestBodyDiffere(t *testing.T, resp *http.Response, expected string) bool {
	if body, err := FetchBody(t, resp); err == nil {
		if body == expected {
			t.Fatalf("error while comparing the body [%s] should be as expected: [%s]", body, expected)
			return false
		}
		return true
	}
	return true
}

func TestStatusCode(t *testing.T, resp *http.Response, expected int) bool {
	if resp.StatusCode != expected {
		t.Errorf("Invalide response status code : [%d] expected: [%d]", resp.StatusCode, expected)
		return false
	}
	return true
}

func TestHeader(t *testing.T, resp *http.Response, key, val string) bool {
	// test existance
	if out, ok := resp.Header[key]; !ok || len(out) == 0 || out[0] != val {
		t.Errorf("Invalid response header [%s] expected: [%s]", out[0], val)
		return false
	}
	return true
}
