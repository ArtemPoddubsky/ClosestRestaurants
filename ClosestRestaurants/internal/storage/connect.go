package storage

import (
	"context"
	"fmt"
	"main/internal/config"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

const (
	attempts = 8
	delay    = 5 * time.Second
)

// Postgres holds a pool of connections and methods to interact with database.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewDB connects to database and returns a new Postgres instance.
func NewDB(cfg *config.Config) *Postgres {
	psqlconn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)

	return &Postgres{connect(psqlconn)}
}

func connect(psqlconn string) *pgxpool.Pool {
	for count, attempt := 1, attempts; attempt > 0; {
		if count != 1 {
			log.Warnf("Trying to connect to database %d time", count)
		}

		ctx, cancel := context.WithTimeout(context.Background(), delay)
		pool, err := pgxpool.Connect(ctx, psqlconn)

		if err != nil {
			cancel()
			time.Sleep(delay)
			attempt--
			count++

			continue
		}

		cancel()

		return pool
	}
	log.Fatalln("Couldn't connect to database ")

	return nil
}
