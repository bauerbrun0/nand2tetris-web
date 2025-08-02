package handlers

import (
	"fmt"
	"reflect"
)

func ParseVerificationMethod(s any) (VerificationMethod, bool) {
	switch v := s.(type) {
	case string:
		return parseVerificationString(v)
	case fmt.Stringer:
		return parseVerificationString(v.String())
	default:
		rv := reflect.ValueOf(s)
		if rv.Kind() == reflect.String {
			return parseVerificationString(rv.String())
		}
		return "", false
	}
}

func parseVerificationString(str string) (VerificationMethod, bool) {
	switch str {
	case string(VerificationPassword), string(VerificationGitHub), string(VerificationGoogle):
		return VerificationMethod(str), true
	default:
		return "", false
	}
}

func ParseAction(s any) (Action, bool) {
	switch v := s.(type) {
	case string:
		return parseActionString(v)
	case fmt.Stringer:
		return parseActionString(v.String())
	default:
		rv := reflect.ValueOf(s)
		if rv.Kind() == reflect.String {
			return parseActionString(rv.String())
		}
		return "", false
	}
}

func parseActionString(str string) (Action, bool) {
	switch str {
	case
		string(ActionChangePassword),
		string(ActionChangeEmail),
		string(ActionChangeEmailSendCode),
		string(ActionCreatePassword),
		string(ActionDeleteAccount),
		string(ActionLinkGitHubAccount),
		string(ActionLinkGoogleAccount),
		string(ActionUnlinkGitHubAccount),
		string(ActionUnlinkGoogleAccount):
		return Action(str), true
	default:
		return "", false
	}
}
