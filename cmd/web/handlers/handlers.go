package handlers

import (
	"errors"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
)

type Handlers struct {
	*application.Application
}

func NewHandlers(app *application.Application) *Handlers {
	return &Handlers{
		Application: app,
	}
}

// VerificationMethod represents the method a user uses to re-authenticate
// before performing a sensitive action (e.g. deleting their account).
type VerificationMethod string

const (
	VerificationPassword VerificationMethod = "password"
	VerificationGitHub   VerificationMethod = "github"
	VerificationGoogle   VerificationMethod = "google"
)

func ParseVerificationMethod(s string) (VerificationMethod, bool) {
	switch s {
	case string(VerificationPassword), string(VerificationGitHub), string(VerificationGoogle):
		return VerificationMethod(s), true
	default:
		return "", false
	}
}

// AuthenticatedAction represents user-initiated actions that require re-authentication.
// These actions are typically sensitive operations like account management or OAuth changes.
type AuthenticatedAction string

const (
	ActionDeleteAccount       AuthenticatedAction = "delete-account"
	ActionLinkGitHubAccount   AuthenticatedAction = "link-github-account"
	ActionLinkGoogleAccount   AuthenticatedAction = "link-google-account"
	ActionUnlinkGitHubAccount AuthenticatedAction = "unlink-github-account"
	ActionUnlinkGoogleAccount AuthenticatedAction = "unlink-google-account"
)

var ErrInvalidActionInSession = errors.New("handlers: invalid authenticated action stored in session")
