package transport

import (
	"encoding/json"
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	resolver "github.com/example/regional-authoritative-dns/internal/resolver/application"
	resdomain "github.com/example/regional-authoritative-dns/internal/resolver/domain"
	zone "github.com/example/regional-authoritative-dns/internal/zone/application"
	zoned "github.com/example/regional-authoritative-dns/internal/zone/domain"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type HTTP struct {
	zones   *zone.Service
	dns     *resolver.Service
	logger  *slog.Logger
	apiKey  string
	maxBody int64
}

func NewHTTP(z *zone.Service, r *resolver.Service, l *slog.Logger, key string, max int64) *HTTP {
	return &HTTP{zones: z, dns: r, logger: l, apiKey: key, maxBody: max}
}
func (h *HTTP) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { write(w, 200, map[string]any{"status": "ok"}) })
	m.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { write(w, 200, map[string]any{"ready": true}) })
	m.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("dns_queries_total 0\ndns_publish_total 0\n"))
	})
	m.HandleFunc("/v1/zones", h.zonesRoot)
	m.HandleFunc("/v1/zones/", h.zonePath)
	m.HandleFunc("/v1/dnssec/verify", h.verify)
	m.HandleFunc("/v1/sync-runs", func(w http.ResponseWriter, r *http.Request) { write(w, 200, []any{}) })
	return recoverer(Chain(m, h.logger, h.maxBody, 5*time.Second, h.apiKey))
}
func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.Handler().ServeHTTP(w, r) }
func (h *HTTP) zonesRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		z, e := h.zones.List(r.Context())
		if e != nil {
			jsonError(w, 500, e.Error())
			return
		}
		write(w, 200, z)
		return
	}
	if r.Method != "POST" {
		jsonError(w, 405, "method not allowed")
		return
	}
	var in struct {
		Name, Description string
		DefaultTTL        uint32
		Views             []string
		Records           domain.Set
	}
	if e := json.NewDecoder(r.Body).Decode(&in); e != nil {
		jsonError(w, 400, e.Error())
		return
	}
	z, e := h.zones.Create(r.Context(), zoned.Zone{Name: in.Name, Description: in.Description, DefaultTTL: in.DefaultTTL, Views: in.Views}, in.Records)
	if e != nil {
		jsonError(w, 400, e.Error())
		return
	}
	write(w, 201, z)
}
func (h *HTTP) zonePath(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(p) < 3 {
		jsonError(w, 404, "zone id required")
		return
	}
	id := p[2]
	switch {
	case len(p) == 4 && p[3] == "records":
		z, rs, e := h.zones.Get(r.Context(), id)
		if e != nil {
			jsonError(w, 404, e.Error())
			return
		}
		write(w, 200, map[string]any{"zone": z, "records": rs})
	case len(p) == 4 && p[3] == "publish":
		if r.Method != "POST" {
			jsonError(w, 405, "method not allowed")
			return
		}
		z, e := h.zones.Publish(r.Context(), id)
		if e != nil {
			jsonError(w, 400, e.Error())
			return
		}
		write(w, 200, z)
	case len(p) == 4 && p[3] == "rollback":
		var in struct{ Serial uint32 }
		_ = json.NewDecoder(r.Body).Decode(&in)
		z, e := h.zones.Rollback(r.Context(), id, in.Serial)
		if e != nil {
			jsonError(w, 400, e.Error())
			return
		}
		write(w, 200, z)
	case len(p) == 4 && p[3] == "freeze":
		z, e := h.zones.Freeze(r.Context(), id)
		if e != nil {
			jsonError(w, 400, e.Error())
			return
		}
		write(w, 200, z)
	case len(p) == 4 && p[3] == "versions":
		var in struct {
			Records  domain.Set
			Expected int64
		}
		if e := json.NewDecoder(r.Body).Decode(&in); e != nil {
			jsonError(w, 400, e.Error())
			return
		}
		z, e := h.zones.AddVersion(r.Context(), id, in.Records, in.Expected)
		if e != nil {
			jsonError(w, 409, e.Error())
			return
		}
		write(w, 200, z)
	default:
		jsonError(w, 404, "not found")
	}
}
func (h *HTTP) verify(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ZoneID string
		Name   string
		Type   domain.Type
		View   string
	}
	if e := json.NewDecoder(r.Body).Decode(&in); e != nil {
		jsonError(w, 400, e.Error())
		return
	}
	if in.Type == "" {
		in.Type = domain.A
	}
	res, e := h.dns.Resolve(r.Context(), in.ZoneID, resdomain.Query{Name: in.Name, Type: in.Type, View: in.View})
	if e != nil {
		jsonError(w, 400, e.Error())
		return
	}
	write(w, 200, map[string]any{"valid": !res.ServFail, "result": res})
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
