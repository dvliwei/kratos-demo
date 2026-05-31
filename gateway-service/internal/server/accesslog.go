package server

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"gateway-service/internal/pkg/requestid"

	"github.com/go-kratos/kratos/v2/log"
)

const maxAccessLogBodyBytes = 4096

type accessLogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *accessLogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *accessLogResponseWriter) Write(body []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	if w.body.Len() < maxAccessLogBodyBytes {
		remaining := maxAccessLogBodyBytes - w.body.Len()
		if len(body) > remaining {
			w.body.Write(body[:remaining])
		} else {
			w.body.Write(body)
		}
	}
	return w.ResponseWriter.Write(body)
}

type networkRequest struct {
	Method       string      `json:"method"`
	Path         string      `json:"path"`
	Query        string      `json:"query"`
	Header       http.Header `json:"header"`
	UserToken    string      `json:"user_token"`
	RequestId    string      `json:"request_id"`
	Status       int         `json:"status"`
	CostMs       int64       `json:"cost_ms"`
	RequestBody  string      `json:"request_body"`
	ResponseBody string      `json:"response_body"`
	CreatedAt    string      `json:"created_at"`
}

func accessLogMiddleware(logger log.Logger, publisher *AccessLogPublisher) func(http.Handler) http.Handler {
	helper := log.NewHelper(logger)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestBody := readAndRestoreRequestBody(r)
			recorder := &accessLogResponseWriter{ResponseWriter: w}

			next.ServeHTTP(recorder, r)

			statusCode := recorder.statusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
			}
			cost := time.Since(start)
			accessLog := networkRequest{
				Method:       r.Method,
				Path:         r.URL.Path,
				Query:        r.URL.RawQuery,
				Header:       r.Header.Clone(),
				UserToken:    r.Header.Get(userTokenHeader),
				RequestId:    requestid.FromContext(r.Context()),
				Status:       statusCode,
				CostMs:       cost.Milliseconds(),
				RequestBody:  truncateAccessLogBody(requestBody),
				ResponseBody: truncateAccessLogBody(recorder.body.Bytes()),
				CreatedAt:    start.Format(time.RFC3339Nano),
			}
			helper.Infof(
				"access method=%s path=%s query=%s user_token=%s request_id=%s status=%d cost=%s request_body=%s response_body=%s",
				accessLog.Method,
				accessLog.Path,
				accessLog.Query,
				accessLog.UserToken,
				accessLog.RequestId,
				accessLog.Status,
				cost,
				accessLog.RequestBody,
				accessLog.ResponseBody,
			)
			publisher.Publish(accessLog)
		})
	}
}

func readAndRestoreRequestBody(r *http.Request) []byte {
	if r.Body == nil {
		return nil
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		r.Body = io.NopCloser(bytes.NewReader(nil))
		return nil
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

func truncateAccessLogBody(body []byte) string {
	if len(body) <= maxAccessLogBodyBytes {
		return string(body)
	}
	return string(body[:maxAccessLogBodyBytes]) + "...(truncated)"
}
