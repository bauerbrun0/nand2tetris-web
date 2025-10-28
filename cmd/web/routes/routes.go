package routes

import (
	"net/http"
	"time"

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

	dynamicChain := alice.New(app.SessionManager.LoadAndSave, m.Language, m.Authenticate, m.GetToasts)
	if app.Config.Env == "production" || app.Config.Env == "test" {
		dynamicChain = dynamicChain.Append(m.NoSurf)
	}
	mux.Handle("GET /{$}", dynamicChain.ThenFunc(h.Home))

	requireUnauthenticatedChain := dynamicChain.Append(m.RequireUnathenticated)

	mux.Handle("GET /user/register", requireUnauthenticatedChain.ThenFunc(h.User.UserRegister))
	mux.Handle("POST /user/register", requireUnauthenticatedChain.
		Append(m.GetIPRateLimiter(50, time.Hour)).ThenFunc(h.User.UserRegisterPost))
	mux.Handle("GET /user/login", requireUnauthenticatedChain.ThenFunc(h.User.UserLogin))
	mux.Handle("POST /user/login", requireUnauthenticatedChain.Append(m.UserLoginPostRateLimiter).ThenFunc(h.User.UserLoginPost))
	mux.Handle("GET /user/login/github", requireUnauthenticatedChain.ThenFunc(h.User.UserLoginGithub))
	mux.Handle("GET /user/login/google", requireUnauthenticatedChain.ThenFunc(h.User.UserLoginGoogle))
	mux.Handle("GET /user/oauth/github/callback/login", requireUnauthenticatedChain.ThenFunc(h.User.UserLoginGithubCallback))
	mux.Handle("GET /user/oauth/google/callback/login", requireUnauthenticatedChain.ThenFunc(h.User.UserLoginGoogleCallback))
	mux.Handle("GET /user/reset-password/send-code", requireUnauthenticatedChain.ThenFunc(h.User.UserResetPasswordSendCode))
	mux.Handle("POST /user/reset-password/send-code", requireUnauthenticatedChain.
		Append(m.GetIPRateLimiter(10, time.Minute)).ThenFunc(h.User.UserResetPasswordSendCodePost))
	mux.Handle("GET /user/reset-password/enter-code", requireUnauthenticatedChain.ThenFunc(h.User.UserResetPasswordEnterCode))
	mux.Handle("POST /user/reset-password/enter-code", requireUnauthenticatedChain.
		Append(m.GetIPRateLimiter(10, time.Minute)).ThenFunc(h.User.UserResetPasswordEnterCodePost))
	mux.Handle("GET /user/reset-password", requireUnauthenticatedChain.ThenFunc(h.User.UserResetPassword))
	mux.Handle("POST /user/reset-password", requireUnauthenticatedChain.
		Append(m.GetIPRateLimiter(10, time.Minute)).ThenFunc(h.User.UserResetPasswordPost))

	requireUnverifiedEmailChain := dynamicChain.Append(m.RequireUnverifiedEmail)

	mux.Handle("GET /user/verify-email", requireUnverifiedEmailChain.ThenFunc(h.User.UserVerifyEmail))
	mux.Handle("POST /user/verify-email", requireUnverifiedEmailChain.
		Append(m.GetIPRateLimiter(10, time.Minute)).ThenFunc(h.User.UserVerifyEmailPost))
	mux.Handle("GET /user/verify-email/send-code", requireUnverifiedEmailChain.ThenFunc(h.User.UserVerifyEmailResendCode))
	mux.Handle("POST /user/verify-email/send-code", requireUnverifiedEmailChain.
		Append(m.GetIPRateLimiter(10, time.Minute)).ThenFunc(h.User.UserVerifyEmailResendCodePost))

	protectedChain := dynamicChain.Append(m.RequireAuthentication)

	mux.Handle("POST /user/logout", protectedChain.ThenFunc(h.User.UserLogoutPost))
	mux.Handle("GET /user/settings", protectedChain.ThenFunc(h.User.UserSettings))
	mux.Handle("POST /user/settings", protectedChain.
		Append(m.GetIPRateLimiter(30, 15*time.Minute)).ThenFunc(h.User.UserSettingsPost))
	mux.Handle("GET /user/oauth/github/callback/action", protectedChain.ThenFunc(h.User.UserGithubActionCallback))
	mux.Handle("GET /user/oauth/google/callback/action", protectedChain.ThenFunc(h.User.UserGoogleActionCallback))
	mux.Handle("GET /user/oauth/google/callback/link", protectedChain.ThenFunc(h.User.UserLinkGoogleCallback))
	mux.Handle("GET /user/oauth/github/callback/link", protectedChain.ThenFunc(h.User.UserLinkGithubCallback))

	mux.Handle("GET /hardware-simulator", protectedChain.ThenFunc(h.HardwareSimulator))

	mux.Handle("GET /api/projects", protectedChain.ThenFunc(h.Project.HandleGetProjects))
	mux.Handle("GET /api/projects/{id}", protectedChain.ThenFunc(h.Project.HandleGetProject))
	mux.Handle("GET /api/projects/{slug}/by-slug", protectedChain.ThenFunc(h.Project.HandleGetProjectBySlug))
	mux.Handle("DELETE /api/projects/{id}", protectedChain.ThenFunc(h.Project.HandleDeleteProject))
	mux.Handle("PATCH /api/projects/{id}", protectedChain.ThenFunc(h.Project.HandleUpdateProject))
	mux.Handle("POST /api/projects", protectedChain.ThenFunc(h.Project.HandleCreateProject))

	mux.Handle("POST /api/projects/{projectId}/chips", protectedChain.ThenFunc(h.Chip.HandleCreateChip))
	mux.Handle("GET /api/projects/{projectId}/chips", protectedChain.ThenFunc(h.Chip.HandleGetChips))
	mux.Handle("DELETE /api/projects/{projectId}/chips/{chipId}", protectedChain.ThenFunc(h.Chip.HandleDeleteChip))
	mux.Handle("PATCH  /api/projects/{projectId}/chips/{chipId}", protectedChain.ThenFunc(h.Chip.HandleUpdateChip))

	mux.Handle("GET /projects", protectedChain.ThenFunc(h.Projects))

	commonChain := alice.New(m.RecoverPanic, m.LogRequest, m.CommonHeaders)

	return commonChain.Then(mux)
}
