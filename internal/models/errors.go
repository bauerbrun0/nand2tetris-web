package models

import "errors"

var (
	ErrDuplicateEmail    = errors.New("db: duplicate email")
	ErrDuplicateUsername = errors.New("db: duplicate username")
	ErrUserDoesNotExist  = errors.New("db: user does not exist")
	ErrProjectTitleTaken = errors.New("db: project title taken")
	ErrChipNameTaken     = errors.New("db: chip name taken")
)

const (
	ErrorCodeUniqueViolation = "23505"
)

const (
	ConstraintNameUsersUniqueEmail    = "users_unique_constraint_email"
	ConstraintNameUsersUniqueUsername = "users_unique_constraint_username"
)
