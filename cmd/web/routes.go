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

	dynamicChain := alice.New(app.sessionManager.LoadAndSave, app.authenticate)

	mux.Handle("GET /{$}", dynamicChain.ThenFunc(app.home))
	mux.Handle("GET /user/register", dynamicChain.ThenFunc(app.userRegister))
	mux.Handle("POST /user/register", dynamicChain.ThenFunc(app.userRegisterPost))
	mux.Handle("GET /user/verify-email", dynamicChain.ThenFunc(app.userVerifyEmail))
	mux.Handle("POST /user/verify-email", dynamicChain.ThenFunc(app.userVerifyEmailPost))
	mux.Handle("GET /user/verify-email/send-code", dynamicChain.ThenFunc(app.userVerifyEmailResendCode))
	mux.Handle("POST /user/verify-email/send-code", dynamicChain.ThenFunc(app.userVerifyEmailResendCodePost))
	mux.Handle("GET /user/login", dynamicChain.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamicChain.ThenFunc(app.userLoginPost))
	mux.Handle("GET /user/reset-password/send-code", dynamicChain.ThenFunc(app.userResetPasswordSendCode))
	mux.Handle("POST /user/reset-password/send-code", dynamicChain.ThenFunc(app.userResetPasswordSendCodePost))
	mux.Handle("GET /user/reset-password/enter-code", dynamicChain.ThenFunc(app.userResetPasswordEnterCode))
	mux.Handle("POST /user/reset-password/enter-code", dynamicChain.ThenFunc(app.userResetPasswordEnterCodePost))
	mux.Handle("GET /user/reset-password", dynamicChain.ThenFunc(app.userResetPassword))
	mux.Handle("POST /user/reset-password", dynamicChain.ThenFunc(app.userResetPasswordPost))

	protectedChain := dynamicChain.Append(app.requireAuthentication)

	mux.Handle("POST /user/logout", protectedChain.ThenFunc(app.userLogoutPost))

	commonChain := alice.New(app.recoverPanic, app.logRequest, app.commonHeaders)

	return commonChain.Then(mux)
}
