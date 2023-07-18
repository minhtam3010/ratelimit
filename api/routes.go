package api

import "net/http"

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/api/v1/health", healthHandler)

	return router
}
