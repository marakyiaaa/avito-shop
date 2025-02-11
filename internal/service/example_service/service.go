package example_service

type Repository interface {
}

type Service struct {
	repo Repository
}

func New() *Service {
	return &Service{}
}
