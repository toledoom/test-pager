package monitoredservice

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
	alert     *Alert
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
	}
}

func (ms *MonitoredService) Healthy() bool {
	return ms.alert == nil || ms.alert.healthyAt > 0
}

func (ms *MonitoredService) Acknowledged() bool {
	return ms.alert != nil && ms.alert.acknowledgedAt > 0
}

func (ms *MonitoredService) TurnToUnhealthy(alert *Alert) error {
	if !ms.Healthy() {
		return nil
	}

	ms.alert = alert
	return nil
}

func (ms *MonitoredService) TurnToHealthy() {
	ms.alert = nil
}

func (ms *MonitoredService) EscalatePolicy(maxEscalationPolicyLevel int) (int, error) {
	if ms.Healthy() {
		return -1, CannotEscalateHealthyService{}
	}

	if ms.alert.escalationPolicyLevel >= maxEscalationPolicyLevel-1 {
		return -1, MaxEscalatePolicyReached{}
	}

	ms.alert.escalationPolicyLevel++

	return ms.alert.escalationPolicyLevel, nil
}

func (ms *MonitoredService) AcknowledgeAlert(timestamp uint64) {
	if ms.Healthy() {
		return
	}

	if ms.alert.acknowledgedAt != 0 {
		return
	}
	ms.alert.acknowledgedAt = timestamp
}

func (ms *MonitoredService) Status() Status {

	status := Status{serviceID: ms.serviceID}

	if ms.alert != nil {
		status.escalationPolicyLevel = ms.alert.escalationPolicyLevel
	}

	if ms.Acknowledged() {
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
