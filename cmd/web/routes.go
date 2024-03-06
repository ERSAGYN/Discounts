package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/myshops", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showMyShops))
	mux.Get("/shop/:id", dynamicMiddleware.ThenFunc(app.showShop))

	mux.Get("/createProduct", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createProduct))
	mux.Post("/createProduct", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createProduct))

	mux.Get("/createShop", dynamicMiddleware.Append(app.requireAuthentication).Append(app.requireAdmin).ThenFunc(app.createShop))
	mux.Post("/createShop", dynamicMiddleware.Append(app.requireAuthentication).Append(app.requireAdmin).ThenFunc(app.createShop))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
