package httpapi

import "net/http"

type response struct {
	status int
	header http.Header
	body   []byte
}

func (resp response) write(w http.ResponseWriter) {
	dst := w.Header()
	for key, values := range resp.header {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
	status := resp.status
	if status == 0 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	if len(resp.body) == 0 {
		return
	}
	_, _ = w.Write(resp.body)
}

func jsonHeader() http.Header {
	h := make(http.Header)
	h.Set("Content-Type", "application/json")
	return h
}
