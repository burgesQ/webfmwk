package webfmwk

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_testPort = ":6666"
	_testAddr = "http://127.0.0.1" + _testPort
)

type (
	// HandlerForTest implement the function signature used to check the req/resp
	HandlerForTest = func(t *testing.T, resp *http.Response)
)

//
// helpers methods
//

func stopServer(t *testing.T, s *Server) {
	var ctx = s.GetContext()
	t.Log("closing server ...")
	ctx.Done()
	s.Shutdown()
	s.WaitAndStop()
	Shutdown()
	t.Log("server closed")
}

func wrapperPost(t *testing.T, route, routeReq string,
	content []byte,
	handlerRoute HandlerFunc, handlerTest HandlerForTest) {
	t.Helper()

	var s = InitServer(CheckIsUp(), DisableKeepAlive())

	t.Log("init server...")
	defer stopServer(t, s)

	s.POST(route, handlerRoute)
	go s.Start(_testPort)
	<-s.isReady
	t.Log("server inited")

	pushAndTestAPI(t, _testAddr+routeReq, content, handlerTest)
}

func wrapperGet(t *testing.T, route, routeReq string,
	handlerRoute HandlerFunc, handlerTest HandlerForTest) {
	t.Helper()

	var s = InitServer(CheckIsUp(), DisableKeepAlive())

	t.Log("init server...")
	defer stopServer(t, s)

	s.GET(route, handlerRoute)
	go s.Start(_testPort)
	<-s.isReady
	t.Log("server inited")

	requestAndTestAPI(t, _testAddr+routeReq, handlerTest)
}

//
// CRUD
//

func DeleteAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()
	deleteAndTestAPI(t, url, handler)
}

func deleteAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()

	var resp = deleteAPI(t, url)
	defer resp.Body.Close()
	handler(t, resp)
}

func deleteAPI(t *testing.T, url string) *http.Response {
	t.Helper()
	// create request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("error creating the  api request : %s", err.Error())
		return nil
	}

	// fetch the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
		return nil
	}

	return resp
}

// requestAndTestAPI request an API then run the test handler
func RequestAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()
	requestAndTestAPI(t, url, handler)
}

func requestAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()
	var resp = requestAPI(t, url)
	defer resp.Body.Close()
	handler(t, resp)
}

func requestAPI(t *testing.T, url string) *http.Response {
	t.Helper()
	var resp, err = http.Get(url)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return resp
}

// pushAndTestAPI post to an API then run the test handler
// The sub method try to send an `application/json` encoded content
func pushAndTestAPI(t *testing.T, path string, content []byte, handler HandlerForTest, headers ...Header) {
	t.Helper()
	var resp = pushAPI(t, path, content, headers...)
	defer resp.Body.Close()
	handler(t, resp)
}

func pushAPI(t *testing.T, url string, content []byte, h ...Header) *http.Response {
	t.Helper()
	var req, err = http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		t.Fatalf("can't post the new request : %s", err.Error())
	}

	if len(h) == 0 {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for i := range h {
			req.Header.Set(h[i][0], h[i][1])
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return resp
}

//
// assertion methods
//

// assertBody fetch and assert that the body of the http.Response is the same than expected
func assertBody(t *testing.T, expected string, resp *http.Response) {
	t.Helper()
	assert.Equal(t, expected, fetchBody(t, resp))
}

// BodyDiffere fetch and assert that the body of the http.Response differ than expected
func assertBodyDiffere(t *testing.T, expected string, resp *http.Response) {
	t.Helper()
	assert.NotEqual(t, expected, fetchBody(t, resp))
}

func assertBodyContain(t *testing.T, expected string, resp *http.Response) {
	t.Helper()
	assert.Contains(t, fetchBody(t, resp), expected, "asserting response contains %q", expected)
}

// assertStatusCode assert the status code of the response
func assertStatusCode(t *testing.T, expected int, resp *http.Response) {
	t.Helper()

	assert.Equal(t, expected, resp.StatusCode)
}

// Header assert value of the given header key:vak in the htt.Response param
func assertHeader(t *testing.T, key, val string, resp *http.Response) {
	t.Helper()

	// test existence
	//"request header doesn't contains "+key, label...

	assert.Contains(t, resp.Header, key, "asserting header have a %q", key)
	assert.Equal(t, val, resp.Header[key][0], "asserting header value for %q", key)
}

// Header assert value of the given header key:vak in the htt.Response param
func assertHeaders(t *testing.T, resp *http.Response, kv ...[2]string) {
	t.Helper()
	for i := range kv {
		k, v := kv[i][0], kv[i][1]
		assertHeader(t, k, v, resp)
	}
}

func fetchBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	//	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("reading body: " + err.Error())
	}

	return string(bodyBytes)
}
