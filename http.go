package quaderno

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type LogLevel int

const (
	LogLevelNone LogLevel = iota
	LogLevelHeaders
	LogLevelBody
)

type httpLogger struct {
	transport http.RoundTripper
	level     LogLevel
}

func (l *httpLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	if l.level == LogLevelNone {
		return l.transport.RoundTrip(req)
	}

	start := time.Now()
	log.Printf("--> %s %s", req.Method, req.URL)

	if l.level >= LogLevelHeaders {
		for k, v := range req.Header {
			val := v[0]

			if strings.EqualFold(k, "Authorization") {
				if strings.HasPrefix(val, "Basic ") {
					val = "Basic " + strings.Repeat("*", len(val)-6)
				} else {
					val = strings.Repeat("*", len(val))
				}
			}

			log.Printf("%s: %s", k, val)
		}
	}

	if l.level >= LogLevelBody && req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err == nil && len(bodyBytes) > 0 {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			log.Printf("\n%s", string(bodyBytes))
		}
	}
	log.Printf("--> END %s\n\n", req.Method)

	res, err := l.transport.RoundTrip(req)
	if err != nil {
		log.Printf("<-- HTTP ERROR: %v", err)
		return nil, err
	}

	duration := time.Since(start)
	log.Printf("<-- %s %s (%v)", res.Status, req.URL, duration)

	if l.level >= LogLevelHeaders {
		for k, v := range res.Header {
			log.Printf("%s: %s", k, v[0])
		}
	}

	if l.level >= LogLevelBody && res.Body != nil {
		bodyBytes, ioErr := io.ReadAll(res.Body)
		if ioErr == nil && len(bodyBytes) > 0 {
			res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			log.Printf("\n%s", string(bodyBytes))
		}
	}
	log.Printf("<-- END HTTP\n\n")

	return res, nil
}
