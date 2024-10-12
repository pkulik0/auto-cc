package service

// Service is an interface for the service layer.
//
//go:generate mockgen -destination=../mock/service.go -package=mock . Service
type Service interface {
}

var _ Service = &service{}

type service struct {
}

func New() Service {
	return &service{}
}
