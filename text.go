package viewmaster

import (
	"fmt"
	"io"
	"text/template"
)

type textViewEngine struct {
	*baseViewEngine
}

func Text() ViewEngine {
	ve := &textViewEngine{
		baseViewEngine: newBaseViewEngine(),
	}

	return ve
}

func (ve *textViewEngine) Resolver(r TemplateResolver) ViewEngine {
	ve.resolver = r
	return ve
}

func (ve *textViewEngine) Funcs(funcs FuncMap) ViewEngine {
	ve.funcs = funcs
	return ve
}

// Render executes the named template and writes
// output to the writer.
func (ve *textViewEngine) Render(writer io.Writer, name string, data ...interface{}) error {

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
