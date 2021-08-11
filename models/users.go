package models

import (
	"errors"
	"strings"
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/go-courses/lenslocked.com/hash"
	"gitlab.com/go-courses/lenslocked.com/rand"
)

var (
	// ErrNotFound is returned when a record cannot be found
	ErrNotFound  = errors.New("models: resource not found")

	// ErrIDInvalid is returned when an ID is 0, for example
	ErrIDInvalid = errors.New("models: ID provided was invalid")

	// ErrPasswordIncorrect is returned when an invalid password is used
	// when authenticating a user
	ErrPasswordIncorrect = errors.New("models: incorrect password provided")

	// ErrEmailRequired is returned when an email address is not provided
	// when creating a user
	ErrEmailRequired = errors.New("models: Email address is required")

	// ErrEmailInvalid is returned when an email is not properly formatted
	ErrEmailInvalid = errors.New("models: Email address is in valid")

	// emailRegex is used to match email address aligned with top level domains
	// of 2 to 16 characters in length, always alfa chars only.
	// emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)

	// ErrEmailTaken is used to verify if an email is already in use during an update
	// or create of a user
	ErrEmailTaken = errors.New("models: Email address is already taken")

	// ErrPasswordRequired is return when creating a user and no password
	// is provided
	ErrPasswordRequired = errors.New("models: Password is required")

	// ErrPasswordTooShort is used to insure password has minimum length
	ErrPasswordTooShort = errors.New("models: Password is too short")

	// ErrRememberTooShort is used to insure remember token is at least 32 bytes
	ErrRememberTooShort = errors.New("models: Remember token must be 32 bytes.")

	// ErrRememberHash is returned when a create or update
	// is attempted without a valid user remember token hash.
	ErrRememberRequired = errors.New("models: Remember hash is reuired.")
)

const userPwPepper = "some-secret"
const hmacSecretKey = "secret-hmac-key"

// User represents the user model stored in our database
// This is used for user accounts, storing both an email and password
// so users can log in and gain access to their content.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember      string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}


// UserDB interface is used for interacting with the users database.
//
// For prettry much all single user queries:
// If the user is found, we will return a nil error
// If the user is not found, we will return a ErrNotFound
// If there is another error, we will return an error with
// additional information about what went wrong. This may not be an error
// generated by the model package.
//
// For single user queries, any error but ErrNotFound should
// probably result in a 500 error.
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(id string) (*User, error)
	ByRemember(toke string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB Connection
	Close() error

	// Migration Helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods to manipulate and work with the user model
type UserService interface {
	// Authenticate will verify the provided email address and password
	// are correct. Tf they are coreect, the corresponding user will be
	// returned. Otherwiseyou will recieve either:
	// ErrNotFound, ErrPasswordIncorrect or another error if something goes wrong.
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(connectionInfo string) (UserService, error) {
  ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
  hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}, nil
}

// Compiler check to make sure userService implements UserService
var _ UserService = &userService{}

type userService struct {
  UserDB
}

// Authenticates a user login request
// takes an email and Password
// If the email doesn't exist
//   return nil and ErrNotFound
// If the password provided doesn't match the hased password
//   return nil and an ErrPasswordIncorrect
// If the email and password are both valid
//   return the user and nil
// Otherwise another system error was encountered
//   return nil and the error
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err =bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password + userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrPasswordIncorrect
		default:
			return nil, err
		}
	}
	 return foundUser, nil

}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}


// Compiler check to make sure userValidator implements UserDB
var _ UserDB = &userValidator{}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac: 	    hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
	emailRegex *regexp.Regexp
}

// ByEmail will normalize the email address before calling ByEmail on the
// UserDB field.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember will hash the remember token and then call
// ByRemember on the subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

//Creates a user in the database and will backfill
// related meta-data like ID, CreatedAt...
func (uv *userValidator) Create(user *User) error {
  err := runUserValFuncs(user,
    uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
	  uv.normalizeEmail,
	  uv.requireEmail,
	  uv.emailFormat,
	  uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Update will hash a remember token if it is provided.
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
    uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
	  uv.normalizeEmail,
	  uv.requireEmail,
	  uv.emailFormat,
	  uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete a user in the database
func (uv *userValidator) Delete(id uint) error {
  var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// bcryptPassword will hash a user's password with a predefinded pepper
// (userPwPepper) and bcrypt if the Passwprd field is not the empty string.
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes :=[]byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

// uv method hmac
func (uv *userValidator) hmacRemember(user *User) error{
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error{
		if user.ID <= n {
			return ErrIDInvalid
		}
		return nil
	})
}

func (uv *userValidator) normalizeEmail(user *User) error {
  user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		// Email address is not taken
		return nil
	}
	if err != nil {
		return err
	}
  // We found a user w/ a email address...
	// If the found user has the same IS as this user, it is
	// an update of the same user's email address
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

// Compiler check that type matches interface
var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	// hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{
		db:   db,
		// hmac: hmac,
	}, nil
}

type userGorm struct {
	db   *gorm.DB
	// hmac hash.HMAC
}

// ByID method allows us to find a user
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// Lookup a user by Email in the database
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with a given remember token and returns
// that user. The method expects the remember token to already be hashed.
// Errors are the same as ByEmail.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	db := ug.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// Lookup a users by Age Range in the database
func (ug *userGorm) InAgeRange(minAge, maxAge uint) ([]User, error) {
	var users []User
	db := ug.db.Where("age BETWEEN ? and ?", minAge, maxAge)
	err := find(db, &users)
	return users, err
}

//Creates a user in the database and will backfill
// related meta-data like ID, CreatedAt...
func (ug *userGorm) Create(user *User) error {
	// pwBytes :=[]byte(user.Password + userPwPepper)
	// hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	// if err != nil {
	// 	return err
	// }
	// user.PasswordHash = string(hashedBytes)
	// user.Password = ""
	//
	// if user.Remember == "" {
	// 	token, err := rand.RememberToken()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	user.Remember = token
	// }
	// user.RememberHash = ug.hmac.Hash(user.Remember)
	return ug.db.Create(user).Error
}

// Update a user in the database
func (ug *userGorm) Update(user *User) error {
	// if user.Remember != "" {
	// 	user.RememberHash = ug.hmac.Hash(user.Remember)
	// }
	return ug.db.Save(user).Error
}

// Delete a user in the database
func (ug *userGorm) Delete(id uint) error {
	// if id == 0 {
	// 	return ErrIDInvalid
	// }
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Closes the UserService database connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset drops the user table and rebuilds it
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// Automigrate will attempt to automatically migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil

}

// first is a help function for Lookups and it will get the first item
// returned and place into destination.
// orig: func first(db *gorm.DB, user *User) error {
func first(db *gorm.DB, destination interface{}) error {
	err := db.First(destination).Error
	if err == gorm.ErrRecordNotFound {
	  return ErrNotFound
	}
	return err
}

// find is a help function for Lookups and it will get the first item
// returned and place into destination.
// orig: func first(db *gorm.DB, user *User) error {
func find(db *gorm.DB, destination interface{}) error {
	err := db.Find(destination).Error
	if err == gorm.ErrRecordNotFound {
	  return ErrNotFound
	}
	return err
}
