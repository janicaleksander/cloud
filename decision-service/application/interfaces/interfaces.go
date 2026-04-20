package interfaces

type DecisionPublisher interface {
	Publish(exchange string, message any) error
}
