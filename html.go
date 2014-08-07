package viewmaster

import (
	"fmt"
	"html/template"
	"io"
)

type htmlViewEngine struct {
	*baseViewEngine
}

// Html is a ViewEngine that uses html/template
// to execute the templates.
func Html() ViewEngine {
	ve := &htmlViewEngine{
		baseViewEngine: newBaseViewEngine(),
	}

	return ve
}

// Resolver sets the TemplateResolver to use.
func (ve *htmlViewEngine) Resolver(r TemplateResolver) ViewEngine {
	ve.resolver = r
	return ve
}

// Funcs sets the FuncMap to use.
func (ve *htmlViewEngine) Funcs(funcs FuncMap) ViewEngine {
	ve.funcs = funcs
	return ve
}

// Render executes the named template and writes
// output to the writer.
func (ve *htmlViewEngine) Render(writer io.Writer, name string, data ...interface{}) error {

	if ve.Resolver == nil {
		return fmt.Errorf("resolver is nil")
	}

	rootName, set, err := ve.getParseTrees(name)
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("").Parse(""))
	tmpl.Funcs(template.FuncMap(ve.funcs))

	for k, v := range set {
		tmpl.AddParseTree(k, v)
	}

	//for _, t := range tmpl.Templates() {
	//	println(t.Name())
	//}

	err = tmpl.ExecuteTemplate(writer, rootName, data)
	if err != nil {
		return err
	}

	return nil
}
