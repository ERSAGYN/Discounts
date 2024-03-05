package web

import (
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	templateCache map[string]*template.Template
}
