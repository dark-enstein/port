package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dark-enstein/port/config"
	"github.com/dark-enstein/port/db"
	"github.com/dark-enstein/port/db/model"
	"github.com/rs/zerolog"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	VanillaUser    = Role{}
	Administrator  = RoleSet{}
	Developer      = RoleSet{}
	ulog           = config.NewLogger(zerolog.GlobalLevel().String())
	KindUser       = "user"
	UserDB         = config.DefaultDBName
	UserTable      = "users"
	UserCollection = "users"
)

// RoleSet defines a list of roles that can be attached to an application user; eg PowerUser on AWS, which consist EC2 Admin, IAM Admin, etc
type RoleSet []Role

// ToBinary returns the binary representation of a RoleSet, where 1 == true, 0 == false
func (rs *RoleSet) ToBinary() []map[string]*int64 {
	res := make([]map[string]*int64, len(*rs))
	for _, v := range *rs {
		res = append(res, v.ToBinary())
	}
	return res
}

// Role defines a list of permissions that can be attached to an application user; eg IAM Administrator role, EC2 creator permission
type Role struct {
	User PermissionSet
}

// ToBinary returns the binary representation of a role, where 1 == true, 0 == false
func (r *Role) ToBinary() map[string]*int64 {
	res := make(map[string]*int64, len(r.User))
	for _, v := range r.User {
		res[v.Name] = v.ToBinary()
	}
	return res
}

// PermissionSet is a list of permissions present in a role;
type PermissionSet []Permission

// ToBinary converts a permission set struct into it's binary representation, where 1 == true, 0 == false
func (ps *PermissionSet) ToBinary() map[string]*int64 {
	res := make(map[string]*int64, len(*ps))
	for _, v := range *ps {
		res[v.Name] = v.ToBinary()
	}
	return res
}

// Permission is a structure for a permission unit. This defines the permissions on a particular resource. It consist the name of the resource, and the permissions boolean
type Permission struct {
	Name   string
	Create bool
	Read   bool
	Update bool
	Delete bool
}

// ToBinary converts a permission struct into it's binary representation, where 1 == true, 0 == false
func (p *Permission) ToBinary() *int64 {
	bin := ""
	if p.Create {
		bin += "1"
	} else {
		bin += "0"
	}

	if p.Read {
		bin += "1"
	} else {
		bin += "0"
	}

	if p.Update {
		bin += "1"
	} else {
		bin += "0"
	}

	if p.Delete {
		bin += "1"
	} else {
		bin += "0"
	}

	if bin == "" {
		b := int64(0)
		return &b
	}

	binInt, _ := strconv.ParseInt(bin, 2, 64)
	return &binInt
}

