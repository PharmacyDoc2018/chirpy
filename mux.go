package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const filepathRoot = "."
const port = "8080"

func initMux(cfg *apiConfig) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("/app/assets/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})

	mux.HandleFunc("GET /admin/metrics/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl := template.Must(template.ParseFiles("./admin/metrics/index.html"))
		err := tmpl.Execute(w, struct {
			Hits int32
		}{
			Hits: cfg.fileserverHits.Load(),
		})
		if err != nil {
			fmt.Println(err)
		}
	})

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, req *http.Request) {
		if cfg.platfrom != "dev" {
			w.WriteHeader(403)
			w.Write([]byte("403 Forbidden"))
			return
		}

		err := cfg.db.ResetUsers(req.Context())
		if err != nil {
			w.WriteHeader(500)
			return
		}

		cfg.fileserverHits.Store(0)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("hit count and users reset!"))
	})

	handleLogin(mux, cfg)
	handleResourseChirps(mux, cfg)
	handleResourceUsers(mux, cfg)
	handleWebhooks(mux, cfg)

	return mux
}
