package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendErrorJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			t.Log("http err request", r.URL)
			SendErrorJSON(w, r, 500, errors.New("error 500"), "error details 123456")
			return
		}
		w.WriteHeader(404)
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/error")
	require.Nil(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	assert.Equal(t, `{"details":"error details 123456","error":"error 500"}`+"\n", string(body))
}

func TestErrorDetailsMsg(t *testing.T) {
	callerFn := func() {
		req, err := http.NewRequest("GET", "https://example.com/test?k1=v1&k2=v2", nil)
		require.Nil(t, err)
		req.RemoteAddr = "1.2.3.4"
		msg := errDetailsMsg(req, 500, errors.New("error 500"), "error details 123456")
		assert.Equal(t, "error details 123456 - error 500 - 500 - 1.2.3.4 - https://example.com/test?k1=v1&k2=v2 [caused by go-pkgz/rest/httperrors_test.go:44 rest.TestErrorDetailsMsg]", msg)
	}
	callerFn()
}

func TestErrorDetailsMsgWithUser(t *testing.T) {
	callerFn := func() {
		req, err := http.NewRequest("GET", "https://example.com/test?k1=v1&k2=v2", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		require.Nil(t, err)
		msg := errDetailsMsg(req, 500, errors.New("error 500"), "error details 123456")
		assert.Equal(t, "error details 123456 - error 500 - 500 - 127.0.0.1 - https://example.com/test?k1=v1&k2=v2 [caused by go-pkgz/rest/httperrors_test.go:55 rest.TestErrorDetailsMsgWithUser]", msg)
	}
	callerFn()
}

func TestSendError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			t.Log("http err request", r.URL)
			SendError(r, w, 500, "error 123456", false)
			return
		}
		w.WriteHeader(404)
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/error")
	require.Nil(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	assert.Equal(t, "error 123456", string(body))
}
