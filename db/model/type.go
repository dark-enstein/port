package model

import (
	"time"
)

type Unit interface {
	Kind() string
	GetTime() time.Time
	FindFilter() interface{}
}
