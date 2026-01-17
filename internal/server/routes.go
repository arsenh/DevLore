package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arsenh/DevLore/internal/auth"
	"github.com/arsenh/DevLore/internal/config"
	"github.com/arsenh/DevLore/internal/logger"
	customMiddlewares "github.com/arsenh/DevLore/internal/middleware"
	"github.com/arsenh/DevLore/internal/model"
	"github.com/arsenh/DevLore/internal/service"
	"github.com/arsenh/DevLore/internal/templates"
	"github.com/arsenh/DevLore/internal/views"

	"github.com/asaskevich/govalidator"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

//go:embed public
var publicFS embed.FS

type Routes struct {
	articleService *service.ArticleService
	userService    *service.UserService
	rateLimiter    *RateLimiter
}

func NewRoutes(articleSvc *service.ArticleService, userSvc *service.UserService, limiter *RateLimiter) *Routes {
	return &Routes{
		articleService: articleSvc,
		userService:    userSvc,
		rateLimiter:    limiter,
	}
}

func (r *Routes) GetRoutes() http.Handler {
	router := chi.NewRouter()

	// setup logrus for requests
	chiLogger := middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: logger.L(),
	})

	router.Use(middleware.RealIP)
	router.Use(chiLogger)

	// setup static files
	publicFiles, err := fs.Sub(publicFS, config.PublicDir)
	if err != nil {
		panic(fmt.Errorf("parse static files: %w", err))
	}

	router.Handle(fmt.Sprintf("/%s/*", config.StaticDir), http.StripPrefix(fmt.Sprintf("/%s/", config.StaticDir), http.FileServerFS(publicFiles)))

	// global 404 not found page
	router.NotFound(r.notFoundPage)

	router.Get("/", r.rootHandler)
	router.Get("/articles/new", r.showNewArticleHandler)
	router.Post("/articles/new", r.createNewArticleHandler)
	router.Post("/articles/{id}/edit", r.editArticleHandler)
	router.Get("/search", r.showSearchHandler)
	router.Get("/email-exists", r.emailExistHandler)

	router.Group(func(router chi.Router) {
		router.Use(customMiddlewares.JWTAuthMiddleware(string(config.JWTSecretKey)))
		router.Get("/auth/register", r.showRegisterHandler)
		router.Post("/auth/register", r.registerUserHandler)
		router.Get("/auth/login", r.showLoginHandler)
		router.Post("/auth/login", r.loginUserHandler)
		router.Get("/dashboard", r.dashboardHandler)
		router.Post("/auth/logout", r.logoutHandler)

		router.Get("/articles/{id}", r.viewArticleHandler)
		router.Get("/articles/{id}/edit", r.viewEditArticleHandler)
		router.Post("/articles/{id}/delete", r.deleteArticleHandler)
	})

	return router
}

func (r *Routes) retrieveId(request *http.Request) (int, error) {
	idStr := chi.URLParam(request, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, logger.LogAndErr("id cannot be parsed into integer")
	}
	return id, nil
}

func (r *Routes) setJWTTokenAsCookie(writer http.ResponseWriter, token string) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false, // MUST be false on HTTP
	})
}

func (r *Routes) deleteJWTTokenFromCookie(writer http.ResponseWriter) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Jan 1, 1970
		MaxAge:   -1,              // also forces deletion
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // true if HTTPS
	})
}

func (r *Routes) isUserAuthenticated(ctx context.Context) bool {
	authErr := ctx.Value(customMiddlewares.CtxAuthError)
	return authErr == nil
}

func (r *Routes) getIP(request *http.Request) string {
	if request.RemoteAddr != "" {
		return strings.Split(request.RemoteAddr, ":")[0]
	}
	return "unknown" // to not broke rate limiter, if address not found, fallback will be used
}

func (r *Routes) getAuthenticatedUserIfAny(ctx context.Context, writer http.ResponseWriter) *model.User {
	if !r.isUserAuthenticated(ctx) {
		// delete JWT token
		r.deleteJWTTokenFromCookie(writer)
		return nil
	}

	userID := ctx.Value(customMiddlewares.CtxUserID).(int)
	//TODO: handle err from service
	user, _ := r.userService.GetUserByID(ctx, userID)
	if user == nil {
		r.deleteJWTTokenFromCookie(writer)
		return nil
	} else {
		return user
	}
}

func (r *Routes) getUserNameIfNotNil(user *model.User) string {
	userName := ""
	if user != nil {
		userName = user.FullName
	}
	return userName
}

func (r *Routes) notPermitted(fullName string, articleID int, writer http.ResponseWriter) {
	notPermittedView := struct {
		UserName string
		ID       int
	}{
		UserName: fullName,
		ID:       articleID,
	}
	if err := templates.Render(writer, http.StatusOK, templates.NotPermitted, notPermittedView); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) notFoundPage(writer http.ResponseWriter, request *http.Request) {
	templates.NotFound(writer)
}

