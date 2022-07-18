package model

type HealthyEvent struct {
	serviceID string
}

func NewHealthyEvent(serviceID string) *HealthyEvent {
	return &HealthyEvent{
		serviceID: serviceID,
	}
}

func (he *HealthyEvent) ServiceID() string {
	return he.serviceID
}
