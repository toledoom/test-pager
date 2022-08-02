package monitoredservice

import "fmt"

type Alert struct {
	_type                 string
	serviceID             string
	message               string
	escalationPolicyLevel int
	createdAt             uint64
	acknowledgedAt        uint64
	healthyAt             uint64
}

func NewAlert(serviceID, _type, message string, createdAt uint64) *Alert {
	return &Alert{
		_type:     _type,
		serviceID: serviceID,
		message:   message,
		createdAt: createdAt,
	}
}

func (a *Alert) ID() string {
	return fmt.Sprintf("%s:%s", a.serviceID, a._type)
}

func (a *Alert) Type() string {
	return a._type
}
