package views

import (
	"html/template"
	"path/filepath"
)

var (
	// LayoutDir gives the layout directory
	LayoutDir = "views/layouts/"
	// TemplateExt gives the template extention
	TemplateExt = ".gohtml"
)

// NewView creates new view templates with default layouts.
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)

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

// layoutFiles returns a slice of strings representing
// the layout files used in our application.
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}
