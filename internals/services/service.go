package services

type Service interface {
	PerformWithOutput(input *ServiceInput) (ServiceOuptput, error)
	Perform(input *ServiceInput) error
}

type ServiceInput any

type ServiceOuptput any
