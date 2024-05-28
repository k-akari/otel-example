package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Client struct {
	db *sql.DB
}

func NewClient(user, pass, host, name string, port int) (*Client, func(), error) {
	db, err := otelsql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, pass, host, port, name), otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	otelsql.ReportDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))

	cleanup := func() {
		if err := db.Close(); err != nil {
			log.Print("failed to close database connection")
		} else {
			log.Print("successfully closed database connection")
		}
	}

	return &Client{db: db}, cleanup, nil
}

func (c *Client) Query(ctx context.Context) error {
	const instrumentationName = "github.com/k-akari/otel-example/internal/infra/database"
	ctx, span := otel.Tracer(instrumentationName).Start(ctx, "Query")
	defer span.End()

	err := c.query(ctx)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (c *Client) query(ctx context.Context) error {
	rows, err := c.db.QueryContext(ctx, `SELECT CURRENT_TIMESTAMP`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var currentTime time.Time
	for rows.Next() {
		err = rows.Scan(&currentTime)
		if err != nil {
			return err
		}
	}
	fmt.Println(currentTime)
	return nil
}
