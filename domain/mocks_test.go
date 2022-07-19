package domain_test

import (
	"github.com/toledoom/test-pager/v2/domain/escalationpolicy"
	"github.com/toledoom/test-pager/v2/domain/monitoredservice"
)

type spyTimer struct {
	setTimeOutCalls,
	timeoutInSeconds int
}

func (spyT *spyTimer) SetTimeout(serviceID string, timeoutInSeconds int) {
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

func (ssn *spySmsNotifier) Notify(targets ...escalationpolicy.Target) {
	ssn.calls++
}

func (ssn *spySmsNotifier) CalledTimes() int {
	return ssn.calls
}

type spyMailNotifier struct {
	calls int
}

func (smn *spyMailNotifier) Notify(targets ...escalationpolicy.Target) {
	smn.calls++
}

func (smn *spyMailNotifier) CalledTimes() int {
	return smn.calls
}

type dummyEscalationPolicyRepository struct {
	serviceID string
}

func (depr *dummyEscalationPolicyRepository) GetByServiceID(serviceID string) escalationpolicy.EscalationPolicy {
	levels := []escalationpolicy.Level{
		escalationpolicy.NewLevel(
			[]escalationpolicy.Target{
				escalationpolicy.NewEmailTarget("target1@test.com"),
				escalationpolicy.NewSmsTarget("+341111111"),
			},
		),
		escalationpolicy.NewLevel(
			[]escalationpolicy.Target{
				escalationpolicy.NewSmsTarget("+342222222"),
				escalationpolicy.NewSmsTarget("+343333333"),
			},
		),
		escalationpolicy.NewLevel(
			[]escalationpolicy.Target{
				escalationpolicy.NewEmailTarget("angry-cto@test.com"),
				escalationpolicy.NewSmsTarget("+349999999"),
			},
		),
	}
	return escalationpolicy.New(depr.serviceID, levels)
}

type dummyMonitoredServiceRepository struct {
	ms *monitoredservice.MonitoredService
}

func (dmsr *dummyMonitoredServiceRepository) FindByServiceID(serviceID string) *monitoredservice.MonitoredService {
	if dmsr.ms == nil {
		dmsr.ms = monitoredservice.New(serviceID)
	}

	return dmsr.ms
}
