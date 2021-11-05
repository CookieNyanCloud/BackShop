package service

import (
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/internal/repository"
)

type EventsService struct {
	repo  repository.Events
	Zones repository.Zones
}

func NewEventsService(repo repository.Events, zones repository.Zones) *EventsService {
	return &EventsService{repo, zones}
}

func (s *EventsService) GetEvent() ([]domain.Event, error) {
	return s.repo.GetAllEvents()
}

func (s *EventsService) GetEventById(id int) (domain.Event, error) {
	return s.repo.GetEventById(id)
}
