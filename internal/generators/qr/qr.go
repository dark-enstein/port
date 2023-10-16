package qr

import (
	"context"
	"fmt"
	"github.com/dark-enstein/port/config"
	amazon "github.com/dark-enstein/port/internal/cloud/aws"
	"github.com/dark-enstein/port/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/skip2/go-qrcode"
	"os"
	"path/filepath"
	"reflect"
)

type Type string

func (t *Type) String() string {
	ref := reflect.ValueOf(*t)
	if ref.Kind() != reflect.String {
		return ""
	}
	return ref.String()
}

var (
	Factory    = func(filename string) string { return filepath.Join(DefaultDir, filename) }
	GetFactory = func(filename string, log *zerolog.Logger) (*os.File, error) {
		if err := os.MkdirAll(DefaultDir, 0755); err != nil {
			log.Error().Err(fmt.Errorf("failed to create qr directory: %w", err))
			return nil, err
		}
		return os.OpenFile(Factory(filename), os.O_RDWR|os.O_CREATE, 0755)
	}
	DefaultFilename = uuid.New().String() + ".png"
	DefaultDir      = filepath.Join(".qr", "generated")
)

// QR defines the structure of a QRcode
type QR struct {
	id            string
	content       string
	size          int
	recoveryLevel qrcode.RecoveryLevel
	Code          *qrcode.QRCode
	ctx           context.Context
	uploadedLoc   string
}

// NewQR generates an empty QR. It receives no arguments, and returns a QR pointer defined by empty fields
func NewQR(ctx context.Context) *QR {
	return &QR{}
}

// NewQRWithArgs generates a QR defined by the arguments passed in.
func NewQRWithArgs(ctx context.Context, id, content string, size int, recoveryLevel qrcode.RecoveryLevel) *QR {
	return &QR{
		id:            id,
		content:       content,
		size:          size,
		recoveryLevel: recoveryLevel,
		Code:          &qrcode.QRCode{},
		ctx:           ctx,
	}
}

// WithID persists the passed in ID into the receiver QR.
func (q *QR) WithID(id string) *QR {
	q.id = id
	return q
}

// WithContent persists the passed in ID into the receiver QR.
func (q *QR) WithContent(content string) *QR {
	q.content = content
	return q
}

// WithSize persists the passed in ID into the receiver QR.
func (q *QR) WithSize(size int) *QR {
	q.size = size
	return q
}

// WithRecoveryLevel persists the passed in ID into the receiver QR.
func (q *QR) WithRecoveryLevel(rec qrcode.RecoveryLevel) *QR {
	q.recoveryLevel = rec
	return q
}

// Generate encodes the content from QR into a QRcode, and saves it on disk/or in buffer.
func (q *QR) Generate() (string, error) {
	return q.upload()
}

func (q *QR) upload() (string, error) {
	log := util.RetrieveLoggerFromCtx(q.ctx).WithMethod("Generate()")
	err := q.generate()
	if err != nil {
		return "", err
	}
	q.ctx = context.WithValue(q.ctx, util.QRLocInContext, q.uploadedLoc)

	comp := amazon.NewCompose(config.DefaultFlagLOC, amazon.S3E, util.UPLOAD)
	interact, err := comp.NewSessionWithOptions(q.ctx)
	if err != nil {
		log.Error().Err(fmt.Errorf("encountered error while trying to upload qr code: %w", err))
		return "", err
	}

	resp := interact.Do(q.ctx)
	if resp.Err == context.DeadlineExceeded {
		log.Error().Err(fmt.Errorf("file upload failed due to: %w", err))
		return "", context.DeadlineExceeded
	}

	return resp.Name, resp.Err
}

// generate encodes the content from QR into a QRcode, and saves it on disk/or in buffer.
func (q *QR) generate() error {
	log := util.RetrieveLoggerFromCtx(q.ctx).WithMethod("Generate()")
	var err error
	q.Code, err = qrcode.New(q.content, q.recoveryLevel)
	if err != nil {
		log.Error().Msgf("qrcode.New() failed with error: %w", err)
		return err
	}

	if q.id != "" {
		DefaultFilename = q.id + ".png"
	}
	factory, err := GetFactory(DefaultFilename, log)
	if err != nil {
		log.Error().Msgf("GetFactory() failed with error: %w", err.Error())
		return err
	}
	q.uploadedLoc = factory.Name()

	defer func(factory *os.File) {
		err = factory.Close()
		if err != nil {
			log.Error().Msgf("Closing open file failed with error: %w", err.Error())
			return
		}
	}(factory)

	err = q.Code.Write(q.size, factory)
	if err != nil {
		log.Error().Msgf("Writing qrcode image to file failed with error: %w", err)
		return err
	} //write to file and buffer

	return nil
}
