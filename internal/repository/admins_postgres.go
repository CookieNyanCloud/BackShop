package repository

import (
	"context"
	"fmt"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type AdminsRepo struct {
	db *sqlx.DB
}

func NewAdminsRepo(db *sqlx.DB) *AdminsRepo {
	return &AdminsRepo{db: db}
}

func (r *AdminsRepo) IsAdmin(ctx context.Context, id string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE id='$1')", adminsTable)
	err := r.db.Get(&id, query, exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *AdminsRepo) AddNewAdmin(ctx context.Context, email, passwordHash string) error {
	id := uuid.New().String()
	query := fmt.Sprintf("INSERT INTO %s (id, email, password_hash) values ($1, $2, $3)",
		adminsTable)
	_, err := r.db.Exec(query, id, email, passwordHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *AdminsRepo) CreateEvent(ctx context.Context, event domain.Event) error {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (time, description, mapfile) values ($1, $2, $3) RETURNING id",
		eventsTable)
	err := r.db.QueryRow(query, event.Time, event.Description, event.MapFile).Scan(&id)
	if err != nil {
		return err
	}
	return r.AddZones(ctx, id, event.Zones)
}

func (r *AdminsRepo) AddZones(ctx context.Context, eventId int, zones []domain.Zone) error {
	query := fmt.Sprintf("INSERT INTO %s (eventId, taken, price) values ($1, $2, $3)",
		zonesTable)
	for _, zone := range zones {
		_, err := r.db.Exec(query, eventId, "", zone.Price)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AdminsRepo) DeleteEvent(ctx context.Context, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", eventsTable)
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	//todo:check if need
	//query = fmt.Sprintf("DELETE FROM %s WHERE id=$1", zonesTable)
	//_, err = r.db.Exec(query, id)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (r *AdminsRepo) DeleteUser(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", usersTable)
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *AdminsRepo) ChangeEvent(ctx context.Context, id int, event domain.Event) error {
	//todo:check race
	query := fmt.Sprintf("UPDATE %s SET time=$1, description=$2, mapfile=$3 WHERE id=$4",
		eventsTable)
	_, err := r.db.Exec(query, event.Time, event.Description, event.MapFile, event.Id)
	if err != nil {
		return err
	}
	for _, zone := range event.Zones {
		query := fmt.Sprintf("UPDATE %s SET taken=$1, price=$2 WHERE id=$3",
			zonesTable)
		_, err := r.db.Exec(query, zone.Taken, zone.Price, zone.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AdminsRepo) IsDuplicate(email string) bool {
	var id string
	query := fmt.Sprintf("SELECT id FROM %s WHERE email='$1'", usersTable)
	_ = r.db.Get(&id, query, email)
	query = fmt.Sprintf("SELECT id FROM %s WHERE 'email'='$1'", usersTable)
	if id == "" {
		return true
	}
	return false
}

func (r *AdminsRepo) GetByCredentials(ctx context.Context, email, passwordHash string) (string, error) {
	var id string
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", adminsTable)
	err := r.db.Get(&id, query, email, passwordHash)
	return id, err
}

func (r *AdminsRepo) SetSession(ctx context.Context, id string, session domain.Session) error {
	query := fmt.Sprintf("UPDATE %s SET refreshtoken = $1, expiresat = $2, lastvisitat = $3 WHERE id = $4",
		sessionsTable)
	_, err := r.db.Exec(query, session.RefreshToken, session.ExpiresAt, time.Now(), id)
	return err
}

func (r *AdminsRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	var id string
	query := fmt.Sprintf("SELECT id FROM %s WHERE refreshtoken=$1", sessionsTable)
	err := r.db.Get(&id, query, refreshToken)
	return id, err
}

