package main

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	var fs http.Handler

	if app.config.env == "production" {
		fs = http.FileServerFS(ui.StaticFiles)
	} else {
		fs = http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static")))
	}

	mux.Handle("GET /static/", app.disableCacheInDevMode(fs))

	dynamicChain := alice.New(app.sessionManager.LoadAndSave, app.authenticate, app.getToasts)

	mux.Handle("GET /{$}", dynamicChain.ThenFunc(app.home))

	mux.Handle("GET /user/reset-password/send-code", dynamicChain.ThenFunc(app.userResetPasswordSendCode))
	mux.Handle("POST /user/reset-password/send-code", dynamicChain.ThenFunc(app.userResetPasswordSendCodePost))
	mux.Handle("GET /user/reset-password/enter-code", dynamicChain.ThenFunc(app.userResetPasswordEnterCode))
	mux.Handle("POST /user/reset-password/enter-code", dynamicChain.ThenFunc(app.userResetPasswordEnterCodePost))
	mux.Handle("GET /user/reset-password", dynamicChain.ThenFunc(app.userResetPassword))
	mux.Handle("POST /user/reset-password", dynamicChain.ThenFunc(app.userResetPasswordPost))

	requireUnauthenticatedChain := dynamicChain.Append(app.requireUnathenticated)

	mux.Handle("GET /user/register", requireUnauthenticatedChain.ThenFunc(app.userRegister))
	mux.Handle("POST /user/register", requireUnauthenticatedChain.ThenFunc(app.userRegisterPost))
	mux.Handle("GET /user/login", requireUnauthenticatedChain.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", requireUnauthenticatedChain.ThenFunc(app.userLoginPost))

	requireUnverifiedEmailChain := dynamicChain.Append(app.requireUnverifiedEmail)

	mux.Handle("GET /user/verify-email", requireUnverifiedEmailChain.ThenFunc(app.userVerifyEmail))
	mux.Handle("POST /user/verify-email", requireUnverifiedEmailChain.ThenFunc(app.userVerifyEmailPost))
	mux.Handle("GET /user/verify-email/send-code", requireUnverifiedEmailChain.ThenFunc(app.userVerifyEmailResendCode))
	mux.Handle("POST /user/verify-email/send-code", requireUnverifiedEmailChain.ThenFunc(app.userVerifyEmailResendCodePost))

	protectedChain := dynamicChain.Append(app.requireAuthentication)

	mux.Handle("POST /user/logout", protectedChain.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /user/settings", protectedChain.ThenFunc(app.userSettings))
	mux.Handle("POST /user/settings", protectedChain.ThenFunc(app.userSettingsPost))

	commonChain := alice.New(app.recoverPanic, app.logRequest, app.commonHeaders)

	return commonChain.Then(mux)
}
