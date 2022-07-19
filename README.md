# test-pager

## Requirements
- Docker (Docker version 20.10.14, build a224086)
- GNU Make 3.81

## How to run the tests
> make build && make test

## Architecture
I've tried to follow a ports and adapters approach to create the domain. According to the diagram shown in the [test description](https://github.com/aircall/technical-test-pager#problem), there are several adapters that interact with the Pager's domain and they are modelled through interfaces. In the next section, I'm describing each file content and how it maps the architecture.

## Folder content
- escalationpolicy/escalationpolicy.go. It contains the definition of all the entities needed to model the **_EscalationPolicy_** entity. This entity is retrieved through the **__escalationpolicy.Repository_**

- monitoredservice/alert.go. It contains the definition of the **_Alert_** object value. The state of this object value is managed by the **_MonitoredService_** entity, so it may suffer a little bit of "anemic model"

- monitoredservice/healthyevent.go. The value object **_HealthyEvent_** is defined here. It's sent to the domain when an alarm is solved. This value object is probably not too meaningful (it could be replaced by a simple string), but I think it makes the code more readable

- monitoredservice/monitoredservice.go. It contains the definition of the **_MonitoredService_** entity. I created this entity because it matches quite well the language used in the description of the problem. From a purely coding standpoint, it's a "wrapper" of the **_Alert_** object value, but the business logic handled by the **_MonitoredService_** makes a lot more sense to me to be in this entity than being handled by the **_Alert_** object value

- notifier/notifier.go. I've created a **_notifier.Composite_** whose dependencies are a **_notifier.Sms_** and a **_notifier.Mail_**, these latter are interfaces that are easily mockable ("spy-able") for testing purposes. These interfaces model the sms and mail adapters.

- timer/timer.go. It describes the interfaces (**_Timer_**) needed to communicate with the Timer external service.

- pager.go. This file offers the public API of the Pager service. It orchestrates the logic between all the adapters that represent the external dependencies alongside the entities, mainly **_MonitoredService_**. The API offered by the Pager service would be defined by the next use cases: SendAlert, AcknowledgeAlert, SendHealthyEvent, NotifyAckTimeout and Status. These are called by different primary/driving adapters (web console, alert system, timer system)

- pager_test.go. It contains the unit tests, as you can imagine. There is a test for each use case scenario defined in the [test description](https://github.com/aircall/technical-test-pager#use-case-scenarios)

- mocks_test.go. It contains the test doubles needed to execute the tests. There are basically spies to assert that some adapters were called and dummyRepositories that return the entities in a specific state suitable for that use case scenario

## Concurrency issues
When two alerts from the same service are sent at the same time, we could face data race conditions at a service level. There are two different scenarios: the Pager service is running in a single server or it's running in a distributed fashion (e.g. a fleet of servers, serverless service like AWS Lambda).

For the former scenario, we could assume that the single machine can run concurrent threads (either OS threads or other form of lightweight threads; in Go, this is performed through goroutines). In this case, I'd go for identifying the critical sections and then I'd use the mutual exclusion mechanisms to ensure atomicity.

For the latter, I'd go for a distributed lock pattern. When the first alarm reaches the Pager's service, it would acquire a lock (this could be implemented in [Redis](https://redis.io/docs/reference/patterns/distributed-locks/), for example). So when the second alarm is sent, the lock is already acquired and therefore the alarm is discarded. It's also worthwhile to mention the unit-of-work pattern to perform business transactions. In this way, the unit-of-work takes the responsibility of tracking the different changes of the entities during the transaction and commits them at the end of the request. In our scenario, the SendAlert use case would modify the MonitoredService state so that the second Alert received would wait to be processed. Since SendAlert method is idempotent, the second alert wouldn't change the service's state.