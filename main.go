package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	fileServer := http.FileServer(http.Dir("."))
	handler := http.StripPrefix("/app", fileServer)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)

		htmlTemplate := `<html>
	  <body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	  </body>
	</html>`

		hits := apiCfg.fileserverHits.Load()
		html := fmt.Sprintf(htmlTemplate, hits)
		w.Write([]byte(html))
	})
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	fmt.Println("Starting server at http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop the server")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
