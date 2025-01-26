package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"greenlight.Arkitecth.github/internal/data"
)

const version = "1.0.0"

type config struct {
	port int 
	env string
	db struct {
		dsn string
		maxOpenConnections int
		maxIdleConnecttions int
		maxIdleTimeout time.Duration
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Model
}


func main() {
	var cfg config 

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.IntVar(&cfg.db.maxOpenConnections, "db-max-open-conn", 25, "PostgresSQL Max Open Connections")
	flag.IntVar(&cfg.db.maxOpenConnections, "db-max-idle-conn", 25, "PostgresSQL Max Idle Connections")
	flag.DurationVar(&cfg.db.maxIdleTimeout, "db-max-idle-timeout", 15 * time.Minute, "PostgresSQL Max Idle Time Out")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")


	flag.Parse() 

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	defer db.Close() 

	logger.Info("database connection pool established")

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server on", "addr", srv.Addr, "env", cfg.env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}



func openDB(config config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.db.dsn) 
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.db.maxOpenConnections)

	db.SetMaxIdleConns(config.db.maxIdleConnecttions)

	db.SetConnMaxIdleTime(config.db.maxIdleTimeout)

	
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second) 
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil

}