package domain

import (
	"github.com/toledoom/test-pager/v2/domain/model"
	"github.com/toledoom/test-pager/v2/domain/service"
)

type Pager struct {
	monitoredServiceRepository model.MonitoredServiceRepository
	escalationPolicyRepository model.EscalationPolicyRepository
	notifier                   service.Notifier
	timer                      service.Timer
	timerConfiguration         service.TimerConfiguration
}

func NewPager(
	monitoredServiceRepository model.MonitoredServiceRepository,
	escalationPolicyRepository model.EscalationPolicyRepository,
	notifier service.Notifier,
	timer service.Timer,
	timerConfiguration service.TimerConfiguration) *Pager {
	return &Pager{
		monitoredServiceRepository: monitoredServiceRepository,
		escalationPolicyRepository: escalationPolicyRepository,
		notifier:                   notifier,
		timer:                      timer,
		timerConfiguration:         timerConfiguration,
	}
}

func (p *Pager) SendAlert(alert *model.Alert) {
	monitoredService := p.monitoredServiceRepository.FindByServiceID(alert.ServiceID())

	if !monitoredService.Healthy() || monitoredService.Acknowledged() {
		return
	}

	monitoredService.TurnToUnhealthy(alert)

	timeoutInSeconds := p.timerConfiguration.GetTimeoutInSeconds()
	p.timer.SetTimeout(timeoutInSeconds)

	ep := p.escalationPolicyRepository.GetByServiceID(alert.ServiceID())
	targets := ep.GetTargetsByLevel(0)
	p.notifier.Notify(targets...)
}

func (p *Pager) AcknowledgeAlert(serviceID string, timestamp uint64) {
	monitoredService := p.monitoredServiceRepository.FindByServiceID(serviceID)
	monitoredService.AcknowledgeAlert(timestamp)
}

func (p *Pager) SendHealthyEvent(healthyEvent *model.HealthyEvent) {
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

	_, maxEscalatePolicyReached := err.(*model.MaxEscalatePolicyReached)
	if maxEscalatePolicyReached {
		return
	}

	timeoutInSeconds := p.timerConfiguration.GetTimeoutInSeconds()
	p.timer.SetTimeout(timeoutInSeconds)

	targets := escalationPolicy.GetTargetsByLevel(nextEscalationPolicyLevel)
	p.notifier.Notify(targets...)
}

func (p *Pager) Status(serviceID string) model.MonitoredServiceStatus {
	ms := p.monitoredServiceRepository.FindByServiceID(serviceID)
	return ms.Status()
}
