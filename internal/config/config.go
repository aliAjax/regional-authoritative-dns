package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr     string
	DNSAddr      string
	APIKey       string
	MaxBody      int64
	QueryTimeout time.Duration
	DatabaseURL  string
	LogLevel     string
}

func Default() Config {
	return Config{HTTPAddr: ":8080", DNSAddr: ":5353", APIKey: "dev-api-key", MaxBody: 1 << 20, QueryTimeout: 5 * time.Second, DatabaseURL: "postgres://dns:dns@localhost:5432/dns?sslmode=disable", LogLevel: "INFO"}
}

func Load(path string) Config {
	c := Default()
	if path == "" {
		path = os.Getenv("DNS_CONFIG")
	}
	if path == "" {
		path = "configs/config.yaml"
	}
	_ = loadFlatYAML(path, &c)
	if v := os.Getenv("DNS_HTTP_ADDR"); v != "" {
		c.HTTPAddr = v
	}
	if v := os.Getenv("DNS_ADDR"); v != "" {
		c.DNSAddr = v
	}
	if v := os.Getenv("DNS_API_KEY"); v != "" {
		c.APIKey = v
	}
	if v := os.Getenv("DNS_DATABASE_URL"); v != "" {
		c.DatabaseURL = v
	}
	if v := os.Getenv("DNS_MAX_BODY"); v != "" {
		if n, e := strconv.ParseInt(v, 10, 64); e == nil {
			c.MaxBody = n
		}
	}
	if v := os.Getenv("DNS_QUERY_TIMEOUT"); v != "" {
		if d, e := time.ParseDuration(v); e == nil {
			c.QueryTimeout = d
		}
	}
	return c
}

func loadFlatYAML(path string, c *Config) error {
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		p := strings.SplitN(strings.TrimSpace(s.Text()), ":", 2)
		if len(p) != 2 {
			continue
		}
		k := strings.TrimSpace(p[0])
		v := strings.Trim(strings.TrimSpace(p[1]), "\"'")
		switch k {
		case "http_addr":
			c.HTTPAddr = v
		case "dns_addr":
			c.DNSAddr = v
		case "api_key":
			c.APIKey = v
		case "database_url":
			c.DatabaseURL = v
		case "max_body":
			if n, e := strconv.ParseInt(v, 10, 64); e == nil {
				c.MaxBody = n
			}
		case "query_timeout":
			if d, e := time.ParseDuration(v); e == nil {
				c.QueryTimeout = d
			}
		}
	}
	return s.Err()
}
