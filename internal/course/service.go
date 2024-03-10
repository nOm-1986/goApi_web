package course

import (
	"log"
	"time"
)

type (
	Service interface {
		Create(name, startDate, endDate string) (*Course, error)
		Get(id string) (*Course, error)
		GetAll(filters Filters, offset, limit int) ([]Course, error)
		Update(id string, name, startDate, endDate *string) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		Name      string
		StartDate string
		EndDate   string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(name, startDate, endDate string) (*Course, error) {

	//Parseo de data string a date
	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	course := &Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}
	if err := s.repo.Create(course); err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]Course, error) {
	courses, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (s service) Get(id string) (*Course, error) {
	course, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) Update(id string, name, startDate, endDate *string) error {
	//Parseo de data string a date
	var startDateParsed, endDateParsed *time.Time

	if startDate != nil {
		date, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Println(err)
			return err
		}
		startDateParsed = &date
	}

	if endDate != nil {
		date, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Println(err)
			return err
		}
		endDateParsed = &date
	}

	return s.repo.Update(id, name, startDateParsed, endDateParsed)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
