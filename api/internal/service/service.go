package service

type Service interface {
}

var _ Service = &service{}

type service struct {
}

func New() Service {
	return &service{}
}
