package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
)

const baseLayoutTemplatePath = "template/base.tmpl"

func (h *handler) registerTemplates() error {
	componentTemplates, err := filepath.Glob("template/components/*.tmpl")
	if err != nil {
		return err
	}

	pageTemplates, err := filepath.Glob("template/pages/*.tmpl")
	if err != nil {
		return err
	}

	for _, path := range pageTemplates {
		tmpl := template.Must(template.ParseFiles(path, baseLayoutTemplatePath))
		tmpl = template.Must(tmpl.ParseFiles(componentTemplates...))
		h.templates[filepath.Base(path)] = tmpl
	}

	return nil
}

func (h *handler) RenderTemplate(w http.ResponseWriter, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.templates[name].ExecuteTemplate(w, "base.tmpl", data)
	return nil
}
