package monitoredservice

type HealthyEvent struct {
	serviceID, alertType string
}

func NewHealthyEvent(serviceID, alertType string) *HealthyEvent {
	return &HealthyEvent{
		serviceID: serviceID,
		alertType: alertType,
	}
}

func (he *HealthyEvent) ServiceID() string {
	return he.serviceID
}

func (he *HealthyEvent) AlertType() string {
	return he.alertType
}
