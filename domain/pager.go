package domain

import (
	"github.com/toledoom/test-pager/v2/domain/escalationpolicy"
	"github.com/toledoom/test-pager/v2/domain/monitoredservice"
	"github.com/toledoom/test-pager/v2/domain/notifier"
	"github.com/toledoom/test-pager/v2/domain/timer"
)

type Pager struct {
	monitoredServiceRepository monitoredservice.Repository
	escalationPolicyRepository escalationpolicy.Repository
	notifier                   notifier.Notifier
	timer                      timer.Timer
	timerConfiguration         timer.Configuration
}

func NewPager(
	monitoredServiceRepository monitoredservice.Repository,
	escalationPolicyRepository escalationpolicy.Repository,
	notifier notifier.Notifier,
	timer timer.Timer,
	timerConfiguration timer.Configuration) *Pager {
	return &Pager{
		monitoredServiceRepository: monitoredServiceRepository,
		escalationPolicyRepository: escalationPolicyRepository,
		notifier:                   notifier,
		timer:                      timer,
		timerConfiguration:         timerConfiguration,
	}
}

func (p *Pager) SendAlert(alert *monitoredservice.Alert) {
	monitoredService := p.monitoredServiceRepository.FindByServiceID(alert.ServiceID())

	if !monitoredService.Healthy() || monitoredService.Acknowledged() {
		return
	}

	monitoredService.TurnToUnhealthy(alert)

	timeoutInSeconds := p.timerConfiguration.GetTimeoutInSeconds()
	p.timer.SetTimeout(alert.ServiceID(), timeoutInSeconds)

	ep := p.escalationPolicyRepository.GetByServiceID(alert.ServiceID())
	targets := ep.GetTargetsByLevel(0)
	p.notifier.Notify(targets...)
}

func (p *Pager) AcknowledgeAlert(serviceID string, timestamp uint64) {
	monitoredService := p.monitoredServiceRepository.FindByServiceID(serviceID)
	monitoredService.AcknowledgeAlert(timestamp)
}

func (p *Pager) SendHealthyEvent(healthyEvent *monitoredservice.HealthyEvent) {
	monitoredService := p.monitoredServiceRepository.FindByServiceID(healthyEvent.ServiceID())
	monitoredService.TurnToHealthy()
}

func (p *Pager) NotifyAckTimeout(serviceID string) {
	monitoredService := p.monitoredServiceRepository.FindByServiceID(serviceID)
	if monitoredService.Acknowledged() {
		return
	}

	escalationPolicy := p.escalationPolicyRepository.GetByServiceID(serviceID)

	nextEscalationPolicyLevel, err := monitoredService.EscalatePolicy(escalationPolicy.MaxLevel())

	_, maxEscalatePolicyReached := err.(*monitoredservice.MaxEscalatePolicyReached)
	if maxEscalatePolicyReached {
		return
	}

	timeoutInSeconds := p.timerConfiguration.GetTimeoutInSeconds()
	p.timer.SetTimeout(serviceID, timeoutInSeconds)

	targets := escalationPolicy.GetTargetsByLevel(nextEscalationPolicyLevel)
	p.notifier.Notify(targets...)
}

func (p *Pager) Status(serviceID string) monitoredservice.Status {
	ms := p.monitoredServiceRepository.FindByServiceID(serviceID)
	return ms.Status()
}
