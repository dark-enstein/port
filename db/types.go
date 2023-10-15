package db

import (
	"context"
	"github.com/dark-enstein/port/db/model"
)

type DBType string

type DB interface {
	Ping() bool

	Kind() string
	Host() string

	Create(context.Context, model.Unit, model.Opts) *model.DBResponse
	CreateAll(context.Context, []model.Unit, model.CreateOpts) *model.DBResponse
	//Read(ReadOpts) *DBResponse
	//Update(UpdateOpts) *DBResponse
	//Delete(DeleteOpts) *DBResponse

	// Ensure the CRUD dependents is all set up, including databases, collections, tables, etc.
	// This is DB engine specific. The override flag is used to decide if the missing scaffold chould be created or not
	EnsureDBScaffold(ctx context.Context, override bool) error
}

type Unit interface {
	String() bool
}
