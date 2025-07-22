package main

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const authenticatedUserInfoKey = contextKey("authenticatedUserInfo")
const initialToastsKey = contextKey("initialToasts")
