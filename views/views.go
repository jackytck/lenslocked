package views

import "html/template"

// NewView creates new view templates with default layouts.
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
	}
}

// View represents a html template.
type View struct {
	Template *template.Template
}
