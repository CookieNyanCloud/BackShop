package repository

import (
	"fmt"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/jmoiron/sqlx"
)

//todo:delete old events

type EventsRepo struct {
	db *sqlx.DB
}

func NewEventsRepo(db *sqlx.DB) *EventsRepo {
	return &EventsRepo{db: db}
}

func (r *EventsRepo) GetEventById(id int) (domain.Event, error) {
	var event domain.Event
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1",
		eventsTable)
	err := r.db.Select(&event, query,id)
	if err != nil {
		return domain.Event{}, err
	}
	return event, err
}

func (r *EventsRepo) GetAllEvents() ([]domain.Event, error) {
	var events []domain.Event
	query := fmt.Sprintf("SELECT * FROM %s",
		eventsTable)
	err := r.db.Select(&events, query)
	return events, err
}