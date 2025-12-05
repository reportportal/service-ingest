package handler

import (
	"io"
	"net/http"
	"strings"
)

// ParseMultipartForm is a middleware that parses multipart/form-data requests,
// extracts the JSON part from the "json_request_part" field, and modifies the request
// to have a JSON body for further processing.
// 32 MB is the maximum memory used to parse the form. Files larger than this will be stored in temporary files on disk.
func ParseMultipartForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if strings.HasPrefix(contentType, "multipart/form-data") {
			r.ParseMultipartForm(32 << 20)

			jsonPart := r.FormValue("json_request_part")
			r.Body = io.NopCloser(strings.NewReader(jsonPart))

			r.Header.Set("Content-Type", "application/json")
			r.ContentLength = int64(len(jsonPart))
		}

		next.ServeHTTP(w, r)
	})
}
