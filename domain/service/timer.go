package service

type Timer interface {
	SetTimeout(timeoutInSeconds int)
}

type TimerConfiguration interface {
	GetTimeoutInSeconds() int
}
