package handlers

import "errors"

// VerificationMethod represents the method a user uses to re-authenticate
// before performing a sensitive action (e.g. deleting their account).
type VerificationMethod string

const (
	VerificationPassword VerificationMethod = "password"
	VerificationGitHub   VerificationMethod = "github"
	VerificationGoogle   VerificationMethod = "google"
)

// Action represents a user-initiated operation, typically triggered from a multi-form page
// like user settings. It is used to identify which specific action the user intended to perform,
// such as changing a password, linking an OAuth account, or deleting their account.
type Action string

const (
	ActionChangePassword      Action = "change-password"
	ActionChangeEmail         Action = "change-email"
	ActionChangeEmailSendCode Action = "change-email-send-code"
	ActionCreatePassword      Action = "create-password"
	ActionDeleteAccount       Action = "delete-account"
	ActionLinkGitHubAccount   Action = "link-github-account"
	ActionLinkGoogleAccount   Action = "link-google-account"
	ActionUnlinkGitHubAccount Action = "unlink-github-account"
	ActionUnlinkGoogleAccount Action = "unlink-google-account"
)

var ErrInvalidActionInSession = errors.New("handlers: invalid action stored in session")
