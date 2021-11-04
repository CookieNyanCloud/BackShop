package repository

import (
	"context"
	"fmt"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) IsDuplicate(email string) bool {
	var id string
	query := fmt.Sprintf("SELECT id FROM %s WHERE email='$1'", usersTable)
	_ = r.db.Get(&id, query, email)
	query = fmt.Sprintf("SELECT id FROM %s WHERE 'email'='$1'", usersTable)
	if id == "" {
		return true
	}
	return false
}

func (r *UsersRepo) CreateUser(ctx context.Context, user domain.User) (string, error) {
	id := uuid.New().String()
	if is := r.IsDuplicate(user.Email); is {
		return "", ErrUserAlreadyExists
	}
	query := fmt.Sprintf("INSERT INTO %s (id, email, password_hash) values ($1, $2, $3)",
		usersTable)
	_, err := r.db.Exec(query, user.Email, user.Password)
	if err != nil {
		return "", err
	}
	query = fmt.Sprintf("INSERT INTO %s (id, refreshtoken, expiresat) values ($1, $2, $3)",
		sessionsTable)
	_, err = r.db.Exec(query, id, nil, nil)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *UsersRepo) SetVerCode(ctx context.Context, id uuid.UUID, code, email string) error {
	query := fmt.Sprintf("INSERT INTO %s (id, code, state ) values ($1, $2, $3)",
		verificationTable)
	_, err := r.db.Exec(query, id, code, email)
	if err != nil {
		return err
	}
	go r.DeleteIfNotVer(id)
	return nil
}

func (r *UsersRepo) DeleteIfNotVer(id uuid.UUID) {
	time.After(time.Minute * 10)
	var state bool
	query := fmt.Sprintf("SELECT state FROM %s WHERE id=$1", verificationTable)
	_ = r.db.Get(&state, query, id)
	if !state {
		query := fmt.Sprintf("DELETE FROM %s WHERE id=$1",
			verificationTable)
		_, _ = r.db.Exec(query, id)
		query = fmt.Sprintf("DELETE FROM %s WHERE id=$1",
			usersTable)
		_, _ = r.db.Exec(query, id)
		query = fmt.Sprintf("DELETE FROM %s WHERE id=$1",
			sessionsTable)
		_, _ = r.db.Exec(query, id)
	}
}

//func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (uuid.UUID, error) {
//	var id uuid.UUID
//	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
//	err := r.db.Get(&id, query, email, password)
//	return id, err
//}

func (r *UsersRepo) GetUserInfo(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT email FROM %s where id = $1`, usersTable)
	err := r.db.Get(&user.Email, query, id)
	return user, err
}

func (r *UsersRepo) SetSession(ctx context.Context, userId uuid.UUID, session domain.Session) error {
	query := fmt.Sprintf("UPDATE %s SET refreshtoken = $1, expiresat = $2, lastvisitat = $3 WHERE id = $4",
		sessionsTable)
	_, err := r.db.Exec(query, session.RefreshToken, session.ExpiresAt, time.Now(), userId)
	return err
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("SELECT id FROM %s WHERE refreshtoken=$1", sessionsTable)
	err := r.db.Get(&id, query, refreshToken)
	return id, err
}

func (r *UsersRepo) Verify(ctx context.Context, userId string, hash string) error {
	var code string
	query := fmt.Sprintf("SELECT code FROM %s WHERE id=$1", codesTable)
	err := r.db.Get(&code, query, userId)
	if err != nil {
		return ErrUserNotFound
	}
	if code == hash {
		query := fmt.Sprintf("UPDATE %s SET verification = $1 id = $2",
			usersTable)
		_, err = r.db.Exec(query, true, userId)
		if err != nil {
			return ErrUserNotFound
		}
	} else {
		return err
	}
	return nil
}
