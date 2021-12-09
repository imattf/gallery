// Handles all view logic
package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gorilla/csrf"
	"gitlab.com/go-courses/lenslocked.com/context"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

type View struct {
	Template *template.Template
	Layout   string
}

// Render is used to render the View with a predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
		//do nothing
	default:
		vd = Data{
			Yield: data,
		}
	}
	if alert := getAlert(r); alert != nil {
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
		http.Error(w, "Something went wrong. If the problem persists, please email support@faulkners.io", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// ServeHTTP acts as a interface handler
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates and it prepends
// the TemaplteDir directory to each string in the slaice
//
// e.g. the input {"home"} would result in the output
// {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings
// representing file paths for temapltes and it appends
// the TemplateExt extension to each string in the slice
//
// e.g. the input {"home"} would result in the output
// {"home.gohtml"} if TemplateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

// func NewView(layout string, files ...string) *View {
// 	addTemplatePath(files)
// 	addTemplateExt(files)
// 	files = append(files, layoutFiles()...)
// 	t, err := template.New("").Funcs(template.FuncMap{
// 		"csrfField": func() (template.HTML, error) {
// 			return "", errors.New("csrfField is not implemented")
// 		},
// 	}).ParseFiles(files...)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return &View{
// 		Template: t,
// 		Layout:   layout,
// 	}
// }

// Had to use book impl... video impl above is ¡no bueno!
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrfField is not implemented")
		},
		"pathEscape": func(s string) string {
			return url.PathEscape(s)
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

// listFiles returns a slice of strings showing
// the layout files used in app
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	// logging to console
	// fmt.Println(files)
	return files
}
