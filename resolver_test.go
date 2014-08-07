package viewmaster

type testResolver struct {
	templates map[string][]byte
}

func newTestResolver() *testResolver {
	return &testResolver{
		templates: map[string][]byte{},
	}
}

func (r *testResolver) Resolve(name string) ([]byte, error) {
	return r.templates[name], nil
}

func (r *testResolver) ResolveName(name string, parentName string) string {
	return name
}

func (r *testResolver) Add(name string, template string) {
	r.templates[name] = []byte(template)
}
