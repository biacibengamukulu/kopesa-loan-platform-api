package cassandra

import (
	"fmt"
	"strings"

	"github.com/gocql/gocql"

	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/config"
)

func NewCluster(cfg config.Config, withKeyspace bool) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(cfg.CassandraHosts...)
	cluster.Port = cfg.CassandraPort
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * 1e9
	cluster.ConnectTimeout = 10 * 1e9
	cluster.ProtoVersion = 4
	cluster.DisableInitialHostLookup = true
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	if cfg.CassandraDC != "" {
		cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(
			gocql.DCAwareRoundRobinPolicy(cfg.CassandraDC),
		)
	}

	if withKeyspace {
		cluster.Keyspace = cfg.CassandraKeyspace
	}

	if strings.TrimSpace(cfg.CassandraUser) != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.CassandraUser,
			Password: cfg.CassandraPass,
		}
	}
	return cluster
}

func NewSession(cfg config.Config, withKeyspace bool) (*gocql.Session, error) {
	session, err := NewCluster(cfg, withKeyspace).CreateSession()
	if err != nil {
		return nil, fmt.Errorf("create cassandra session: %w", err)
	}
	return session, nil
}