func (r *Routes) dashboardHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	user := r.getAuthenticatedUserIfAny(ctx, writer)

	//TODO: add limit on articles count
	articleItems, err := r.articleService.GetDashboardData(request.Context())
	if err != nil {
		templates.InternalServerError(writer, err)
		return
	}

	view := &views.DashboardView{
		BaseView: views.BaseView{
			UserName: r.getUserNameIfNotNil(user),
		},
		Articles: articleItems,
	}

	if err := templates.Render(writer, http.StatusOK, templates.DashboardTemplate, view); err != nil {
		templates.InternalServerError(writer, err)
		return
	}
}

func (r *Routes) rootHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/dashboard", http.StatusPermanentRedirect)
}

func (r *Routes) showNewArticleHandler(writer http.ResponseWriter, request *http.Request) {
	// render form for new article creation
	if err := templates.Render(writer, http.StatusOK, templates.NewArticleTemplate, nil); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) createNewArticleHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	if err := request.ParseForm(); err != nil {
		templates.BadRequest(writer)
		return
	}
	//TODO: need to get also UserId which created the article.
	title := request.Form.Get("title")
	content := request.Form.Get("content")
	id, err := r.articleService.SaveArticle(ctx, title, content)
	if err != nil {
		templates.InternalServerError(writer, err)
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/articles/%d", id), http.StatusSeeOther)
}

func (r *Routes) viewArticleHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	user := r.getAuthenticatedUserIfAny(ctx, writer)

	id, err := r.retrieveId(request)
	if err != nil {
		templates.BadRequest(writer)
		return
	}

	view, err := r.articleService.GetArticleById(ctx, id)
	if err != nil {
		templates.NotFound(writer)
		return
	}

	view.UserName = r.getUserNameIfNotNil(user)

	authorUser, err := r.userService.GetUserByID(ctx, view.Article.UserId)
	if err != nil {
		templates.InternalServerError(writer, err)
		return
	}

	// set user name of article author
	view.CreatedBy = authorUser.FullName

	if err := templates.Render(writer, http.StatusOK, templates.ViewArticleTemplate, view); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) deleteArticleHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := r.retrieveId(request)
	if err != nil {
		templates.BadRequest(writer)
	}

	ctx := request.Context()
	user := r.getAuthenticatedUserIfAny(ctx, writer)

	if user == nil {
		http.Redirect(writer, request, "/auth/login", http.StatusSeeOther)
		return
	}

	view, err := r.articleService.GetArticleById(ctx, id)
	if err != nil {
		templates.NotFound(writer)
		return
	}

	if view.Article.UserId != user.ID {
		r.notPermitted(user.FullName, view.Article.ID, writer)
		return
	}

	if err = r.articleService.DeleteArticleById(ctx, id); err != nil {
		templates.InternalServerError(writer, err)
	}

	http.Redirect(writer, request, "/dashboard", http.StatusSeeOther)
}

