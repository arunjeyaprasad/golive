package service

type Service interface {
	Health() bool
}

type baseService struct{}

func NewService() Service {
	return &baseService{}
}

func (s *baseService) Health() bool {
	return true
}
