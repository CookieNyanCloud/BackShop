package repository

import (
	"context"
	"strconv"
	"time"

	//"errors"
	"fmt"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) IsDuplicate(email, name string) (bool, error) {
	id := -1
	println("1", id)
	query := fmt.Sprintf("SELECT id FROM %s WHERE name='$1'", usersTable)
	//err:= r.db.Get(&id, query, name)
	_ = r.db.Get(&id, query, name)
	println("2", id)
	//if err != nil {
	//	return true, err
	//}
	println("3", id)
	query = fmt.Sprintf("SELECT id FROM %s WHERE 'email'='$1'", usersTable)
	//err= r.db.Get(&id, query, email)
	_ = r.db.Get(&id, query, email)
	//if err != nil {
	//	return true, err
	//}
	if id > -1 {
		return true, nil

	}
	return false, nil
}

func (r *UsersRepo) CreateUser(ctx context.Context, user domain.User) (int, error) {
	var id int
	if is, err := r.IsDuplicate(user.Email, user.Name); is || err != nil {
		return 0, ErrUserAlreadyExists
	}
	query := fmt.Sprintf("INSERT INTO %s (email, name, password_hash) values ($1, $2, $3) RETURNING id",
		usersTable)
	row := r.db.QueryRow(query, user.Email, user.Name, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	println("SDSDSDSDSDS")
	query = fmt.Sprintf("INSERT INTO %s (id, refreshtoken, expiresat) values ($1, $2, $3)",
		sessionsTable)
	_, err := r.db.Exec(query, id, nil,nil)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)
	return user, err
}

func (r *UsersRepo) GetUserInfo(ctx context.Context, id int) (domain.User, error) {
	var user domain.User
	println("repid:", id)
	//query := fmt.Sprintf("SELECT * FROM %s where id = $1", usersTable)
	query := fmt.Sprintf(`SELECT id FROM %s where id = $1`, usersTable)
	err := r.db.Get(&user, query, id)
	//var email string
	//row:=r.db.QueryRow(query,id)
	//err := row.Scan(&user)
	println("rep", user.Email)
	return user, err
}

func (r *UsersRepo) SetSession(ctx context.Context, userId string, session domain.Session) error {
	query := fmt.Sprintf("UPDATE %s SET refreshtoken = $1, expiresat = $2, lastvisitat = $3 WHERE id = $4",
		sessionsTable)
	userIdInt, err := strconv.Atoi(userId)
	_, err = r.db.Exec(query, session.RefreshToken, session.ExpiresAt, time.Now(), userIdInt)
	println("asas")
	println(session.RefreshToken)
	return err
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE refreshtoken=$1", sessionsTable)
	err := r.db.Get(&user, query, refreshToken)
	return user, err
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
