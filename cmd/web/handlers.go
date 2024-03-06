package main

import (
	"Discounts/pkg/forms"
	"Discounts/pkg/models"
	"errors"
	"net/http"
)

func (app *application) showMyShops(w http.ResponseWriter, r *http.Request) {
	// Retrieve the shops associated with the currently logged-in user
	//shops, err := app.shops.GetByOwner(app.session.GetInt(r, "authenticatedUserID"))
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}

	// Render a template displaying those shops
	//app.render(w, r, "myshops.page.tmpl", &templateData{
	//	Shops: shops,
	//})
}

// Implement your show shop handler
func (app *application) showShop(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	shop, err := app.shops.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "shop.page.tmpl", &templateData{Shop: shop})
}

// Implement your create product form handler
func (app *application) createProductForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create_product.page.tmpl", nil)
}
func (app *application) createProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("productName", "category", "price", "discount")
	form.MinLength("productName", 3)
	form.MaxLength("productName", 255)
	form.MinLength("category", 3)
	form.MaxLength("category", 255)

	if !form.Valid() {
		app.render(w, r, "create_product_form.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	// Process form data and save the product to the database
	// You'll need to access form values like r.PostForm.Get("fieldname")

	// Redirect the user after saving the product
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Implement your create shop form handler
func (app *application) createShopForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create_shop.page.tmpl", nil)
}

// Implement your create shop handler
func (app *application) createShop(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("shopName", "address")
	form.MinLength("shopName", 3)
	form.MaxLength("shopName", 255)
	form.MinLength("address", 3)
	form.MaxLength("address", 255)

	if !form.Valid() {
		app.render(w, r, "create_shop_form.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	// Process form data and save the shop to the database
	// You'll need to access form values like r.PostForm.Get("fieldname")

	// Redirect the user after saving the shop
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.shops.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Shops: s,
	})
}
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		}
		return
	}
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.MatchesPattern("email", forms.EmailRX)
	user, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
