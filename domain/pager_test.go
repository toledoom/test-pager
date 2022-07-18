package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toledoom/test-pager/v2/domain"
	"github.com/toledoom/test-pager/v2/domain/model"
	"github.com/toledoom/test-pager/v2/domain/service"
)

const myServiceID = "my-service-id"

func TestUseCase1(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	msr := &dummyMonitoredServiceRepository{}
	epr := &dummyEscalationPolicyRepository{}
	smsNotifier := &spySmsNotifier{}
	mailNotifier := &spyMailNotifier{}
	notifier := service.NewCompositeNotifier(smsNotifier, mailNotifier)
	timer := &spyTimer{}
	timerConfiguration := &dummyTimerConfiguration{timeoutInSeconds: 15 * 60}
	pager := domain.NewPager(msr, epr, notifier, timer, timerConfiguration)

	// Act
	now := 10
	alert := model.NewAlert(myServiceID, "a message", uint64(now))
	pager.SendAlert(alert)

	// Assert monitored service is healthy
	status := pager.Status(myServiceID)
	assert.Equal(model.Unhealthy, status.Health())
	// Assert timeout set to timer service
	assert.Equal(1, timer.CalledTimes())
	assert.Equal(15*60, timer.timeoutInSeconds)
	// Assert all sms and mail targets have been notified
	assert.Equal(1, smsNotifier.calls)
	assert.Equal(1, mailNotifier.calls)
}

func TestUseCase2(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	smsNotifier, mailNotifier, timer, pager := arrangeUnhealthyDependencies()

	// Act
	pager.NotifyAckTimeout(myServiceID)

	// Assert timeout set to timer service
	assert.Equal(1, timer.CalledTimes())
	assert.Equal(15*60, timer.timeoutInSeconds)
	// Assert all sms and mail targets have been notified: second level!
	assert.Equal(2, smsNotifier.calls)
	assert.Equal(0, mailNotifier.calls)
}

func TestUseCase3(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	smsNotifier, mailNotifier, timer, pager := arrangeUnhealthyDependencies()

	// Act
	acknowledgeAt := uint64(10)
	pager.AcknowledgeAlert(myServiceID, acknowledgeAt)
	pager.NotifyAckTimeout(myServiceID)

	// Assert no timeout is sent to timer service
	assert.Equal(0, timer.CalledTimes())
	// Assert no targets are notified
	assert.Equal(0, smsNotifier.calls)
	assert.Equal(0, mailNotifier.calls)
}

func TestUseCase4(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	smsNotifier, mailNotifier, timer, pager := arrangeUnhealthyDependencies()

	// Act
	now := 10
	newAlert := model.NewAlert(myServiceID, "a message", uint64(now))
	pager.SendAlert(newAlert)

	// Assert no timeout is sent to timer service
	assert.Equal(0, timer.CalledTimes())
	// Assert no targets are notified
	assert.Equal(0, smsNotifier.calls)
	assert.Equal(0, mailNotifier.calls)
}

func TestUseCase5(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	smsNotifier, mailNotifier, timer, pager := arrangeUnhealthyDependencies()

	// Act
	healthyEvent := model.NewHealthyEvent(myServiceID)
	pager.SendHealthyEvent(healthyEvent)

	// Assert monitored service is healthy
	status := pager.Status(myServiceID)
	assert.Equal(model.Ok, status.Health())
	// Assert no timeout is sent to timer service
	assert.Equal(0, timer.CalledTimes())
	// Assert no targets are notified
	assert.Equal(0, smsNotifier.calls)
	assert.Equal(0, mailNotifier.calls)
}

func arrangeUnhealthyDependencies() (*spySmsNotifier, *spyMailNotifier, *spyTimer, *domain.Pager) {
	ms := model.NewMonitoredService(myServiceID)
	alert := model.NewAlert(myServiceID, "a message", 10)
	ms.TurnToUnhealthy(alert)
	msr := &dummyMonitoredServiceRepository{ms: ms}
	epr := &dummyEscalationPolicyRepository{}
	smsNotifier := &spySmsNotifier{}
	mailNotifier := &spyMailNotifier{}
	notifier := service.NewCompositeNotifier(smsNotifier, mailNotifier)
	timer := &spyTimer{}
	timerConfiguration := &dummyTimerConfiguration{timeoutInSeconds: 15 * 60}
	pager := domain.NewPager(msr, epr, notifier, timer, timerConfiguration)
	return smsNotifier, mailNotifier, timer, pager
}
