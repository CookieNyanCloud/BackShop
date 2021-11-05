package service

import (
	"context"
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

func (s *EventsService) GetEvents(ctx context.Context) ([]domain.Event, error) {
	return s.repo.GetAllEvents(ctx)
}

func (s *EventsService) GetEventById(ctx context.Context,id int) (domain.Event, error) {
	return s.repo.GetEventById(ctx,id)
}
