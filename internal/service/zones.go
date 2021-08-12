package service

import (
	"github.com/cookienyancloud/back/internal/domain"
	"github.com/cookienyancloud/back/internal/repository"
)

type ZonesService struct {
	repo repository.Zones
}

func NewZonesService(repo repository.Zones) *ZonesService {
	return &ZonesService{repo}
}

func (s *ZonesService) GetZonesByEventId(id int) ([]domain.Zone, error) {
	return s.repo.GetZonesByEventId(id)
}

func (s *ZonesService) TakeZoneById(idEvent, idZone, userId int) ([]domain.Zone, error) {
	return s.repo.TakeZoneById(idEvent, idZone, userId)
}