package event

import "github.com/google/uuid"

const (
	TYPE1 = "TYPE1"
	TYPE2 = "TYPE2"
	TYPE3 = "TYPE3"
	TYPE4 = "TYPE4"
)

type Type1Event struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewType1Event() *Type1Event {
	id, _ := uuid.NewRandom()
	return &Type1Event{
		Id:   id.String(),
		Name: TYPE1,
	}
}

type Type2Event struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewType2Event() *Type2Event {
	id, _ := uuid.NewRandom()
	return &Type2Event{
		Id:   id.String(),
		Name: TYPE2,
	}
}

type Type3Event struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewType3Event() *Type3Event {
	id, _ := uuid.NewRandom()
	return &Type3Event{
		Id:   id.String(),
		Name: TYPE3,
	}
}

type Type4Event struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewType4Event() *Type4Event {
	id, _ := uuid.NewRandom()
	return &Type4Event{
		Id:   id.String(),
		Name: TYPE4,
	}
}
