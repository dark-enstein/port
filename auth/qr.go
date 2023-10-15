package auth

import (
	"context"
	"github.com/dark-enstein/port/config"
	"github.com/dark-enstein/port/internal/generators"
	"github.com/dark-enstein/port/internal/generators/qr"
	"github.com/dark-enstein/port/util"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"strconv"
)

type QRDirector struct {
	cfg           *config.Config
	uid           string
	content       string
	size          int
	recoveryLevel int
	ctx           context.Context
	code          *qr.QR
	generator     generators.Generator
}

func (q *QRDirector) Generate() (string, error) {
	// TODO: hide the details of the QR package and only expose its functionality via the Generator interface
	return q.SetUp().code.Generate()
}

func (q *QRDirector) PingDependencies() (bool, error) {
	return pingDependencies()
}

func NewQRDirector(ctx context.Context, uid uuid.UUID, content, recoveryLevel string, size int, config *config.Config) *QRDirector {
	rec, _ := strconv.Atoi(recoveryLevel)
	return &QRDirector{
		ctx:           ctx,
		cfg:           config,
		recoveryLevel: rec,
		size:          size,
		content:       content,
		uid:           uid.String(),
	}
}

// SetUp sets up all the dependent structs and data, and readies director for execution
func (q *QRDirector) SetUp() *QRDirector {
	log := util.RetrieveLoggerFromCtx(q.ctx).WithMethod("QRDirector.SetUp()")
	//if !q.IsEmpty() {
	//	log.Debug().Msgf("director not empty with %v. SetUp() likely already called", q)
	//	return q
	//}
	q.code = qr.NewQRWithArgs(q.ctx, q.uid, q.content, q.size, qrcode.RecoveryLevel(q.recoveryLevel))
	log.Debug().Msg("setting up qr struct")
	log.Debug().Msgf("qr struct setup complete: %v", q.code)
	return q
}

// IsEmpty checks if QRDirector is empty, in other words, if SetUp() has been called.
// IsEmpty implements the Director interface
func (q *QRDirector) IsEmpty() bool {
	if q.uid == "" && q.content == "" && q.size == 0 && q.recoveryLevel == 1 && q.generator == nil && q.code == nil {
		return true
	}
	return false
}

func pingDependencies() (bool, error) {
	// TODO: later
	return true, nil
}
