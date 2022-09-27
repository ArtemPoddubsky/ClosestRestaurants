package storage

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"main/internal/utils"
)

const (
	itemsOnPage     = 10
	amountOfClosest = 3
	earthRadius     = 6371
)

// FillDatabase filling database with data provided in data.csv file, returns nil if database is not empty.
func (db Postgres) FillDatabase() error {
	tag, err := db.pool.Exec(context.Background(), "SELECT * FROM restaurants LIMIT 1")

	if err != nil {
		return fmt.Errorf("exec select: %w", err)
	} else if tag.RowsAffected() == 1 {
		return nil
	}

	log.Infoln("Filling database...")

	_, err = db.pool.Exec(context.Background(),
		`COPY restaurants(ID, Name, Address, Phone, Longitude, Latitude) 
				FROM '/data.csv' 
				DELIMITER '	' 
				CSV HEADER`)

	if err != nil {
		return fmt.Errorf("exec fill: %w", err)
	}

	return nil
}

// GetPage returns page of restaurants from database, based on page number.
func (db Postgres) GetPage(ctx context.Context, pageNum int) (utils.HTMLPlaces, error) {
	content := utils.HTMLPlaces{
		Total:  0,
		Page:   pageNum,
		Places: make([]utils.Restaurant, itemsOnPage),
	}

	err := db.pool.QueryRow(ctx, "SELECT COUNT(*) FROM restaurants").Scan(&content.Total)
	if err != nil {
		return content, fmt.Errorf("select count rows.Scan: %w", err)
	}

	if pageNum > content.Total/itemsOnPage {
		return content, fmt.Errorf("content out of bounds: %w", err)
	}

	rows, err := db.pool.Query(ctx,
		`SELECT name, address, phone FROM restaurants LIMIT $1 OFFSET $2`, itemsOnPage, pageNum*itemsOnPage)
	if err != nil {
		return content, fmt.Errorf("query: %w", err)
	}

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&content.Places[i].Name, &content.Places[i].Address, &content.Places[i].Phone)
		if err != nil {
			return content, fmt.Errorf("content rows.Scan: %w", err)
		}
	}

	return content, nil
}

// GetClosest returns 3 closest restaurants from database, based on provided coordinates.
func (db Postgres) GetClosest(ctx context.Context, lat, lon float64) (utils.JSONPlaces, error) {
	content := utils.JSONPlaces{
		Name:   "Recommendation",
		Places: make([]utils.Restaurant, amountOfClosest),
	}
	rows, err := db.pool.Query(ctx,
		`SELECT * FROM restaurants 
			 ORDER BY ACOS(SIN($1)*SIN(Latitude)+COS($1)*COS(Latitude)*COS($2-Longitude)) * $3 
			 LIMIT 3`, lat, lon, earthRadius,
	)

	if err != nil {
		return content, fmt.Errorf("ThreeClosest query: %w", err)
	}

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&content.Places[i].ID, &content.Places[i].Name, &content.Places[i].Address,
			&content.Places[i].Phone, &content.Places[i].Location.Lon, &content.Places[i].Location.Lat)

		if err != nil {
			return content, fmt.Errorf("rows.Scan: %w", err)
		}
	}

	return content, nil
}
