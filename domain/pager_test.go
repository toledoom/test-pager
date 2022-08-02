package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toledoom/test-pager/v2/domain"
	"github.com/toledoom/test-pager/v2/domain/monitoredservice"
	"github.com/toledoom/test-pager/v2/domain/notifier"
)

const myServiceID = "my-service-id"
const alertType1 = "1"
const alertType2 = "2"

func TestUseCase1(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	smsNotifier, mailNotifier, timer, pager := arrangeHealthyDependencies()

	// Act
	now := 10
	alert := monitoredservice.NewAlert(myServiceID, alertType1, "a message", uint64(now))
	pager.SendAlert(alert)

	// Assert monitored service is healthy
	status := pager.Status(myServiceID)
	assert.Equal(monitoredservice.Unhealthy, status.Health())
	// Assert timeout set to timer service
	assert.Equal(1, timer.CalledTimes())
	// Assert all sms and mail targets have been notified
	assert.Equal(1, smsNotifier.calls)
	assert.Equal(1, mailNotifier.calls)
}

func TestUseCase2(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	smsNotifier, mailNotifier, timer, pager := arrangeUnhealthyDependencies()

	// Act
	pager.NotifyAckTimeout(myServiceID, alertType1)

	// Assert timeout set to timer service
	assert.Equal(1, timer.CalledTimes())
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
	pager.AcknowledgeAlert(myServiceID, alertType1, acknowledgeAt)
	pager.NotifyAckTimeout(myServiceID, alertType1)

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
	newAlert := monitoredservice.NewAlert(myServiceID, alertType1, "a message", uint64(now))
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
	healthyEvent := monitoredservice.NewHealthyEvent(myServiceID, alertType1)
	pager.SendHealthyEvent(healthyEvent)

	// Assert monitored service is healthy
	status := pager.Status(myServiceID)
	assert.Equal(monitoredservice.Ok, status.Health())
	// Assert no timeout is sent to timer service
	assert.Equal(0, timer.CalledTimes())
	// Assert no targets are notified
	assert.Equal(0, smsNotifier.calls)
	assert.Equal(0, mailNotifier.calls)
}

func TestUseCaseTwoAlertsSameTypeOnlyOneIsNotified(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	smsNotifier, mailNotifier, timer, pager := arrangeHealthyDependencies()

	// Act
	now := 10
	alert := monitoredservice.NewAlert(myServiceID, alertType1, "a message", uint64(now))
	pager.SendAlert(alert)
	pager.SendAlert(alert)

	// Assert monitored service is unhealthy
	status := pager.Status(myServiceID)
	assert.Equal(monitoredservice.Unhealthy, status.Health())
	// Assert timeout set to timer service only once
	assert.Equal(1, timer.CalledTimes())
	// Assert all sms and mail targets have been notified only once
	assert.Equal(1, smsNotifier.calls)
	assert.Equal(1, mailNotifier.calls)
}

func TestUseCaseTwoAlertsDifferentTypeBothAreNotified(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	smsNotifier, mailNotifier, timer, pager := arrangeHealthyDependencies()

	// Act
	now := 10
	alert := monitoredservice.NewAlert(myServiceID, alertType1, "a message", uint64(now))
	pager.SendAlert(alert)
	alert2 := monitoredservice.NewAlert(myServiceID, alertType2, "another message", uint64(now))
	pager.SendAlert(alert2)

	// Assert monitored service is unhealthy
	status := pager.Status(myServiceID)
	assert.Equal(monitoredservice.Unhealthy, status.Health())
	// Assert timeout set to timer service twice
	assert.Equal(2, timer.CalledTimes())
	// Assert all sms and mail targets have been notified twice
	assert.Equal(2, smsNotifier.calls)
	assert.Equal(2, mailNotifier.calls)
}

func arrangeHealthyDependencies() (*spySmsNotifier, *spyMailNotifier, *spyTimer, *domain.Pager) {
	msr := &dummyMonitoredServiceRepository{}
	epr := &dummyEscalationPolicyRepository{}
	smsNotifier := &spySmsNotifier{}
	mailNotifier := &spyMailNotifier{}
	notifier := notifier.NewComposite(smsNotifier, mailNotifier)
	timer := &spyTimer{}
	timerConfiguration := &dummyTimerConfiguration{timeoutInSeconds: 15 * 60}
	pager := domain.NewPager(msr, epr, notifier, timer, timerConfiguration)
	return smsNotifier, mailNotifier, timer, pager
}

func arrangeUnhealthyDependencies() (*spySmsNotifier, *spyMailNotifier, *spyTimer, *domain.Pager) {
	ms := monitoredservice.New(myServiceID)
	alert := monitoredservice.NewAlert(myServiceID, alertType1, "a message", 10)
	ms.TurnToUnhealthy(alert)
	msr := &dummyMonitoredServiceRepository{ms: ms}
	epr := &dummyEscalationPolicyRepository{}
	smsNotifier := &spySmsNotifier{}
	mailNotifier := &spyMailNotifier{}
	notifier := notifier.NewComposite(smsNotifier, mailNotifier)
	timer := &spyTimer{}
	timerConfiguration := &dummyTimerConfiguration{timeoutInSeconds: 15 * 60}
	pager := domain.NewPager(msr, epr, notifier, timer, timerConfiguration)
	return smsNotifier, mailNotifier, timer, pager
}
