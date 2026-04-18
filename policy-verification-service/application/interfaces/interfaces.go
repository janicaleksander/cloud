package interfaces

type PolicyEventPublisher interface {
	Publish(exchange string, msg interface{}) error
}
