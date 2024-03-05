package web

import "Discounts/pkg/forms"

type templateData struct {
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	Success         string
}
