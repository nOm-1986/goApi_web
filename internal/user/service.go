package user

import "log"

type Service interface {
	Create(firstName, lastName, email, phone string) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

//I can improve it making it as a variadic function
func (s service) Create(firstName, lastName, email, phone string) error {
	log.Println("Create user service")
	return nil
}
