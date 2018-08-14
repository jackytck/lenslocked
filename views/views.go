package views

import (
	"html/template"
	"net/http"
	"path/filepath"
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

	t, err := template.ParseFiles(files...)
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
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render is used to render the view with the predefined layout.
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type) {
	case Data:
		// do nothing
	default:
		data = Data{
			Yield: data,
		}
	}
	return v.Template.ExecuteTemplate(w, v.Layout, data)
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
