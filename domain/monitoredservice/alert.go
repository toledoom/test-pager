package monitoredservice

type Alert struct {
	serviceID             string
	message               string
	escalationPolicyLevel int
	createdAt             uint64
	acknowledgedAt        uint64
	healthyAt             uint64
}

func NewAlert(serviceID, message string, createdAt uint64) *Alert {
	return &Alert{
		serviceID: serviceID,
		message:   message,
		createdAt: createdAt,
	}
}

func (a *Alert) ServiceID() string {
	return a.serviceID
}
