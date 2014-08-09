package viewmaster

import (
	"bytes"
	"testing"
)

func assertHtmlOutput(t *testing.T, r TemplateResolver, name, expected string, data interface{}) {
	v := Html().Resolver(r)

	b := &bytes.Buffer{}
	err := v.Render(b, name, data)
	actual := b.String()
	t.Logf("actual: %q", actual)
	if err != nil {
		t.Fatal(err.Error())
	}

	if b.String() != expected {
		t.Fatalf("expected: %q", expected)
	}
}

func TestHtmlSinglePage(t *testing.T) {
	r := newTestResolver()
	r.Add("page.html", "PAGE")

	assertHtmlOutput(t, r, "page.html", "PAGE", nil)
}

func TestHtmlSimpleLayout(t *testing.T) {
	r := newTestResolver()
	r.Add("layout.html", "HEADER-{{template \"body\"}}-FOOTER")
	r.Add("page.html", "{{define \"body\"}}BODY{{end}}{{template \"layout.html\"}}")

	assertHtmlOutput(t, r, "page.html", "HEADER-BODY-FOOTER", nil)
}

func TestHtmlLayoutWithIncludes(t *testing.T) {
	r := newTestResolver()
	r.Add("layout.html", "HEADER-{{template \"body\"}}-FOOTER-{{template \"include.html\"}}")
	r.Add("include.html", "INCLUDE")
	r.Add("page.html", "{{define \"body\"}}BODY-{{template \"include.html\"}}{{end}}{{template \"layout.html\"}}")

	assertHtmlOutput(t, r, "page.html", "HEADER-BODY-INCLUDE-FOOTER-INCLUDE", nil)
}

func TestHtmlNestedLayouts(t *testing.T) {
	r := newTestResolver()
	r.Add("top_layout.html", "HEADER-{{template \"body\"}}-FOOTER")
	r.Add("sub_layout.html", "{{define \"body\"}}WRAP-{{template \"sub_body\"}}-WRAP{{end}}{{template \"top_layout.html\"}}")
	r.Add("page.html", "{{define \"sub_body\"}}BODY{{end}}{{template \"sub_layout.html\"}}")

	assertHtmlOutput(t, r, "page.html", "HEADER-WRAP-BODY-WRAP-FOOTER", nil)
}

func TestHtmlData(t *testing.T) {
	r := newTestResolver()
	r.Add("page.html", "Name: {{.Name}}")

	d := struct {
		Name string
	}{
		Name: "Joshua",
	}

	assertHtmlOutput(t, r, "page.html", "Name: Joshua", d)
}
