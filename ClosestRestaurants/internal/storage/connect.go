package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"main/internal/config"
	"time"
)

const Attempts = 8
const Delay = 5 * time.Second

type Postgres struct {
	pool *pgxpool.Pool
}

func NewDB(cfg config.Config) *Postgres {
	psqlconn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.Db.Username, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Database)
	return &Postgres{Connect(psqlconn)}
}

func Connect(psqlconn string) *pgxpool.Pool {
	for count, attempt := 1, Attempts; attempt > 0; {
		if count != 1 {
			log.Warnf("Trying to connect to database %d time", count)
		}
		ctx, cancel := context.WithTimeout(context.Background(), Delay)
		pool, err := pgxpool.Connect(ctx, psqlconn)
		if err != nil {
			cancel()
			time.Sleep(Delay)
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
