package server

import (
    "net/http"
    "dqc/internal/meta"
    "encoding/json"
)

func registerMetaHandlers(mux *http.ServeMux) {
    mux.HandleFunc("/api/meta/databases", func(w http.ResponseWriter, r *http.Request) {
        hc, err := meta.New(); if err != nil { w.WriteHeader(http.StatusServiceUnavailable); return }
        defer hc.Close()
        dbs, err := hc.Databases(); if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
        w.Header().Set("Content-Type","application/json")
        json.NewEncoder(w).Encode(dbs)
    })
    mux.HandleFunc("/api/meta/tables", func(w http.ResponseWriter, r *http.Request) {
        db := r.URL.Query().Get("db"); if db == "" { w.WriteHeader(http.StatusBadRequest); return }
        hc, err := meta.New(); if err != nil { w.WriteHeader(http.StatusServiceUnavailable); return }
        defer hc.Close()
        ts, err := hc.Tables(db); if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
        w.Header().Set("Content-Type","application/json")
        json.NewEncoder(w).Encode(ts)
    })
    mux.HandleFunc("/api/meta/partitions", func(w http.ResponseWriter, r *http.Request) {
        db := r.URL.Query().Get("db"); table := r.URL.Query().Get("table"); if db == "" || table == "" { w.WriteHeader(http.StatusBadRequest); return }
        hc, err := meta.New(); if err != nil { w.WriteHeader(http.StatusServiceUnavailable); return }
        defer hc.Close()
        ps, err := hc.Partitions(db, table); if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
        w.Header().Set("Content-Type","application/json")
        json.NewEncoder(w).Encode(ps)
    })
}
