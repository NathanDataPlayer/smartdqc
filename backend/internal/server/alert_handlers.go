package server

import (
    "encoding/json"
    "net/http"
    "strconv"
    "dqc/internal/store"
)

func registerAlertHandlers(mux *http.ServeMux) {
    mux.HandleFunc("/api/alerts/persist", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            as, err := store.ListAlerts(); if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
            w.Header().Set("Content-Type","application/json")
            json.NewEncoder(w).Encode(as); return
        }
        if r.Method == http.MethodPost {
            var in struct{ Level string `json:"level"`; Message string `json:"message"` }
            json.NewDecoder(r.Body).Decode(&in)
            a, err := store.CreateAlert(in.Level, in.Message); if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
            w.Header().Set("Content-Type","application/json")
            json.NewEncoder(w).Encode(a); return
        }
        w.WriteHeader(http.StatusMethodNotAllowed)
    })
    mux.HandleFunc("/api/alerts/ack/", func(w http.ResponseWriter, r *http.Request) {
        idStr := r.URL.Path[len("/api/alerts/ack/"):]
        id, _ := strconv.Atoi(idStr)
        a, err := store.AckAlert(id); if err != nil { w.WriteHeader(http.StatusNotFound); return }
        w.Header().Set("Content-Type","application/json")
        json.NewEncoder(w).Encode(a)
    })
}
