import re
p = "main.go"
s = open(p).read()

needle = '''	mux.HandleFunc("POST /api/grade", func(w http.ResponseWriter, r *http.Request) {'''
repl = '''	mux.HandleFunc("POST /api/format", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "POST only", http.StatusMethodNotAllowed); return }
		var body struct{ Code string `json:"code"` }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest); return
		}
		ctx, cancel := context.WithTimeout(r.Context(), g.Timeout); defer cancel()
		formatted, err := g.Format(ctx, []byte(body.Code))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest); return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]string{"formatted": formatted})
	})

''' + needle
assert s.count(needle) == 1, s.count(needle)
open(p, "w").write(s.replace(needle, repl))
print("route added")
