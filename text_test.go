package viewmaster

import (
	"bytes"
	"testing"
)

func assertTextOutput(t *testing.T, r TemplateResolver, name, expected string) {
	v := Text().Resolver(r)

	b := &bytes.Buffer{}
	err := v.Render(b, name, nil)
	actual := b.String()
	t.Logf("actual: %q", actual)
	if err != nil {
		t.Fatal(err.Error())
	}

	if b.String() != expected {
		t.Fatalf("expected: %q", expected)
	}
}

func TestTextSinglePage(t *testing.T) {
	r := newTestResolver()
	r.Add("page.html", "PAGE")

	assertTextOutput(t, r, "page.html", "PAGE")
}

func TestTextSimpleLayout(t *testing.T) {
	r := newTestResolver()
	r.Add("layout.html", "HEADER-{{template \"body\"}}-FOOTER")
	r.Add("page.html", "{{define \"body\"}}BODY{{end}}{{template \"layout.html\"}}")

	assertTextOutput(t, r, "page.html", "HEADER-BODY-FOOTER")
}

func TestTextLayoutWithIncludes(t *testing.T) {
	r := newTestResolver()
	r.Add("layout.html", "HEADER-{{template \"body\"}}-FOOTER-{{template \"include.html\"}}")
	r.Add("include.html", "INCLUDE")
	r.Add("page.html", "{{define \"body\"}}BODY-{{template \"include.html\"}}{{end}}{{template \"layout.html\"}}")

	assertTextOutput(t, r, "page.html", "HEADER-BODY-INCLUDE-FOOTER-INCLUDE")
}

func TestTextNestedLayouts(t *testing.T) {
	r := newTestResolver()
	r.Add("top_layout.html", "HEADER-{{template \"body\"}}-FOOTER")
	r.Add("sub_layout.html", "{{define \"body\"}}WRAP-{{template \"sub_body\"}}-WRAP{{end}}{{template \"top_layout.html\"}}")
	r.Add("page.html", "{{define \"sub_body\"}}BODY{{end}}{{template \"sub_layout.html\"}}")

	assertTextOutput(t, r, "page.html", "HEADER-WRAP-BODY-WRAP-FOOTER")
}
