package amazon

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dark-enstein/port/util"
)

var (
	DefaultProfile = "elvis"
)

const (
	SessionInContext = "session"
	DefaultBucket    = "port-elvis-gargantuan-panda" // https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
)

const (
	S3E = iota
)

type Compose struct {
	credLoc string
	action  *Action
}

type Action struct {
	// Name of AWS Service
	service int
	// Verb denoting action to be performed on service
	verb int
	// not sure what this does yet. Removing when next I see this comment.
	target string
}

type Response struct {
	Data []*s3.Bucket
	Name string
	Err  error
}

type S3UploadFileResponse struct {
	URL string
	Err error
}

func NewCompose(credLoc string, service int, verb int) *Compose {
	return &Compose{
		credLoc: credLoc,
		action: &Action{
			service: service,
			verb:    verb,
			target:  "target",
		},
	}
}

func (c *Compose) NewSessionWithOptions(ctx context.Context) (*Interaction, error) {
	alog := util.RetrieveLoggerFromCtx(ctx).WithMethod("NewSessionWithOptions()")
	sess, err := session.NewSessionWithOptions(session.Options{Profile: DefaultProfile, Config: aws.Config{
		Region: aws.String("us-west-2"),
	}})
	if err != nil {
		alog.Error().Msgf("creating session with aws failed with %w", err)
		return nil, fmt.Errorf("creating session with aws failed with %w", err)
	}
	alog.Debug().Msg("creating session with aws successful")
	return &Interaction{kind: "aws", session: sess}, nil
}

func (c *Compose) Service() int {
	return c.action.service
}

func (c *Compose) Verb() int {
	return c.action.verb
}

type Interaction struct {
	kind    string
	compose *Compose
	session *session.Session
}

func (i *Interaction) Kind() string {
	return i.kind
}

func (i *Interaction) Do(ctx context.Context, srv, verb int) *Response {
	log := util.RetrieveLoggerFromCtx(ctx).WithMethod("Interaction.Do()")
	ctx = context.WithValue(ctx, SessionInContext, i.session)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	switch srv {
	case S3E:
		sss := NewS3(ctx)
		res := make(chan *Response, 1)
		log.Debug().Msgf("created new S3 session: %v", sss)
		switch verb {
		case util.CREATE:
		case util.READ:
		case util.LIST:
			log.Info().Msg("list verb called on S3")
			go sss.List(ctx, res)
			select {
			case resp := <-res:
				return resp
			case <-ctx.Done():
				return &Response{
					Err: ctx.Err(),
				}
			}
		case util.UPDATE:
		case util.DELETE:
		case util.UPLOAD:
			uploadResp := make(chan *S3UploadFileResponse, 1)
			log.Info().Msg("upload verb called on S3")
			go sss.Upload(ctx, i.session, uploadResp)
			select {
			case resp := <-uploadResp:
				return &Response{
					Name: resp.URL,
					Err:  resp.Err,
				}
			case <-ctx.Done():
				return &Response{
					Err: ctx.Err(),
				}
			}
		}
	}
	// create bucket if not exist
	// C: upload
	// R: list
	// U: update
	// D: delete
	return nil
}

type S3 struct {
	// id of the file to be created. it is the request id
	id string
	// store is the connection to s3
	store *s3.S3
	// loc is the location of the file on disk to be uploaded
	loc string
}

func NewS3(ctx context.Context) *S3 {
	return &S3{
		id:    util.RetrieveReqIDFromCtx(ctx),
		store: s3.New(util.RetrieveFromCtx(ctx, SessionInContext).(*session.Session)),
		loc:   util.RetrieveFromCtx(ctx, util.QRLocInContext).(string),
	}
}

func (s *S3) List(ctx context.Context, response chan *Response) {
	log := util.RetrieveLoggerFromCtx(ctx).WithMethod("S3.List()")
	buckets, err := s.store.ListBuckets(nil)
	if err != nil {
		log.Error().Err(fmt.Errorf("encountered error: %w while trying to list s3 buckets", err))
		response <- &Response{
			Data: nil,
			Err:  err,
		}
	}
	log.Debug().Msgf("successfully listed s3 buckets: %#v", buckets)
	response <- &Response{
		Data: buckets.Buckets,
		Err:  err,
	}
}

func (s *S3) listBucket(ctx context.Context) *Response {
	log := util.RetrieveLoggerFromCtx(ctx).WithMethod("S3.listBucket()")
	log.Debug().Msg("handler")
	resp := make(chan *Response, 1)
	go s.List(ctx, resp)
	cliResp := new(Response)
	select {
	case cliResp := <-resp:
		return cliResp
	case <-ctx.Done():
		cliResp.Err = ctx.Err()
		return cliResp
	}
}

