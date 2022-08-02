package monitoredservice

import "fmt"

const (
	Ok = iota
	Unhealthy
	Acknowledged
)

type MaxEscalatePolicyReached struct {
	error
}

type CannotEscalateHealthyService struct {
	error
}

type MonitoredService struct {
	serviceID string
	alerts    map[string]*Alert
}

type Status struct {
	serviceID             string
	health                int
	escalationPolicyLevel int
}

func (mss Status) Health() int {
	return mss.health
}

func (mss Status) EscalationPolicyLevel() int {
	return mss.escalationPolicyLevel
}

func New(serviceID string) *MonitoredService {
	return &MonitoredService{
		serviceID: serviceID,
		alerts:    make(map[string]*Alert),
	}
}

func (ms *MonitoredService) Healthy() bool {
	return len(ms.alerts) == 0
}

func (ms *MonitoredService) HasAlertWithType(alertType string) bool {
	_, ok := ms.alerts[alertType]
	return ok
}

func (ms *MonitoredService) Acknowledged(alertType string) bool {
	alert, ok := ms.alerts[alertType]
	if !ok {
		return false
	}
	return alert.acknowledgedAt > 0
}

func (ms *MonitoredService) TurnToUnhealthy(alert *Alert) error {
	if !ms.Healthy() {
		return nil
	}

	ms.alerts[alert._type] = alert
	return nil
}

func (ms *MonitoredService) TurnToHealthy(alertType string) {
	_, ok := ms.alerts[alertType]
	if !ok {
		return
	}
	delete(ms.alerts, alertType)
}

func (ms *MonitoredService) EscalatePolicy(alertType string, maxEscalationPolicyLevel int) (int, error) {
	if ms.Healthy() {
		return -1, CannotEscalateHealthyService{}
	}

	alert, ok := ms.alerts[alertType]
	if !ok {
		return 0, fmt.Errorf("no alertType for service %s and type %s", ms.serviceID, alertType)
	}

	if alert.escalationPolicyLevel >= maxEscalationPolicyLevel-1 {
		return -1, MaxEscalatePolicyReached{}
	}

	alert.escalationPolicyLevel++
	ms.alerts[alertType] = alert

	return alert.escalationPolicyLevel, nil
}

func (ms *MonitoredService) AcknowledgeAlert(alertType string, timestamp uint64) {
	if ms.Healthy() {
		return
	}

	alert, ok := ms.alerts[alertType]
	if !ok {
		return
	}

	if alert.acknowledgedAt != 0 {
		return
	}
	alert.acknowledgedAt = timestamp
}

func (ms *MonitoredService) Status() Status {

	status := Status{serviceID: ms.serviceID}

	var alert *Alert
	if ms.alerts != nil && len(ms.alerts) > 0 {
		for _, v := range ms.alerts {
			alert = v
		}
	}

	if alert == nil {
		return Status{}
	}

	status.escalationPolicyLevel = alert.escalationPolicyLevel

	if ms.Acknowledged(alert._type) {
		status.health = Acknowledged
		return status
	}

	if !ms.Healthy() {
		status.health = Unhealthy
		return status
	}

	status.health = Ok
	return status
}

type Repository interface {
	FindByServiceID(serviceID string) *MonitoredService
}
