package main

import (
	"net/http"

	"github.com/Pablo-E/snippetbox/ui"
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

	// Take the ui.Files embedded filesystem and convert it to a http.FS type so
	// that it satisfies the http.FileSystem interface. We then pass that to the
	// http.FileServer() function to create the file server handler.
	fileServer := http.FileServer(http.FS(ui.Files))

	// Our static files are contained in the "static" folder of the ui.Files
	// embedded filesystem. So, for example, our CSS stylesheet is located at
	// "static/css/main.css". This means that we now longer need to strip the
	// prefix from the request URL -- any requests that start with /static/ can
	// just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// And then create the routes using the appropriate methods, patterns and
	// handlers.
	router.Handler(http.MethodGet, "/", noSurf(app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.home)))))
	router.Handler(http.MethodGet, "/snippet/view/:id", noSurf(app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.snippetView)))))
	router.Handler(http.MethodGet, "/user/signup", noSurf(app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.userSignup)))))
	router.Handler(http.MethodPost, "/user/signup", noSurf(app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.userSignupPost)))))
	router.Handler(http.MethodGet, "/user/login", noSurf(app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.userLogin)))))
	router.Handler(http.MethodPost, "/user/login", noSurf(app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.userLoginPost)))))

	router.Handler(http.MethodGet, "/snippet/create", noSurf(app.sessionManager.LoadAndSave(app.authenticate(app.requireAuthentication(http.HandlerFunc(app.snippetCreate))))))
	router.Handler(http.MethodPost, "/snippet/create", noSurf(app.sessionManager.LoadAndSave(app.authenticate(app.requireAuthentication(http.HandlerFunc(app.snippetCreatePost))))))
	router.Handler(http.MethodPost, "/user/logout", noSurf(app.sessionManager.LoadAndSave(app.authenticate(app.requireAuthentication(http.HandlerFunc(app.userLogoutPost))))))

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
