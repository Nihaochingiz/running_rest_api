package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := initStore()
	if err != nil {
		log.Fatalf("failed to initialise the store: %s", err)
	}
	defer db.Close()

	e.GET("/", func(c echo.Context) error {
		return rootHandler(db, c)
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.POST("/running_statistics", func(c echo.Context) error {
		return sendHandler(db, c)
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

type Message struct {
	Date     time.Time `json:"date"`
	Distance string    `json:"distance"`
	Time     string    `json:"time"`
}

func initStore() (*sql.DB, error) {
	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	var (
		db  *sql.DB
		err error
	)
	openDB := func() error {
		db, err = sql.Open("postgres", pgConnString)
		return err
	}

	err = backoff.Retry(openDB, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS running_statistics (
			id SERIAL PRIMARY KEY,
			date DATE,
			distance VARCHAR(10),
			time VARCHAR(10),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(date, distance, time)
		)`); err != nil {
		return nil, err
	}

	return db, nil
}
func rootHandler(db *sql.DB, c echo.Context) error {
	r, err := countRecords(db)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, fmt.Sprintf("Hello, Docker! (%d)\n", r))
}

func sendHandler(db *sql.DB, c echo.Context) error {
	m := &Message{}

	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var recordCreated bool

	err := crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {
		result, err := tx.Exec(
			"INSERT INTO running_statistics (date, distance, time) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			m.Date, m.Distance, m.Time,
		)
		if err != nil {
			return err
		}

		rowsAffected, _ := result.RowsAffected()
		recordCreated = rowsAffected > 0

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	response := map[string]interface{}{
		"message":        "Record creation status",
		"record_created": recordCreated,
	}

	return c.JSON(http.StatusOK, response)
}
func countRecords(db *sql.DB) (int, error) {

	rows, err := db.Query("SELECT COUNT(*) FROM running_statistics")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		rows.Close()
	}

	return count, nil
}
