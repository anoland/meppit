package domain

import (
  "time"
  "labix.org/v2/mgo/bson"
)

type Person struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
	Phone     string
	Timestamp time.Time
}

func NewPerson(name string, phone string) (person *Person) {
  return &Person{}
}