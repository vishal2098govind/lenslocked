package views

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTpl *template.Template
}

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t

}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	buf := &strings.Builder{}

	tpl, _ := t.htmlTpl.Clone()
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
	})

	err := tpl.Execute(buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing template", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, buf.String())
}

func ParseFS(fs fs.FS, patterns ...string) (*Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrfField function not implemented")
		},
	})
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	return &Template{htmlTpl: tpl}, nil
}

func Parse(filepath string) (*Template, error) {
	t, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("failed to parse template file: %v", err)
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	errors.Is(err, errors.New(""))
	return &Template{htmlTpl: t}, nil
}