type User struct {
	Name     string `json:"name"`
	Birth    string `json:"dob"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// InternalUser struct for performing various computing before making it ready for the DB
type InternalUser struct {
	name     Name
	birth    DateOfBirth
	username string
	password string
}

func (u InternalUser) Name() *Name {
	return &u.name
}

func (u InternalUser) NameStr() string {
	return u.name.String()
}

func (u InternalUser) BirthDate() *string {
	str := u.birth.String()
	return &str
}

func (u InternalUser) Kind() string {
	return KindUser
}

func (u InternalUser) GetTime() time.Time {
	return time.Now()
}

func (u InternalUser) GetPermissions() PermissionSet {
	// db calls
	return PermissionSet{}
}

func (u InternalUser) GetRoles() RoleSet {
	// db calls
	return RoleSet{}
}

func (u InternalUser) IntoUserModel(ctx context.Context) *model.User {
	log := GetLoggerFromCtx(ctx).With().Str("method", "InternalUser.IntoUserModel()").Logger()
	hash, err := hashPassword(u.password)
	if err != nil {
		log.Error().Err(fmt.Errorf("unable to hash password: %w", err))
		return nil
	}
	mU := model.NewUser(ctx).WithName(
		&model.Name{FirstName: u.Name().RetrieveFirstName(), LastName: u.Name().RetrieveLastName()}).WithBirthDate(
		u.BirthDate()).WithRoleSet(
		&model.RoleSet{VanillaUser.ToBinary()}).WithUsername(u.username).WithPasswordHash(hash)
	log.Debug().Msgf("converted Internal user %v into Model user %v successfully", u, mU)
	return mU
}

// hashPassword hashes the password
func hashPassword(orig string, key ...string) (string, error) {
	keyToBeUsed := []byte(config.DefaultSecretKey)
	if len(key) != 0 {
		keyToBeUsed = []byte(key[0])
	}
	block, err := aes.NewCipher(keyToBeUsed)
	if err != nil {
		return "", err
	}
	b := base64.StdEncoding.EncodeToString([]byte(orig))
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return string(ciphertext), nil
}

// deHashPassword hashes the password
func deHashPassword(orig string, key ...string) (string, error) {
	keyToBeUsed := []byte(config.DefaultSecretKey)
	if len(key) != 0 {
		keyToBeUsed = []byte(key[0])
	}
	origByte := []byte(orig)
	block, err := aes.NewCipher(keyToBeUsed)
	if err != nil {
		return "", err
	}
	if len(origByte) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := origByte[:aes.BlockSize]
	origByte = origByte[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(origByte, origByte)
	data, err := base64.StdEncoding.DecodeString(string(origByte))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Name struct {
	firstName string
	lastName  string
}

func NewName(fullname string) *Name {
	ulog := ulog.With().Str("method", "NewName()").Logger()
	name := strings.Split(fullname, " ")
	if len(name) != 2 {
		ulog.Info().Msgf("name field %v contains less than two strings", fullname)
	}
	return &Name{
		firstName: name[0],
		lastName:  name[1],
	}
}

func (n *Name) String() string {
	return fmt.Sprintf("%v %v", n.firstName, n.lastName)
}

func (n *Name) RetrieveFirstName() string {
	return n.firstName
}

func (n *Name) RetrieveLastName() string {
	return n.lastName
}

type DateOfBirth struct {
	y string
	m string
	d string
}

// NewDateOfBirth creates a new date object taking in a date argument in the format (DD/MM/YYYY).
func NewDateOfBirth(date string) *DateOfBirth {
	ulog := ulog.With().Str("method", "NewDateOfBirth()").Logger()
	dateSlice := strings.Split(date, "/")
	if len(dateSlice) != 3 {
		ulog.Info().Msgf("date field %v doesnt conform with standard (DD/MM/YYYY)", date)
	}
	return &DateOfBirth{
		y: dateSlice[2],
		m: dateSlice[1],
		d: dateSlice[0],
	}
}

func (d *DateOfBirth) Day() string {
	return d.d
}

func (d *DateOfBirth) Month() string {
	return d.m
}

func (d *DateOfBirth) Year() string {
	return d.y
}

func (d *DateOfBirth) String() string {
	return fmt.Sprintf("%v/%v/%v", d.d, d.m, d.y)
}

func NewUser() *User {
	return &User{}
}

func (u *User) String() string {
	return fmt.Sprintf("%s %s", u.Name, u.Birth)
}

func (u *User) IntoInternal() *InternalUser {
	return &InternalUser{
		name:  *NewName(u.Name),
		birth: *NewDateOfBirth(u.Birth),
	}
}

// UserDirector defines a master that can perform major resource creation functions in the Auth package
type UserDirector struct {
	log    *zerolog.Logger
	ReqCtx context.Context
	users  map[InternalUser]*model.UserOptions
	mUsers []model.User
	sync.Mutex
	db db.DB
}

func NewUserDirector(ctx context.Context) *UserDirector {
	return &UserDirector{ReqCtx: ctx, db: GetDBFromCtx(ctx), log: GetLoggerFromCtx(ctx)}
}

func (d *UserDirector) Create() ([]string, map[string]error) {
	log := d.log.With().Str("method", "UserDirector.Create()").Logger()
	cantCreate := make(map[string]error, len(d.users))
	createdIDs := make([]string, len(d.users))
	for k, v := range d.users {
		u := k.IntoUserModel(d.ReqCtx)
		opts := v
		dbResp := d.db.Create(d.ReqCtx, u, opts)
		if dbResp.Err != nil {
			cantCreate[u.NameStr()] = dbResp.Err
			log.Info().Msgf("cannot create the user %v due to error: %v. \ncontinuing..", u.NameStr(), dbResp.Err)
		} else {
			createdIDs = append(createdIDs, dbResp.ID)
		}
	}
	return createdIDs, cantCreate
}

//// CreateAll creates multiple users in the Database
//func (d *UserDirector) CreateAll() ([]byte, error) {
//	_ = db.CreateOpts{TargetTable: db.UserDB}
//	for i := 0; i < len(d.users); i++ {
//		dbResp := d.db.Create(d.users[i], opts)
//	}
//	return nil, nil
//}

// CreateUsers creates a group of users in the DB
// Errors that come through this method will be processing errors, not data errors.
// Ideally all data errors would be captured early during the request validation process
func (d *UserDirector) CreateUsers(u []InternalUser) func() ([]string, error) {

	return func() ([]string, error) {
		log := d.log.With().Str("method", "UserDirector.CreateUsers()").Logger()
		d.users = make(map[InternalUser]*model.UserOptions, len(u))
		for i := 0; i < len(u); i++ {
			d.users[u[i]] = resolveOpts(u[i].Kind()).(*model.UserOptions)
		}
		createIDs, errGrp := d.Create()
		if len(errGrp) > 1 {
			log.Info().Msgf("errors while creating users: %v", errGrp)
		}
		return createIDs, nil // these errors are server errors, or for tracing, not one to be returned to the client
	}
}