func (r *Routes) viewEditArticleHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := r.retrieveId(request)
	if err != nil {
		templates.BadRequest(writer)
		return
	}

	ctx := request.Context()
	user := r.getAuthenticatedUserIfAny(ctx, writer)

	if user == nil {
		http.Redirect(writer, request, "/auth/login", http.StatusSeeOther)
		return
	}

	view, err := r.articleService.GetArticleById(request.Context(), id)
	if err != nil {
		templates.NotFound(writer)
		return
	}

	if view.Article.UserId != user.ID {
		r.notPermitted(user.FullName, view.Article.ID, writer)
		return
	}

	if err := templates.Render(writer, http.StatusOK, templates.EditArticleTemplate, view); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) editArticleHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := r.retrieveId(request)
	if err != nil {
		templates.BadRequest(writer)
		return
	}

	_, err = r.articleService.GetArticleById(request.Context(), id)
	if err != nil {
		templates.BadRequest(writer)
		return
	}

	if err := request.ParseForm(); err != nil {
		templates.BadRequest(writer)
		return
	}

	title := request.Form.Get("title")
	content := request.Form.Get("content")

	view, err := r.articleService.EditArticleById(request.Context(), id, title, content)
	if err != nil {
		templates.InternalServerError(writer, err)
		return
	}

	if err := templates.Render(writer, http.StatusOK, templates.ViewArticleTemplate, view); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) showSearchHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	query := strings.TrimSpace(request.URL.Query().Get("q"))

	view, err := r.articleService.SearchArticlesByQuery(ctx, query)
	if err != nil {
		templates.InternalServerError(writer, err)
		return
	}
	if err := templates.Render(writer, http.StatusOK, templates.SearchTemplate, view); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) showRegisterHandler(writer http.ResponseWriter, request *http.Request) {

	// redirect to dashboard if user already authenticated
	if r.isUserAuthenticated(request.Context()) {
		http.Redirect(writer, request, "/dashboard", http.StatusSeeOther)
	}

	if err := templates.Render(writer, http.StatusOK, templates.RegisterTemplate, nil); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) registerUserHandler(writer http.ResponseWriter, request *http.Request) {
	// redirect to dashboard if user already authenticated
	if r.isUserAuthenticated(request.Context()) {
		http.Redirect(writer, request, "/dashboard", http.StatusSeeOther)
	}

	// need to get user email, password, confirm password, full name
	if err := request.ParseForm(); err != nil {
		templates.BadRequest(writer)
		return
	}

	email := request.Form.Get("email")
	fullName := request.Form.Get("full_name")
	password := request.Form.Get("password")
	passwordConfirm := request.Form.Get("password_confirm")

	if !govalidator.IsEmail(email) ||
		(fullName == "") ||
		(password == "") ||
		(passwordConfirm == "") ||
		(password != passwordConfirm) {
		templates.BadRequest(writer)
		return
	}

	ctx := request.Context()

	user := r.userService.GetUserByEmail(ctx, email)
	if user != nil {
		// The presence of an email address is checked using JavaScript
		// If the user bypasses JavaScript, an invalid request will be displayed.
		templates.BadRequest(writer)
		return
	}

	registeredUser, err := r.userService.RegisterUser(ctx, email, password, fullName)
	if err != nil {
		templates.InternalServerError(writer, err)
		return
	}

	token, err := auth.GenerateJWTToken(registeredUser.ID, registeredUser.Email, registeredUser.FullName)
	if err != nil {
		templates.InternalServerError(writer, err)
		return
	}

	r.setJWTTokenAsCookie(writer, token)
	http.Redirect(writer, request, "/dashboard", http.StatusSeeOther)
}

func (r *Routes) loginUserHandler(writer http.ResponseWriter, request *http.Request) {
	// redirect to dashboard if user already authenticated
	if r.isUserAuthenticated(request.Context()) {
		http.Redirect(writer, request, "/dashboard", http.StatusSeeOther)
	}

	// need to get user email, password
	if err := request.ParseForm(); err != nil {
		templates.BadRequest(writer)
		return
	}

	email := request.Form.Get("email")
	password := request.Form.Get("password")

	if !govalidator.IsEmail(email) || (password == "") {
		templates.BadRequest(writer)
		return
	}

	ctx := request.Context()

	user := r.userService.GetUserByEmail(ctx, email)

	if user == nil {
		view := struct {
			UserName string
			Email    string
			Error    string
		}{
			UserName: "",
			Email:    email,
			Error:    "Authentication failed. Incorrect login or password.",
		}
		templates.Render(writer, http.StatusUnauthorized, templates.LoginTemplate, view)
		return
	}

	// check user password
	ok := r.userService.CheckUserPassword(ctx, user.ID, password)
	if !ok {
		templates.BadRequest(writer)
		return
	}

	// user exists, password correct, but token is expired or not exist
	// create new token
	token, err := auth.GenerateJWTToken(user.ID, user.Email, user.FullName)
	if err != nil {
		templates.InternalServerError(writer, err)
		return
	}

	r.setJWTTokenAsCookie(writer, token)
	http.Redirect(writer, request, "/dashboard", http.StatusSeeOther)
}

func (r *Routes) showLoginHandler(writer http.ResponseWriter, request *http.Request) {
	// redirect to dashboard if user already authenticated
	if r.isUserAuthenticated(request.Context()) {
		http.Redirect(writer, request, "/dashboard", http.StatusSeeOther)
	}

	if err := templates.Render(writer, http.StatusOK, templates.LoginTemplate, nil); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) emailExistHandler(writer http.ResponseWriter, request *http.Request) {
	//rate limiter to give 5 requests per minute
	limiter := r.rateLimiter.GetLimiterIp(r.getIP(request))

	if !limiter.Allow() {
		http.Error(writer, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	raw := request.URL.Query().Get("email")

	email := strings.TrimSpace(strings.ToLower(raw))

	if len(email) > 255 {
		email = ""
	}

	logger.L().Infoln("email: ", email)
	logger.L().Infoln("email len: ", len(email))

	exists := false
	if user := r.userService.GetUserByEmail(request.Context(), email); user != nil {
		exists = true
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]bool{
		"exists": exists,
	})
}

func (r *Routes) logoutHandler(writer http.ResponseWriter, request *http.Request) {
	if r.isUserAuthenticated(request.Context()) {
		r.deleteJWTTokenFromCookie(writer)
	}
	http.Redirect(writer, request, "/dashboard", http.StatusSeeOther)
}
