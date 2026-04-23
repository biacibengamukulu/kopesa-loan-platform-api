package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port              string
	JWTSecret         string
	CassandraHosts    []string
	CassandraPort     int
	CassandraUser     string
	CassandraPass     string
	CassandraKeyspace string
	CassandraDC       string
	ProviderBaseURL   string
	DropboxBaseURL    string
	WhatsAppBaseURL   string
	WhatsAppAPIKey    string
	WhatsAppInstance  string
	EmailFrom         string
}

func Load() Config {
	return Config{
		Port:              env("PORT", "8080"),
		JWTSecret:         env("JWT_SECRET", "kopesa-dev-secret"),
		CassandraHosts:    split(env("CASSANDRA_HOSTS", "127.0.0.1")),
		CassandraPort:     envInt("CASSANDRA_PORT", 9042),
		CassandraUser:     env("CASSANDRA_USERNAME", ""),
		CassandraPass:     env("CASSANDRA_PASSWORD", ""),
		CassandraKeyspace: env("CASSANDRA_KEYSPACE", "kopesa_loan_platform"),
		CassandraDC:       env("CASSANDRA_DATACENTER", "datacenter1"),
		ProviderBaseURL:   env("PROVIDER_BASE_URL", "https://cloudcalls.easipath.com/backend-email-service/api/v1"),
		DropboxBaseURL:    env("DROPBOX_BASE_URL", "https://cloudcalls.easipath.com/backend-biatechdropbox/api"),
		WhatsAppBaseURL:   env("WHATSAPP_BASE_URL", "http://safer.easipath.com:8080"),
		WhatsAppAPIKey:    env("WHATSAPP_API_KEY", "mUtombo8544e4EGG25841serEEESSA"),
		WhatsAppInstance:  env("WHATSAPP_INSTANCE", "biacibenga"),
		EmailFrom:         env("EMAIL_FROM", "no-reply@mails.biacibenga.co.za"),
	}
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	raw := env(key, "")
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func split(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	if len(out) == 0 {
		return []string{"127.0.0.1"}
	}
	return out
}
