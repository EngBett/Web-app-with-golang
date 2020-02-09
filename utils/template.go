package utils

import "html/template"

var templates *template.Template

func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

templates.ExecuteTemplate(w, "login.html", nil)