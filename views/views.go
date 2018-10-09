package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
	"github.com/jackytck/lenslocked/context"
)

var (
	// LayoutDir gives the layout directory
	LayoutDir = "views/layouts/"
	// TemplateDir is the prefix of the views directory
	TemplateDir = "views/"
	// TemplateExt gives the template extention
	TemplateExt = ".gohtml"
)

// NewView creates new view templates with default layouts.
func NewView(layout string, files ...string) *View {
	files = append(addTemplateExt(addTemplatePath(files)), layoutFiles()...)

	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrfField is not implemented")
		},
	}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View represents a html template.
type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Render is used to render the view with the predefined layout.
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		// do nothing
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	if alert := getAlert(r); alert != nil && vd.Alert == nil {
		vd.Alert = alert
		clearAlert(w)
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// layoutFiles returns a slice of strings representing
// the layout files used in our application.
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func addTemplatePath(files []string) []string {
	ret := make([]string, len(files))
	for i, f := range files {
		ret[i] = TemplateDir + f
	}
	return ret
}

func addTemplateExt(files []string) []string {
	ret := make([]string, len(files))
	for i, f := range files {
		ret[i] = f + TemplateExt
	}
	return ret
}
