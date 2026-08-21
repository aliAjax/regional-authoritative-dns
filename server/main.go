package main

import (
	"context"
	"errors"
	"github.com/example/regional-authoritative-dns/internal/config"
	dnssecapp "github.com/example/regional-authoritative-dns/internal/dnssec/application"
	"github.com/example/regional-authoritative-dns/internal/dnssec/infrastructure"
	resolverapp "github.com/example/regional-authoritative-dns/internal/resolver/application"
	"github.com/example/regional-authoritative-dns/internal/transport"
	zoneapp "github.com/example/regional-authoritative-dns/internal/zone/application"
	zoneinfra "github.com/example/regional-authoritative-dns/internal/zone/infrastructure"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := config.Load("")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	zones := zoneinfra.New()
	records := zoneinfra.NewRecordMemory()
	zs := zoneapp.New(zones, records)
	resolver := resolverapp.New(zs)
	keys := infrastructure.New()
	_ = dnssecapp.New(keys)
	h := transport.NewHTTP(zs, resolver, logger, c.APIKey, c.MaxBody)
	httpServer := &http.Server{Addr: c.HTTPAddr, Handler: h, ReadHeaderTimeout: 3 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	dnsServer := transport.NewDNS(c.DNSAddr, resolver, logger)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if e := dnsServer.Start(ctx); e != nil {
		logger.Error("dns start failed", "error", e)
		os.Exit(1)
	}
	go func() {
		logger.Info("http listening", "addr", c.HTTPAddr)
		if e := httpServer.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			logger.Error("http failed", "error", e)
		}
	}()
	logger.Info("authoritative dns started", "dns_addr", c.DNSAddr)
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = httpServer.Shutdown(shutdownCtx)
	_ = dnsServer.Close()
	logger.Info("server stopped")
}
