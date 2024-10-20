package middleware

import (
	"bytes"
	"io"
	"net/http"
)

type httpWriter struct {
	w          http.ResponseWriter
	statusCode int
	data       *bytes.Buffer
	writer     io.Writer
}

func newHttpWriter(w http.ResponseWriter) *httpWriter {
	buf := new(bytes.Buffer)
	return &httpWriter{
		w:      w,
		data:   buf,
		writer: io.MultiWriter(w, buf),
	}
}

func (h *httpWriter) Header() http.Header {
	return h.w.Header()
}

func (h *httpWriter) Write(data []byte) (int, error) {
	return h.writer.Write(data)
}

func (h *httpWriter) WriteHeader(statusCode int) {
	h.statusCode = statusCode
	h.w.WriteHeader(statusCode)
}

func (h *httpWriter) StatusCode() int {
	if h.statusCode == 0 {
		return http.StatusOK
	}
	return h.statusCode
}
