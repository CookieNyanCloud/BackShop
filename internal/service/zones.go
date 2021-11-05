package service

import (
	"context"
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/internal/repository"
)

type ZonesService struct {
	repo repository.Zones
}

func NewZonesService(repo repository.Zones) *ZonesService {
	return &ZonesService{repo}
}

func (s *ZonesService) GetZonesByEventId(ctx context.Context,id int) ([]domain.Zone, error) {
	return s.repo.GetZonesByEventId(ctx,id)
}

func (s *ZonesService) TakeZonesById(ctx context.Context,idEvent int, idZones []int, userId string) ([]domain.Zone, error) {
	return s.repo.TakeZonesById(ctx,idEvent, idZones, userId)
}
