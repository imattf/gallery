package models

import (
  "fmt"
  "testing"
  "time"
)

func testingUserService() (*UserService, error) {
  const (
  	host   = "localhost"
  	port   = 5432
  	user   = "matthew"
  	dbname = "lenslocked_test"
  )
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  us, err := NewUserService(psqlInfo)
  if err != nil {
    return nil, err
  }
  // Toggle logging
  us.db.LogMode(false)
  //Clear user table between tests
  us.DestructiveReset()
  return us, nil
}

func TestCreateUser(t *testing.T) {
  us, err := testingUserService()
  if err != nil {
    t.Fatal(err)
  }
  user := User{
		Name: "Mike Scott",
		Email: "m@aol.io",
	}
  err = us.Create(&user)
  if err != nil {
    t.Fatal(err)
  }
  if user.ID == 0 {
    t.Errorf("Expected ID > 0. Recieved %d", user.ID)
  }
  if time.Since(user.CreatedAt) > time.Duration(5*time.Second) {
    t.Errorf("Expected CreactedAt to be recent. Recieved %s", user.CreatedAt)
  }
  if time.Since(user.UpdatedAt) > time.Duration(5*time.Second) {
    t.Errorf("Expected UpdatedAt to be recent. Recieved %s", user.UpdatedAt)
  }
}
