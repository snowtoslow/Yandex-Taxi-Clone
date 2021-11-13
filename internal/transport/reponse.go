package transport

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func createResponse(err error, response []byte, req *http.Request) *http.Response {
	resp := &http.Response{}
	if err != nil {
		resp.Status = http.StatusText(http.StatusInternalServerError)
		resp.StatusCode = http.StatusInternalServerError
		resp.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		resp.Body = ioutil.NopCloser(ioutil.NopCloser(bytes.NewBufferString(err.Error())))
	} else {
		resp.Status = http.StatusText(http.StatusOK)
		resp.StatusCode = http.StatusOK
		resp.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		resp.Body = ioutil.NopCloser(ioutil.NopCloser(bytes.NewReader(response)))
		resp.Request = req
	}

	return resp
}

func createStatus() {

}
