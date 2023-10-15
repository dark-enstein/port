package model

import "github.com/dark-enstein/port/util"

var (
	Tables         = map[string]string{}
	UserDB         = "users"
	UserCollection = "user-info"
)

type DBResponse struct {
	Err      error
	ID       string
	Metadata Metadata
}

type Metadata interface {
	String() string
}

type CreateOpts struct {
	TargetTable string
}

func (c *CreateOpts) IsValid() bool {
	if !util.IsIn(c.TargetTable, FlattenTables()) {
		return false
	}
	return true
}

func LoadTablesInCache() {
	Tables = map[string]string{
		"users": UserDB,
	}
}

func FlattenTables() []string {
	return util.FlattenMapToString(Tables)
}

type Opts interface {
	IsValid() bool
	RetrieveDatabase() string
	RetrieveTable() string
	RetrieveCollection() string
	RetrieveOverride() bool
}
