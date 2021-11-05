package repository

import (
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

func (r *ZonesRepo) GetZonesByEventId(id int) ([]domain.Zone, error) {
	var zones []domain.Zone
	query := fmt.Sprintf("SELECT * FROM %s WHERE eventid= $1",
		zonesTable)
	err := r.db.Select(&zones, query, id)
	return zones, err
}

func (r *ZonesRepo) TakeZonesById(idEvent int, idZones []int, userId string) ([]domain.Zone, error) {
	//todo: payment
	var zones []domain.Zone
	for _, zoneId := range idZones {
		var zone domain.Zone
		query := fmt.Sprintf("SELECT * FROM %s WHERE eventid = $1 AND id = $2",
			zonesTable)
		err := r.db.Select(&zone, query, idEvent, zoneId)
		if err != nil {
			return []domain.Zone{}, err
		}
		if zone.Taken == 0 {
			query := fmt.Sprintf("UPDATE %s SET taken = $1 WHERE id = $2 AND eventid = $3",
				zonesTable)
			_, err = r.db.Exec(query, userId, zoneId, idEvent)
			if err != nil {
				return []domain.Zone{}, err
			}
		}
	}
	return zones, nil
}
