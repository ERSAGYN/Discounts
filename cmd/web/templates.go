package main

import (
	"Discounts/pkg/forms"
	"Discounts/pkg/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	IsAdmin         bool
	Success         string
	CurrentYear     int
	User            *models.User
	Users           []*models.User
	Product         *models.Product
	Products        []*models.Product
	Shop            *models.Shop
	Shops           []*models.Shop
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
func calculateDiscountedPrice(price, discount int) int {
	return price - (price * discount / 100)
}

var functions = template.FuncMap{
	"calculateDiscountedPrice": calculateDiscountedPrice,
}