func (s *S3) Upload(ctx context.Context, sess *session.Session, response chan *S3UploadFileResponse) {
	log := util.RetrieveLoggerFromCtx(ctx).WithMethod("S3.Upload()")
	awsDefaultBucket := aws.String(DefaultBucket)

	// get the list of all the buckets in the account configured in the credentials
	listResp := s.listBucket(ctx)
	if listResp.Err != nil {
		log.Error().Err(fmt.Errorf("received err from S3.List handler: %w", listResp.Err))
		response <- &S3UploadFileResponse{
			Err: listResp.Err,
		}
	}

	// check if specified bucket is created already
	//var bucket *s3.Bucket TODO: Enable Bucket CLS and uploaded object public access
	found := false
	for i := 0; i < len(listResp.Data); i++ {
		if *listResp.Data[i].Name == DefaultBucket {
			found = !found
			//bucket = listResp.Data[i]
		}
	}

	if !found {
		log.Debug().Msg("bucket not found, proceeding to create one")

		inp := &s3.CreateBucketInput{
			Bucket: awsDefaultBucket,
		}
		_, err := s.store.CreateBucket(inp)
		if err != nil {
			log.Error().Err(fmt.Errorf("encountered error: %w while trying to create s3 buckets", err))
			response <- &S3UploadFileResponse{
				Err: err,
			}
		}

		err = s.store.WaitUntilBucketExists(&s3.HeadBucketInput{
			Bucket: awsDefaultBucket,
		})
		if err != nil {
			log.Error().Err(fmt.Errorf("encountered error: %w while waiting for bucket to get created", err))
			response <- &S3UploadFileResponse{
				Err: err,
			}
		}

		log.Debug().Msgf("bucket created: %s", DefaultBucket)
		//} else { TODO: Enable Bucket CLS and uploaded object public access
		//	//res, err := s.store.GetBucketAcl(&s3.GetBucketAclInput{Bucket: bucket.Name})
		//	//if err != nil {
		//	//	util.ExitOnErrorf("retrieving Bucket ACL failed with: %w", err)
		//	//}
		//	//res.Grants
		//	s.store.PutBucketAcl(&s3.PutBucketAclInput{Bucket: aws.String(DefaultBucket), ACL: aws.String("public-read")})
		//	s.store.PutBucketAcl()
		//}

		log.Debug().Msgf("preparing to upload to bucket: %s", DefaultBucket)
		uploader := s3manager.NewUploader(sess)

		// open file for upload
		log.Debug().Msgf("opening file to be uploaded: %s", DefaultBucket)
		file, err := os.Open(s.loc)
		if err != nil {
			log.Error().Err(fmt.Errorf("error while trying to open file %v for upload: %w", s.loc, err))
			response <- &S3UploadFileResponse{Err: err}
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Error().Err(fmt.Errorf("error while trying to close file: %v", s.loc))
			}
		}(file)

		// begin upload
		_, err = uploader.Upload(&s3manager.UploadInput{
			ACL:         aws.String("public-read"),
			Bucket:      awsDefaultBucket,
			Key:         aws.String(s.id),
			Body:        file,
			ContentType: aws.String("image/jpg"),
		})
		if err != nil {
			log.Error().Err(fmt.Errorf("error while trying to prepare file %v for upload: %w", s.loc, err))
			response <- &S3UploadFileResponse{Err: err}
		}
		log.Debug().Msgf("file %v uploaded successfully", s.loc)

		// retrieve the url to the file
		log.Debug().Msgf("retrieving url to uploaded object")
		req, _ := s3.New(sess).GetObjectRequest(&s3.GetObjectInput{Bucket: awsDefaultBucket, Key: aws.String(s.id)})
		rest.Build(req)
		urlLocation := req.HTTPRequest.URL.String()
		log.Debug().Msgf("file uploaded to %s", urlLocation)
		response <- &S3UploadFileResponse{
			URL: urlLocation,
			Err: err,
		}
	}
}

func (s *S3) CreateBucket(ctx context.Context, response chan *Response) {
	//log := util.RetrieveLoggerFromCtx(ctx).WithMethod("S3.CreateBucket()")
	//
	//inp := &s3.CreateBucketInput{
	//	Bucket: &s.id,
	//}
	//buckets, err := s.store.CreateBucket(inp)
	//if err != nil {
	//	log.Error().Err(fmt.Errorf("encountered error: %w while trying to list s3 buckets", err))
	//	response <- &Response{
	//		Data: nil,
	//		Err:  err,
	//	}
	//}
	//log.Debug().Msgf("successfully listed s3 buckets: %#v", buckets)
	//response <- &Response{
	//	Data: buckets.Location,
	//	Err:  err,
	//}
}

// TODO: create bucket function, create bucket if not exist function, upload bucket function
