# test-pager

## Requirements
- Docker (Docker version 20.10.14, build a224086)
- GNU Make 3.81

## How to run the tests
> make build && make test

## Architecture
I've tried to follow a ports and adapters approach to create the domain. According to the diagram shown in the [test description](https://github.com/aircall/technical-test-pager#problem), there are several adapters that interact with the Pager's domain and they are modelled through interfaces. In the next section, I'm describing each file content and how it mmaps the architecture.

## Folder content

1. Folder model. It contains the entities' and value objects' definitions and the repositories that will return those entities.

    - model/alert.go. It contains the definition of the **_Alert_** object value
    
    - model/escalationpolicy.go. It contains the definition of all the entities needed to model the **_EscalationPolicy_** entity. This entity is retrieved through the **__EscalationPolicyRepository_**

    - model/healthyevent.go. The value object **_HealthyEvent_** is defined here. It's sent to the domain when an alarm is solved. This value object is probably not too meaningful (it could be replaced by a simple string), but I think it makes the code more readable

    - model/monitoredservice.go. It contains the definition of the **_MonitoredService_** entity. I created this entity because it matches quite well the language used in the description of the problem. From a modelling standpoint, it's a "wrapper" of the **_Alert_** object value, but the business logic handled by the **_MonitoredService_** makes a lot more sense to me than being handled by the **_Alert_** entity.

2. Folder service. It defines those services that do not have any entity linked to them.

    - service/notifier.go. I've created a **_CompositeNotifier_** whose dependencies are a **_SmsNotifier_** and a **_MailNotifier_**, these latter are interfaces that are easily mockable ("spy-able") for testing purposes. These interfaces model the sms and mail adapters.

    - service/timer.go. It describes the interfaces (**_Timer_**) needed to communicate with the Timer external service.

3. File pager.go. This file offers the public API of the Pager service. It orchestrates the logic between all the adapters defined by the services and the repositories that represent the external dependencies (adapters) alongside the entities, mainly **_MonitoredService_**.

4. File pager_test.go. It contains the tests, as you can imagine. There is a test for each use case scenario defined in the [test description](https://github.com/aircall/technical-test-pager#use-case-scenarios)

5. File mocks_test.go. It contains the test doubles needed to execute the tests.

## Concurrency issues
- Unit of work: business transactions
- [Distributed locks](https://redis.io/docs/reference/patterns/distributed-locks/)