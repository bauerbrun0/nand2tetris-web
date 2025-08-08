package routes

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/middleware"
	"github.com/bauerbrun0/nand2tetris-web/ui"
	"github.com/justinas/alice"
)

func GetRoutes(app *application.Application, m *middleware.Middleware, h *handlers.Handlers) http.Handler {
	mux := http.NewServeMux()

	var fs http.Handler

	if app.Config.Env == "production" {
		fs = http.FileServerFS(ui.StaticFiles)
	} else {
		fs = http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static")))
	}

	mux.Handle("GET /static/", m.DisableCacheInDevMode(fs))

	mux.Handle("GET /ping", http.HandlerFunc(h.Ping))

	dynamicChain := alice.New(app.SessionManager.LoadAndSave, m.Language, m.Authenticate, m.GetToasts, m.NoSurf)

	mux.Handle("GET /{$}", dynamicChain.ThenFunc(h.Home))

	requireUnauthenticatedChain := dynamicChain.Append(m.RequireUnathenticated)

	mux.Handle("GET /user/register", requireUnauthenticatedChain.ThenFunc(h.UserRegister))
	mux.Handle("POST /user/register", requireUnauthenticatedChain.ThenFunc(h.UserRegisterPost))
	mux.Handle("GET /user/login", requireUnauthenticatedChain.ThenFunc(h.UserLogin))
	mux.Handle("POST /user/login", requireUnauthenticatedChain.ThenFunc(h.UserLoginPost))
	mux.Handle("GET /user/login/github", requireUnauthenticatedChain.ThenFunc(h.UserLoginGithub))
	mux.Handle("GET /user/login/google", requireUnauthenticatedChain.ThenFunc(h.UserLoginGoogle))
	mux.Handle("GET /user/oauth/github/callback/login", requireUnauthenticatedChain.ThenFunc(h.UserLoginGithubCallback))
	mux.Handle("GET /user/oauth/google/callback/login", requireUnauthenticatedChain.ThenFunc(h.UserLoginGoogleCallback))
	mux.Handle("GET /user/reset-password/send-code", requireUnauthenticatedChain.ThenFunc(h.UserResetPasswordSendCode))
	mux.Handle("POST /user/reset-password/send-code", requireUnauthenticatedChain.ThenFunc(h.UserResetPasswordSendCodePost))
	mux.Handle("GET /user/reset-password/enter-code", requireUnauthenticatedChain.ThenFunc(h.UserResetPasswordEnterCode))
	mux.Handle("POST /user/reset-password/enter-code", requireUnauthenticatedChain.ThenFunc(h.UserResetPasswordEnterCodePost))
	mux.Handle("GET /user/reset-password", requireUnauthenticatedChain.ThenFunc(h.UserResetPassword))
	mux.Handle("POST /user/reset-password", requireUnauthenticatedChain.ThenFunc(h.UserResetPasswordPost))

	requireUnverifiedEmailChain := dynamicChain.Append(m.RequireUnverifiedEmail)

	mux.Handle("GET /user/verify-email", requireUnverifiedEmailChain.ThenFunc(h.UserVerifyEmail))
	mux.Handle("POST /user/verify-email", requireUnverifiedEmailChain.ThenFunc(h.UserVerifyEmailPost))
	mux.Handle("GET /user/verify-email/send-code", requireUnverifiedEmailChain.ThenFunc(h.UserVerifyEmailResendCode))
	mux.Handle("POST /user/verify-email/send-code", requireUnverifiedEmailChain.ThenFunc(h.UserVerifyEmailResendCodePost))

	protectedChain := dynamicChain.Append(m.RequireAuthentication)

	mux.Handle("POST /user/logout", protectedChain.ThenFunc(h.UserLogoutPost))
	mux.Handle("GET /user/settings", protectedChain.ThenFunc(h.UserSettings))
	mux.Handle("POST /user/settings", protectedChain.ThenFunc(h.UserSettingsPost))
	mux.Handle("GET /user/oauth/github/callback/action", protectedChain.ThenFunc(h.UserGithubActionCallback))
	mux.Handle("GET /user/oauth/google/callback/action", protectedChain.ThenFunc(h.UserGoogleActionCallback))
	mux.Handle("GET /user/oauth/google/callback/link", protectedChain.ThenFunc(h.UserLinkGoogleCallback))
	mux.Handle("GET /user/oauth/github/callback/link", protectedChain.ThenFunc(h.UserLinkGithubCallback))

	commonChain := alice.New(m.RecoverPanic, m.LogRequest, m.CommonHeaders)

	return commonChain.Then(mux)
}
