package model

import (
	"context"
	"fmt"
	"time"
)

var (
	UnitUser = "user"
	glog     = &Logger{}
)

// User holds the user information, and it is ready for working with the DB
type User struct {
	Name  *Name    `bson:"name,inline"`
	Birth *string  `bson:"name"`
	Roles *RoleSet `bson:"role_set,omitempty"`
}

// UserOptions holds the user request options, and it is ready for working with the DB
type UserOptions struct {
	Database         string
	Collection       string
	Table            string
	CreateOnNotExist bool
}

func NewUserOptions() *UserOptions {
	return &UserOptions{}
}

func (u *UserOptions) IsValid() bool {
	//TODO
	return true
}

func (u *UserOptions) RetrieveDatabase() string {
	return u.Database
}

func (u *UserOptions) RetrieveTable() string {
	return u.Table
}

func (u *UserOptions) RetrieveCollection() string {
	return u.Collection
}

func (u *UserOptions) RetrieveOverride() bool {
	return u.CreateOnNotExist
}

func NewUser(ctx context.Context) *User {
	glog = RetrieveLoggerFromCtx(ctx)
	return &User{}
}

func NewUserComplete(ctx context.Context, name *Name, birthdate *string, rs *RoleSet) *User {
	u := NewUser(ctx)
	mlog := glog.WithMethod("NewUserComplete()")
	u.Name, u.Birth, u.Roles = name, birthdate, rs
	mlog.Info().Msgf("successfully created user model %v", u.Name)
	mlog.Debug().Str("created model user with params", fmt.Sprintf("%v", u)).Send()
	return u
}

func (u *User) WithName(name *Name) *User {
	mlog := glog.WithMethod("WithName()")
	u.Name = name
	mlog.Info().Msgf("added name %v to user model", u.Name)
	return u
}

func (u *User) WithBirthDate(birthdate *string) *User {
	mlog := glog.WithMethod("WithBirthDate()")
	u.Birth = birthdate
	mlog.Info().Msgf("added birthdate %v to user model", u.Birth)
	return u
}

func (u *User) WithRoleSet(rs *RoleSet) *User {
	mlog := glog.WithMethod("WithRoleSet()")
	u.Roles = rs
	mlog.Info().Msgf("added roleset %v to user model", u.Roles)
	return u
}

func (u *User) GetTime() time.Time {
	return time.Now()
}

func (u *User) Kind() string {
	return UnitUser
}

func (u *User) NameStr() string {
	return fmt.Sprintf("%v %v", u.Name.FirstName, u.Name.LastName)
}

type Name struct {
	FirstName string
	LastName  string
}

func (n *Name) GetTime() time.Time {
	return time.Now()
}

type RoleSet []map[string]*int64
