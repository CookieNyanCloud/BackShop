package repository

import (
	"context"
	"fmt"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type AdminsRepo struct {
	db *sqlx.DB
}

func NewAdminsRepo(db *sqlx.DB) *AdminsRepo {
	return &AdminsRepo{db: db}
}

func (r *AdminsRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var admin domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err:= r.db.Get(&admin, query,email,password)
	return admin, err
//todo:admin hash
}

func (r *AdminsRepo) SetSession(ctx context.Context, adminId string, session domain.Session) error {
	query :=fmt.Sprintf("UPDATE %s SET refreshtoken = $1, expiresat = $2, lastvisitat = $3 WHERE id = $4",
		sessionsTable)
	userIdInt, err := strconv.Atoi(adminId)
	_,err = r.db.Exec(query,session.RefreshToken,session.ExpiresAt,time.Now(),userIdInt)
	println(session.RefreshToken)
	return err
}

func (r *AdminsRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE refreshtoken=$1", sessionsTable)
	err:= r.db.Get(&user, query,refreshToken)
	return user, err
}

func (r *AdminsRepo) CreateEvent(time time.Time, description string, zones []domain.Zone) (int, error){
	var id int
	query := fmt.Sprintf("INSERT INTO %s (time) values ($1) RETURNING id",
		eventsTable)
	row := r.db.QueryRow(query, time)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	for _, zone:= range zones {
		query = fmt.Sprintf("INSERT INTO %s (eventId,taken,price) values ($1, $2, $3)",
			zonesTable)
		_, err := r.db.Exec(query, id, false , zone.Price)
		println("range", zone.Price)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
//	todo:check
}

