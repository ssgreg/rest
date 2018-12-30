package rest

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

// SendError sends msg and status code
func SendError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

// SendErrorJSON makes {error: blah, details: blah} json body and responds with error code
func SendErrorJSON(w http.ResponseWriter, r *http.Request, code int, err error, details string) {
	log.Printf("[DEBUG] %s", errDetailsMsg(r, code, err, details))
	w.WriteHeader(code)
	RenderJSON(w, r, map[string]interface{}{"error": err.Error(), "details": details})
}

func errDetailsMsg(r *http.Request, code int, err error, details string) string {

	q := r.URL.String()
	if qun, e := url.QueryUnescape(q); e == nil {
		q = qun
	}

	srcFileInfo := ""
	if pc, file, line, ok := runtime.Caller(2); ok {
		fnameElems := strings.Split(file, "/")
		funcNameElems := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		srcFileInfo = fmt.Sprintf(" [caused by %s:%d %s]", strings.Join(fnameElems[len(fnameElems)-3:], "/"),
			line, funcNameElems[len(funcNameElems)-1])
	}

	remoteIP := r.RemoteAddr
	if pos := strings.Index(remoteIP, ":"); pos >= 0 {
		remoteIP = remoteIP[:pos]
	}
	return fmt.Sprintf("%s - %v - %d - %s - %s%s", details, err, code, remoteIP, q, srcFileInfo)
}
