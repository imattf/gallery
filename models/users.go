package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/go-courses/lenslocked.com/hash"
	"gitlab.com/go-courses/lenslocked.com/rand"
)

var (
	// ErrNotFound is returned when a record cannot be found
	ErrNotFound  = errors.New("models: resource not found")

	// ErrInvalidID is returned when an ID is 0, for example
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password is used
	// when authenticating a user
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

const userPwPepper = "some-secret"
const hmacSecretKey = "secret-hmac-key"


type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember      string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{
		db:   db,
		hmac: hmac,
	}, nil
}

type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// ByID method allows us to find a user
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// Lookup a user by Email in the database
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// Lookup by Remember looks up user by Remeber token
// this method wil handle hasing the token for us
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := us.hmac.Hash(token)
	db := us.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// Lookup a users by Age Range in the database
func (us *UserService) InAgeRange(minAge, maxAge uint) ([]User, error) {
	var users []User
	db := us.db.Where("age BETWEEN ? and ?", minAge, maxAge)
	err := find(db, &users)
	return users, err
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

//Creates a user in the database and will backfill
// related meta-data like ID, CreatedAt...
func (us *UserService) Create(user *User) error {
	pwBytes :=[]byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)
	return us.db.Create(user).Error
}


// Athenticates a user loging request
// takes an email and Password
// If the email doesn't exist
//   return nil and ErrNotFound
// If the password provided doesn't match the hased password
//   return nil and an ErrInvalidPassword
// If the email and password are both valid
//   return the user and nil
// Otherwise another system error was encountered
//   return nil and the error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err =bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password + userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	 return foundUser, nil

}

// Update a user in the database
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
	return us.db.Save(user).Error
}

// Delete a user in the database
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// Closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// Automigrate will attempt to automatically migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil

}
