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

func (r *ZonesRepo) GetZonesByEventId(id int)([]domain.Zone, error)  {
	var zones []domain.Zone
	query := fmt.Sprintf("SELECT * FROM %s WHERE eventid= $1",
		zonesTable)
	err := r.db.Select(&zones, query, id)
	println(id)
	return zones, err
}

func (r *ZonesRepo) TakeZoneById(idEvent, idZone, userId int) ([]domain.Zone, error)  {
	var zones []domain.Zone
	query := fmt.Sprintf("SELECT * FROM %s WHERE eventid = $1 AND id = $2",
		zonesTable)
	err := r.db.Select(&zones, query, idEvent, idZone)
	if zones[0].Taken == 0 {
		query :=fmt.Sprintf("UPDATE %s SET taken = $1 WHERE id = $2 AND eventid = $3",
			zonesTable)
		_,err = r.db.Exec(query, userId,idZone,idEvent)
	}
	return zones, err
}
