package storage

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"main/internal/utils"
)

const EarthRadius = 6371

func (db Postgres) FillDatabase() error {
	tag, err := db.pool.Exec(context.Background(), "SELECT * FROM restaurants LIMIT 1")
	if err != nil {
		return err
	} else if tag.RowsAffected() == 1 {
		return nil
	}
	log.Infoln("Filling data into database...")
	_, err = db.pool.Exec(context.Background(),
		`COPY restaurants(Id, Name, Address, Phone, Longitude, Latitude) 
				FROM '/data.csv' 
				DELIMITER '	' 
				CSV HEADER`)
	return err
}

func (db Postgres) GetPage(page int) (utils.HTMLPlaces, error) {
	r := utils.HTMLPlaces{
		Page:   page,
		Places: make([]utils.Restaurant, 10),
	}
	row := db.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM restaurants")
	if err := row.Scan(&r.Total); err != nil {
		return r, err
	} else if page > r.Total/10 {
		return r, errors.New("This page doesn't exist ")
	}

	rows, err := db.pool.Query(context.Background(),
		`SELECT name, address, phone FROM restaurants LIMIT 10 OFFSET $1`, page+10)
	if err != nil {
		return r, err
	}
	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&r.Places[i].Name, &r.Places[i].Address, &r.Places[i].Phone)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func (db Postgres) ThreeClosest(lat, lon float64) (utils.JSONPlaces, error) {
	var r = utils.JSONPlaces{
		Name:   "Recommendation",
		Places: make([]utils.Restaurant, 3),
	}
	rows, err := db.pool.Query(context.Background(),
		`SELECT * FROM restaurants 
			 ORDER BY ACOS(SIN($1)*SIN(Latitude)+COS($1)*COS(Latitude)*COS($2-Longitude)) * $3 
			 LIMIT 3`, lat, lon, EarthRadius,
	)
	if err != nil {
		return r, err
	}
	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&r.Places[i].Id, &r.Places[i].Name, &r.Places[i].Address,
			&r.Places[i].Phone, &r.Places[i].Location.Lon, &r.Places[i].Location.Lat)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}
