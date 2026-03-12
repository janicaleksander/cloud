package event

import "github.com/google/uuid"

type Event interface {
	Type1Event | Type2Event | Type3Event | Type4Event
}
type Type1Event struct {
	Id string `json:"id"`
}

func NewType1Event() Type1Event {
	id, _ := uuid.NewRandom()
	return Type1Event{
		Id: id.String(),
	}
}

type Type2Event struct {
	Id string `json:"id"`
}

func NewType2Event() Type2Event {
	id, _ := uuid.NewRandom()
	return Type2Event{
		Id: id.String(),
	}
}

type Type3Event struct {
	Id string `json:"id"`
}

func NewType3Event() Type3Event {
	id, _ := uuid.NewRandom()
	return Type3Event{
		Id: id.String(),
	}
}

type Type4Event struct {
	Id string `json:"id"`
}

func NewType4Event() Type4Event {
	id, _ := uuid.NewRandom()
	return Type4Event{
		Id: id.String(),
	}
}
