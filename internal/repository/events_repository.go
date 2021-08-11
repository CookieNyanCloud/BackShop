package repository

import (
	"fmt"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/jmoiron/sqlx"
)

type EventsRepo struct {
	db *sqlx.DB
}

func NewEventsRepo(db *sqlx.DB) *EventsRepo {
	return &EventsRepo{db: db}
}

func (r *EventsRepo) GetEvent() ([]domain.Event, error) {
	var events []domain.Event
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=0",
		eventsTable)
	err := r.db.Select(&events, query)
	return events, err
}