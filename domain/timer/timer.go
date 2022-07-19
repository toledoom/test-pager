package timer

type Timer interface {
	SetTimeout(serviceID string, timeoutInSeconds int)
}

type Configuration interface {
	GetTimeoutInSeconds() int
}
