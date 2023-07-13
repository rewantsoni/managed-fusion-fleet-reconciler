package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/types"
	"go.uber.org/zap"
)

func GetConnectionString(host, user, password, name string, port int) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		name,
	)
}

type Database struct {
	tables struct {
		providers string
	}
	pool *pgxpool.Pool
	// conn is used to listen for notifications and should be closed
	conn *pgxpool.Conn
}

func NewClient(ctx context.Context, connString string, tables map[string]string) (*Database, error) {
	if _, ok := tables["providers"]; !ok {
		return nil, fmt.Errorf("incomplete table name mapping, missing a table name for \"providers\" table")
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create new instace of pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db := Database{}
	db.pool = pool
	db.tables.providers = tables["providers"]
	return &db, nil
}

func (pg *Database) Close(ctx context.Context) {
	if pg.conn != nil {
		pg.conn.Conn().Close(ctx)
	}
	pg.pool.Close()
}

func (pg *Database) OnProvider(ctx context.Context, logger *zap.Logger, notifyExisting bool, fn func(providerName string)) error {
	var err error
	if pg.conn, err = pg.pool.Acquire(ctx); err != nil {
		return fmt.Errorf("failed to acquire connection: %v", err)
	}

	query := fmt.Sprintf("LISTEN %s;", pg.tables.providers)
	if _, err = pg.conn.Exec(ctx, query); err != nil {
		return err
	}
	go func() {
		for {
			notification, err := pg.conn.Conn().WaitForNotification(ctx)
			if err != nil {
				// If the connection is closed, release the connection and exit the goroutine
				if pg.conn.Conn().IsClosed() {
					pg.conn.Release()
					logger.Error("Connection closed by server:", zap.Error(err))
					return
				}
				logger.Error("failed to wait for notification", zap.Error(err))
				continue
			}
			fn(notification.Payload)
		}
	}()

	if notifyExisting {
		go func() {
			query := fmt.Sprintf("SELECT cluster_id FROM %s", pg.tables.providers)
			rows, err := pg.pool.Query(ctx, query)
			if err != nil {
				logger.Fatal("failed to query database", zap.Error(err))
			}
			defer rows.Close()

			for rows.Next() {
				var clusterID string
				if err := rows.Scan(&clusterID); err != nil {
					logger.Fatal("failed to scan row", zap.Error(err))
				}
				fn(clusterID)
			}
		}()
	}
	return nil
}

func (pg *Database) GetProviderCluster(ctx context.Context, clusterID string) (*types.ProviderCluster, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE cluster_id = $1", pg.tables.providers)
	row := pg.pool.QueryRow(ctx, query, clusterID)

	var c types.ProviderCluster
	if err := row.Scan(&c.ClusterID, &c.AccountID, &c.SatelliteID, &c.MetaData, &c.Spec, &c.Status); err != nil {
		return nil, err
	}

	return &c, nil
}
