package mongo

import (
	"github.com/dark-enstein/port/util"
	"strings"
)

type Unit struct {
	data string
}

func (u *Unit) IsValid() bool {
	return strings.ContainsAny(u.data, util.Forbidden)
}
