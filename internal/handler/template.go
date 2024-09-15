package handler

import (
	"os"
	"path/filepath"
	"text/template"
)

var templates *template.Template

func init() {
	templates = template.New("")
	err := filepath.Walk("template", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err = templates.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
