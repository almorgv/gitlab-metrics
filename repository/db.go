package repository

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"

	"gitlab-metrics/envutil"
	"gitlab-metrics/log"
)

type Db struct {
	migrationsPath string
	*sqlx.DB
	log.Loggable
}

func NewDb(host string, port string, user string, password string, dbname string) (*Db, error) {
	connectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	return NewDbFromUrl(connectionStr)
}

func NewDbFromUrl(url string) (*Db, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}
	return &Db{
		migrationsPath: "file://db/migrations",
		DB:             db,
	}, nil
}

func NewDbFromEnv() (*Db, error) {
	var db *Db
	var err error

	if dbUrl := envutil.GetEnvStr(envutil.EnvDbUrl); len(dbUrl) > 0 {
		db, err = NewDbFromUrl(dbUrl)
	} else {
		dbHost := envutil.MustGetEnvStr(envutil.EnvDbHost)
		dbPort := envutil.MustGetEnvStr(envutil.EnvDbPort)
		dbUser := envutil.MustGetEnvStr(envutil.EnvDbUser)
		dbPassword := envutil.MustGetEnvStr(envutil.EnvDbPassword)
		dbName := envutil.MustGetEnvStr(envutil.EnvDbName)
		db, err = NewDb(dbHost, dbPort, dbUser, dbPassword, dbName)
	}

	return db, err
}

func (d *Db) Migrate() error {
	driver, err := postgres.WithInstance(d.DB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create new migration driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(d.migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create new migration instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %v", err)
	}
	return nil
}
