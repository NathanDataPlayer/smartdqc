package server

import (
	"dqc/internal/store"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

var rules = []Rule{}

var alerts = []Alert{
	{Level: "danger", Message: "SLA 违约：dwd.order_detail 延迟 45min", Time: "2分钟前"},
	{Level: "warn", Message: "空值率超阈：ods.user_event null_rate=2.3%", Time: "12分钟前"},
	{Level: "info", Message: "分区异常：dim.product dt 缺失", Time: "1小时前"},
}

var tables = []Table{
	{DB: "ods", Name: "user_event", Partition: "dt", RuleCount: 6, Health: "良好"},
	{DB: "dwd", Name: "order_detail", Partition: "dt", RuleCount: 8, Health: "关注"},
	{DB: "dim", Name: "product", Partition: "dt", RuleCount: 4, Health: "良好"},
}

type Rule struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Table   string `json:"table"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	LastRun string `json:"lastRun"`
}

type Alert struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

type Table struct {
	DB        string `json:"db"`
	Name      string `json:"name"`
	Partition string `json:"partition"`
	RuleCount int    `json:"ruleCount"`
	Health    string `json:"health"`
}

type Overview struct {
	RuleCount  int    `json:"ruleCount"`
	Alerts24h  int    `json:"alerts24h"`
	SLAIndex   int    `json:"slaIndex"`
	Compliance string `json:"compliance"`
}

func jsonResp(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func Start() {
	if err := store.Init(); err == nil {
		store.EnsureSeed()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) { jsonResp(w, map[string]string{"status": "ok"}) })
	mux.HandleFunc("/api/overview", func(w http.ResponseWriter, r *http.Request) {
		rs, _ := store.ListRules()
		jsonResp(w, Overview{RuleCount: len(rs), Alerts24h: 12, SLAIndex: 97, Compliance: "99.2%"})
	})
	mux.HandleFunc("/api/rules", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			rs, err := store.ListRules()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			jsonResp(w, rs)
			return
		}
		if r.Method == http.MethodPost {
			var in Rule
			json.NewDecoder(r.Body).Decode(&in)
			out, err := store.CreateRule(store.Rule{Name: in.Name, Table: in.Table, Type: in.Type, Status: "enabled", LastRun: time.Now().Format("15:04")})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			jsonResp(w, Rule{ID: out.ID, Name: out.Name, Table: out.Table, Type: out.Type, Status: out.Status, LastRun: out.LastRun})
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/api/rules/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/api/rules/"):]
		id, _ := strconv.Atoi(idStr)
		if r.Method == http.MethodPut {
			var in Rule
			json.NewDecoder(r.Body).Decode(&in)
			out, err := store.UpdateRule(id, store.Rule{Name: in.Name, Table: in.Table, Type: in.Type, Status: in.Status})
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			jsonResp(w, Rule{ID: out.ID, Name: out.Name, Table: out.Table, Type: out.Type, Status: out.Status, LastRun: out.LastRun})
			return
		}
		if r.Method == http.MethodDelete {
			if err := store.DeleteRule(id); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) { jsonResp(w, alerts) })
	mux.HandleFunc("/api/tables", func(w http.ResponseWriter, r *http.Request) { jsonResp(w, tables) })

	h := cors(mux)
	http.ListenAndServe(":8088", h)
}

func choose(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
