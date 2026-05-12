package broker_service

import (
	"mini-stock-exchange/internal/dto"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"

	"github.com/google/uuid"
)

type BrokerService interface {
	GetBroker(req dto.GetBrokerRequest) (dto.GetBrokerResponse, error)
	GetBrokerByID(id uuid.UUID) (entity.Broker, error)
	Create(req dto.CreateBrokerRequest) (dto.CreateBrokerResponse, error)
}

type brokerService struct {
	brokerRepo repository.BrokerRepository
}

func NewBrokerService(brokerRepo repository.BrokerRepository) BrokerService {
	return &brokerService{
		brokerRepo: brokerRepo,
	}
}

func (s *brokerService) GetBroker(req dto.GetBrokerRequest) (dto.GetBrokerResponse, error) {
	broker, err := s.brokerRepo.GetByID(req.ID)
	if err != nil {
		return dto.GetBrokerResponse{}, err
	}
	return dto.NewGetBrokerResponse(broker)
}

func (s *brokerService) GetBrokerByID(id uuid.UUID) (entity.Broker, error) {
	return s.brokerRepo.GetByID(id)
}

func (s *brokerService) Create(req dto.CreateBrokerRequest) (dto.CreateBrokerResponse, error) {
	broker, err := req.ToBroker()
	if err != nil {
		return dto.CreateBrokerResponse{}, err
	}
	err = s.brokerRepo.Insert(broker)
	if err != nil {
		return dto.CreateBrokerResponse{}, err
	}
	return dto.NewCreateBrokerResponse(broker)
}
