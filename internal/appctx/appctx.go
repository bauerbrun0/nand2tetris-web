package appctx

type contextKey string

const IsAuthenticatedContextKey = contextKey("isAuthenticated")
const AuthenticatedUserInfoKey = contextKey("authenticatedUserInfo")
const InitialToastsKey = contextKey("initialToasts")
const LocalizerKey = contextKey("localizer")
