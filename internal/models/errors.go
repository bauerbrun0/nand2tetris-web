package models

import "errors"

var (
	ErrDuplicateEmail    = errors.New("db: duplicate email")
	ErrDuplicateUsername = errors.New("db: duplicate username")
	ErrUserDoesNotExist  = errors.New("db: user does not exist")
)

const (
	ErrorCodeUniqueViolation = "23505"
)

const (
	ConstraintNameUsersUniqueEmail    = "users_unique_constraint_email"
	ConstraintNameUsersUniqueUsername = "users_unique_constraint_username"
)
