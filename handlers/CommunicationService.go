package handlers

type Communication struct {
	IMService interface{}
}

type IMExchange interface {
	PushedUpdates()
	Notify()
}
