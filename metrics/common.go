package metrics

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bitrise-io/go-utils/httputil"
)

func getContentTypeFromHeader(header http.Header) string {
	contentType, err := httputil.GetSingleValueFromHeader("Content-Type", header)
	if err != nil {
		return fmt.Sprintf("CONTENT-TYPE-ERROR:%s", err)
	}
	return contentType
}

// WrapHandlerFunc ...
func WrapHandlerFunc(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	requestWrap := func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		h(w, req)
		log.Printf(" => %s: %s: %s - %s", req.Method, getContentTypeFromHeader(req.Header), req.RequestURI, time.Since(startTime))
	}
	if newRelicAgent == nil {
		return requestWrap
	}
	return newRelicAgent.WrapHTTPHandlerFunc(requestWrap)
}

// Trace ...
func Trace(name string, fn func()) {
	wrapFn := func() {
		startTime := time.Now()
		fn()
		log.Printf(" ==> TRACE (%s) - %s", name, time.Since(startTime))
	}
	if newRelicAgent == nil {
		wrapFn()
		return
	}
	newRelicAgent.Tracer.Trace(name, wrapFn)
}
