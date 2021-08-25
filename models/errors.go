package models

import "strings"

const (
	// ErrNotFound is returned when a record cannot be found
	ErrNotFound modelError = "models: resource not found"

	// ErrIDInvalid is returned when an ID is 0, for example
	ErrIDInvalid modelError = "models: ID provided was invalid"

	// ErrPasswordIncorrect is returned when an invalid password is used
	// when authenticating a user
	ErrPasswordIncorrect modelError = "models: incorrect password provided"

	// ErrEmailRequired is returned when an email address is not provided
	// when creating a user
	ErrEmailRequired modelError = "models: Email address is required"

	// ErrEmailInvalid is returned when an email is not properly formatted
	ErrEmailInvalid modelError = "models: Email address is in valid"

	// emailRegex is used to match email address aligned with top level domains
	// of 2 to 16 characters in length, always alfa chars only.
	// emailRegex modelError = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`

	// ErrEmailTaken is used to verify if an email is already in use during an update
	// or create of a user
	ErrEmailTaken modelError = "models: Email address is already taken"

	// ErrPasswordRequired is return when creating a user and no password
	// is provided
	ErrPasswordRequired modelError = "models: Password is required"

	// ErrPasswordTooShort is used to insure password has minimum length
	ErrPasswordTooShort modelError = "models: Password must be at least 8 characters long"

	// ErrRememberTooShort is used to insure remember token is at least 32 bytes
	ErrRememberTooShort modelError = "models: Remember token must be 32 bytes."

	// ErrRememberHash is returned when a create or update
	// is attempted without a valid user remember token hash.
	ErrRememberRequired modelError = "models: Remember hash is required."
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}