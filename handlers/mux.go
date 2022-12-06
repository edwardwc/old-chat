package handlers

import "net/http"

func NewProxyMux() *http.ServeMux {
	mux := http.NewServeMux()
	pipelineHandler := http.HandlerFunc(HandleRequest)

	// Regular proxy flow
	mux.Handle("/", pipelineHandler)
	return mux
}
