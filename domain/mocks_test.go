package domain_test

import (
	"github.com/toledoom/test-pager/v2/domain/model"
)

type spyTimer struct {
	setTimeOutCalls,
	timeoutInSeconds int
}

func (spyT *spyTimer) SetTimeout(timeoutInSeconds int) {
	spyT.timeoutInSeconds = timeoutInSeconds
	spyT.setTimeOutCalls++
}

func (spyT *spyTimer) CalledTimes() int {
	return spyT.setTimeOutCalls
}

type dummyTimerConfiguration struct {
	timeoutInSeconds int
}

func (dtc *dummyTimerConfiguration) GetTimeoutInSeconds() int {
	return dtc.timeoutInSeconds
}

type spySmsNotifier struct {
	calls int
}

func (ssn *spySmsNotifier) Notify(targets ...model.Target) {
	ssn.calls++
}

func (ssn *spySmsNotifier) CalledTimes() int {
	return ssn.calls
}

type spyMailNotifier struct {
	calls int
}

func (smn *spyMailNotifier) Notify(targets ...model.Target) {
	smn.calls++
}

func (smn *spyMailNotifier) CalledTimes() int {
	return smn.calls
}

type dummyEscalationPolicyRepository struct {
	serviceID string
}

func (depr *dummyEscalationPolicyRepository) GetByServiceID(serviceID string) model.EscalationPolicy {
	levels := []model.Level{
		model.NewLevel(
			[]model.Target{
				model.NewEmailTarget("target1@test.com"),
				model.NewSmsTarget("+341111111"),
			},
		),
		model.NewLevel(
			[]model.Target{
				model.NewSmsTarget("+342222222"),
				model.NewSmsTarget("+343333333"),
			},
		),
		model.NewLevel(
			[]model.Target{
				model.NewEmailTarget("angry-cto@test.com"),
				model.NewSmsTarget("+349999999"),
			},
		),
	}
	return model.NewEscalationPolicy(depr.serviceID, levels)
}

type dummyMonitoredServiceRepository struct {
	ms *model.MonitoredService
}

func (dmsr *dummyMonitoredServiceRepository) FindByServiceID(serviceID string) *model.MonitoredService {
	if dmsr.ms == nil {
		dmsr.ms = model.NewMonitoredService(serviceID)
	}

	return dmsr.ms
}
