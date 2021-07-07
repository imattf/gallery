package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a record cannot be found
	ErrNotFound  = errors.New("models: resource not found")
	// ErrInvalidID is returned when an ID is 0, for example
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
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

// Lookup a user by Age in the database
func (us *UserService) ByAge(age uint) (*User, error) {
	var user User
	db := us.db.Where("age = ?", age)
	err := first(db, &user)
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
	// return nil   // use for unit test
	return us.db.Create(user).Error
}

// Update a user in the database
func (us *UserService) Update(user *User) error {
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

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Age   uint
}
