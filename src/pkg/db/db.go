package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	currentDatabase = "central_active"

	databasePasswordFile = "/var/run/secrets/stackrox.io/db-password/password"

	databaseConfigSource = `host=central-db.stackrox.svc
port=5432
user=postgres
sslmode=verify-full
sslrootcert=/run/secrets/stackrox.io/certs/ca.pem
statement_timeout=1.2e+06
pool_min_conns=10
pool_max_conns=90
client_encoding=UTF8
password=%s`
)

func getDBConfig() (*pgxpool.Config, error) {
	password, err := os.ReadFile(databasePasswordFile)
	if err != nil {
		return nil, errors.Wrapf(err, "pgsql: could not load password file %q", databasePasswordFile)
	}
	source := fmt.Sprintf(databaseConfigSource, password)
	config, err := pgxpool.ParseConfig(source)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse postgres config")
	}
	config.ConnConfig.Database = currentDatabase
	return config, nil
}

func GetDBConn(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := getDBConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get postgres config")
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get postgres pool")
	}
	_, err = pool.Exec(ctx, "create extension if not exists pg_stat_statements")
	if err != nil {
		return nil, errors.Wrap(err, "Could not create extension pg_stat_statements")
	}
	return pool, nil
}
