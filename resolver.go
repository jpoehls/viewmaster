package viewmaster

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// TemplateResolver implementations are used to locate
// template references that aren't already in the template set.
type TemplateResolver interface {

	// Resolve returns the content of the named template.
	Resolve(name string) ([]byte, error)

	// ResolveName returns the expanded template name that should
	// be used when resolving the template. e.g. FileResolver
	// uses this to expand the name to the absolute file path.
	ResolveName(name string, parentName string) string
}

// FileResolver is an implementation of TemplateResolver
// that resolves template names on the file system.
type FileResolver struct {

	// Dir is the root directory
	// to resolve relative template names from.
	Dir string
}

// Resolve tries to return the content of the file witht
// the specified name. If the file doesn't exist then
// nothing is returned (nil, nil).
func (r *FileResolver) Resolve(name string) ([]byte, error) {
	content, err := ioutil.ReadFile(name)
	if os.IsNotExist(err) {
		// file not found
		// don't bubble the error
		return nil, nil
	}

	return content, err
}

// ResolveName expands the specified name into an absolute file path.
func (r *FileResolver) ResolveName(name string, parentName string) string {
	file := name
	if !filepath.IsAbs(file) {
		// assume file name is relative to the parentName
		parentDir := filepath.Dir(parentName)
		if parentDir == "" {
			parentDir = r.Dir
		}
		file = filepath.Join(parentDir, file)

		file, err := filepath.Abs(file)
		if err != nil {
			// error resolving file path
			// return name unchanged
			return name
		}
		file = filepath.Clean(file)
	}
	return file
}
