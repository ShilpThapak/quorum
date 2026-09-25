package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const maxBody = 256 << 10

func main() {
	base := os.Getenv("QUORUM_QUESTIONS")
	if base == "" {
		wd, _ := os.Getwd()
		base = filepath.Join(wd, "..", "questions-src")
	}
	g := NewGrader(base)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "GET only", http.StatusMethodNotAllowed)
			return
		}
		qs, _ := g.List()
		writeJSON(w, map[string]any{
			"status":    "ok",
			"questions": len(qs),
			"vet":       graderToolchainOK(),
		})
	})
	mux.HandleFunc("/api/questions", func(w http.ResponseWriter, r *http.Request) {
		qs, err := g.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, qs)
	})
	mux.HandleFunc("/api/source/{slug}/{kind}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "GET only", http.StatusMethodNotAllowed)
			return
		}
		slug := r.PathValue("slug")
		kind := r.PathValue("kind") // starter|reference
		if kind != "starter" && kind != "reference" {
			http.Error(w, "kind must be starter|reference", http.StatusBadRequest)
			return
		}
		qdir := filepath.Join(base, slug)
		if fi, err := os.Stat(qdir); err != nil || !fi.IsDir() {
			http.Error(w, "unknown question: "+slug, http.StatusNotFound)
			return
		}
		body, err := os.ReadFile(filepath.Join(qdir, kind+".go"))
		if err != nil {
			http.Error(w, "no "+kind+".go for "+slug, http.StatusNotFound)
			return
		}
		writeJSON(w, map[string]string{"kind": kind, "code": string(body)})
	})
	mux.HandleFunc("/api/grade", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Question string `json:"question"`
			Kind     string `json:"kind"` // reference|starter|user
			Code     string `json:"code"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBody)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), g.Timeout)
		defer cancel()
		rep, err := g.Grade(ctx, req.Question, req.Kind, []byte(req.Code))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, rep)
	})

	addr := ":" + envOr("PORT", "8090")
	log.Printf("quorum grader on http://localhost%s (questions=%s)", addr, base)
	s := &http.Server{Addr: addr, Handler: withCORS(mux), ReadTimeout: 30 * time.Second}
	log.Fatal(s.ListenAndServe())
}

// withCORS lets the separate static web site (different origin) call the API.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode: %v", err)
	}
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
