package viewmaster

import (
	"io"
	"text/template/parse"
)

// FuncMap is the type of the map defining the mapping from names to
// functions. Each function must have either a single return value, or two
// return values of which the second has type error. In that case, if the
// second (error) argument evaluates to non-nil during execution, execution
// terminates and Render returns that error. FuncMap has the same base type
// as FuncMap in "text/template", copied here so clients need not import
// "text/template".
type FuncMap map[string]interface{}

type ViewEngine interface {
	// Resolver sets the TemplateResolver to use.
	// Defaults to a FileResolver.
	Resolver(TemplateResolver) ViewEngine

	// Funcs sets the function mapping to use
	// when executing templates.
	Funcs(FuncMap) ViewEngine

	// Render executes the named template and
	// outputs to the writer.
	Render(writer io.Writer, name string, data interface{}) error
}

type baseViewEngine struct {
	resolver   TemplateResolver
	funcs      FuncMap
	leftDelim  string
	rightDelim string
}

func newBaseViewEngine() *baseViewEngine {
	ve := &baseViewEngine{
		resolver: &FileResolver{},
		funcs: map[string]interface{}{
			"sayHi": func(...interface{}) string {
				return "Hello world!"
			},
		},
	}

	return ve
}

func (ve *baseViewEngine) getParseTrees(name string) (string, map[string]*parse.Tree, error) {
	var err error

	var resolvedName = ve.resolver.ResolveName(name, "")

	content, err := ve.resolver.Resolve(resolvedName)
	if err != nil {
		return "", nil, err
	}
	trees, err := parse.Parse(resolvedName, string(content), ve.leftDelim, ve.rightDelim, ve.funcs)
	if err != nil {
		return "", nil, err
	}

	set := map[string]*parse.Tree{}

	for _, tree := range trees {

		err = walk(tree.Root, func(n parse.Node) error {
			return ve.addIncludes(set, resolvedName, n)
		})
		if err != nil {
			return "", nil, err
		}

		set[tree.Name] = tree
	}

	return resolvedName, set, nil
}

func (ve *baseViewEngine) addIncludes(set map[string]*parse.Tree, parentName string, n parse.Node) error {
	switch n := n.(type) {
	case *parse.TemplateNode:

		var resolvedName = ve.resolver.ResolveName(n.Name, parentName)

		found := false

		// check if we've already included this template
		if set[resolvedName] == nil {
			content, err := ve.resolver.Resolve(resolvedName)
			if err != nil {
				return err
			}

			// TODO: consider using a 'TemplateNotFound' error instead of checking for nil
			if content != nil {
				trees, err := parse.Parse(resolvedName, string(content), ve.leftDelim, ve.rightDelim, ve.funcs)
				if err != nil {
					return err
				}

				for _, tree := range trees {
					err = walk(tree.Root, func(n parse.Node) error {
						return ve.addIncludes(set, resolvedName, n)
					})
					if err != nil {
						return err
					}

					set[tree.Name] = tree
				}

				found = true
			}
		} else {
			found = true
		}

		if found {
			// rewrite the template node to use
			// the resolved template name
			n.Name = resolvedName
		}
	}

	return nil
}
