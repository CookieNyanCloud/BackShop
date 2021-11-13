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
//todo: payment
//todo: races
//todo: orders

func NewZonesRepo(db *sqlx.DB) *ZonesRepo {
	return &ZonesRepo{db: db}
}

func (r *ZonesRepo) GetZonesByEventId(ctx context.Context, id int) ([]domain.Zone, error) {
	var zones []domain.Zone
	query := fmt.Sprintf("SELECT * FROM %s WHERE eventid= $1", zonesTable)
	err := r.db.Select(&zones, query, id)
	return zones, err
}

func (r *ZonesRepo) GetZonesByUserId(ctx context.Context, userId string) ([]domain.Zone, error) {
	var zones []domain.Zone
	query := fmt.Sprintf("SELECT * FROM %s WHERE taken= $1", zonesTable)
	err := r.db.Select(&zones, query, userId)
	return zones, err
}

func (r *ZonesRepo) TakeZonesById(ctx context.Context, idEvent int, idZones []int, userId string) error {


	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, zoneId := range idZones {
		var zone domain.Zone

		query := fmt.Sprintf("SELECT * FROM %s WHERE eventid = $1 AND id = $2 FOR UPDATE", zonesTable)
		err := tx.QueryRowContext(ctx, query, idEvent, zoneId).Scan(&zone)
		if err != nil {
			return err
		}

		if zone.Taken == ""{
			query := fmt.Sprintf("UPDATE %s SET taken = $1 WHERE id = $2 AND eventid = $3", zonesTable)
			_, err = tx.ExecContext(ctx, query, userId, zoneId, idEvent)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

//func (r *ZonesRepo) UntakeZonesById(ctx context.Context, orderId uuid.UUID) error {
//	query := fmt.Sprintf("UPDATE %s SET taken = '' WHERE ", ordersTable)
//	err := tx.QueryRowContext(ctx, query, idEvent, zoneId).Scan(&zone)
//	if err != nil {
//		return err
//	}
//
//
//}

