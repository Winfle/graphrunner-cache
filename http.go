package cacheproxy

import (
	"encoding/json"
	"net/http"
)

type HttpPayload struct {
	Method string
	Body   []byte      `json:"body"`
	Header http.Header `json:"headers"`
}

func (h *HttpPayload) GetResponseBody() string {
	if h.Header.Get("Content-Encoding") == "gzip" {
		bodyStr, _ := DecompressGzip(h.Body)
		return string(bodyStr)
	}

	return string(h.Body)
}

func (h *HttpPayload) HashKey() string {
	store := h.Header.Get("store")
	if store == "" {
		store = "default"
	}
	hash := HashBytes(h.Body)

	return store + ":" + hash
}

func (h *HttpPayload) serialize() ([]byte, error) {
	compressedBody, err := CompressGzip([]byte(h.GetResponseBody()))
	if err != nil {
		return nil, err
	}

	p := HttpPayload{
		Method: h.Method,
		Body:   compressedBody,
		Header: h.Header,
	}

	return json.Marshal(p)
}

func UnserializeHttpPayload(data []byte) (HttpPayload, error) {
	var p HttpPayload
	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		return HttpPayload{}, err
	}

	// Decompress the body
	decompressedBody, err := DecompressGzip(p.Body)
	if err != nil {
		return HttpPayload{}, err
	}
	p.Body = decompressedBody

	return p, nil
}
