package repository

import (
	"context"
	"fmt"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ZonesRepo struct {
	db *sqlx.DB
}

func NewZonesRepo(db *sqlx.DB) *ZonesRepo {
	return &ZonesRepo{db: db}
}

func (r *ZonesRepo) GetZonesByEventId(ctx context.Context, id int) ([]domain.Zone, error) {
	var zones []domain.Zone
	query := fmt.Sprintf("SELECT * FROM %s WHERE eventid= $1",
		zonesTable)
	err := r.db.Select(&zones, query, id)
	return zones, err
}

func (r *ZonesRepo) TakeZonesById(ctx context.Context, idEvent int, idZones []int, userId string) ([]domain.Zone, error) {
	//todo: payment
	//todo: races
	var zones []domain.Zone

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return []domain.Zone{}, err
	}
	defer tx.Rollback()

	for _, zoneId := range idZones {
		var zone domain.Zone

		query := fmt.Sprintf("SELECT * FROM %s WHERE eventid = $1 AND id = $2",
			zonesTable)
		err := tx.QueryRowContext(ctx, query, idEvent, zoneId).Scan(&zone)
		if err != nil {
			return []domain.Zone{}, err
		}

		if zone.Taken == "" {
			query := fmt.Sprintf("UPDATE %s SET taken = $1 WHERE id = $2 AND eventid = $3",
				zonesTable)
			_, err = tx.ExecContext(ctx, query, userId, zoneId, idEvent)
			if err != nil {
				return []domain.Zone{}, err
			}
		} else {
			return []domain.Zone{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []domain.Zone{}, err
	}
	return zones, nil
}
