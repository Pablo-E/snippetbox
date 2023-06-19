package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// Create a handler function which wraps our notFound() helper, and then
	// assign it as the custom handler for 404 Not Found responses. You can also
	// set a custom handler for 405 Method Not Allowed responses by setting
	// router.MethodNotAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Update the pattern for the route for the static files.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// And then create the routes using the appropriate methods, patterns and
	// handlers.
	router.Handler(http.MethodGet, "/", noSurf(app.sessionManager.LoadAndSave(http.HandlerFunc(app.home))))
	router.Handler(http.MethodGet, "/snippet/view/:id", noSurf(app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetView))))
	router.Handler(http.MethodGet, "/user/signup", noSurf(app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignup))))
	router.Handler(http.MethodPost, "/user/signup", noSurf(app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignupPost))))
	router.Handler(http.MethodGet, "/user/login", noSurf(app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogin))))
	router.Handler(http.MethodPost, "/user/login", noSurf(app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLoginPost))))

	router.Handler(http.MethodGet, "/snippet/create", noSurf(app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.snippetCreate)))))
	router.Handler(http.MethodPost, "/snippet/create", noSurf(app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.snippetCreatePost)))))
	router.Handler(http.MethodPost, "/user/logout", noSurf(app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.userLogoutPost)))))

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
