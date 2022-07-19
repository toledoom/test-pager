package escalationpolicy

type EscalationPolicy struct {
	serviceID string
	levels    []Level
}

func New(serviceID string, levels []Level) EscalationPolicy {
	return EscalationPolicy{
		serviceID: serviceID,
		levels:    levels,
	}
}

func (ep *EscalationPolicy) GetTargetsByLevel(level int) []Target {
	return ep.levels[level].targets
}

func (ep *EscalationPolicy) MaxLevel() int {
	return len(ep.levels)
}

type Level struct {
	targets []Target
}

func NewLevel(targets []Target) Level {
	return Level{
		targets: targets,
	}
}

type Target interface {
	GetMetadata() map[string]string
}

type EmailTarget struct {
	emailAddress string
}

func NewEmailTarget(emailAddress string) *EmailTarget {
	return &EmailTarget{
		emailAddress: emailAddress,
	}
}

func (et *EmailTarget) GetMetadata() map[string]string {
	return map[string]string{
		"emailAddress": et.emailAddress,
	}
}

type SmsTarget struct {
	phoneNumber string
}

func NewSmsTarget(phoneNumber string) *SmsTarget {
	return &SmsTarget{
		phoneNumber: phoneNumber,
	}
}

func (et *SmsTarget) GetMetadata() map[string]string {
	return map[string]string{
		"phoneNumber": et.phoneNumber,
	}
}

type Repository interface {
	GetByServiceID(serviceID string) EscalationPolicy
}
