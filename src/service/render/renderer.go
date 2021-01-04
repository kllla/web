package render

import (
	"html/template"
	"io"
	"log"
)

type Renderer interface {
	GetPageTemplate(pageTitle string) *template.Template
	RenderPageWithOptions(w io.Writer, page Page, options *PageRenderOptions)
}

type PageRenderOptions struct {
	Banner        string
	AuthedNavBar  bool
	AuthedContent bool
	Indata        []interface{}
}

type renderImpl struct {
	Pages     map[string]Page
	Templates map[string]*template.Template
}

func NewRender() Renderer {
	templates := make(map[string]*template.Template)
	renderer := &renderImpl{Pages: MakePages(), Templates: templates}
	renderer.renderPageTemplates()
	return renderer
}

// GetPageTemplate returns the populated template
func (r *renderImpl) GetPageTemplate(pageTitle string) *template.Template {
	return r.Templates[pageTitle]
}

// GetPageTemplate returns the populated template
func (r *renderImpl) RenderPageWithOptions(w io.Writer, page Page, options *PageRenderOptions) {
	t := r.GetPageTemplate(page.GetTemplateTitle())
	err := t.Execute(w, page.RenderPageDataWithOptions(options))
	if err != nil {
		log.Print(err.Error())
	}
}

func (r *renderImpl) renderPageTemplates() {
	for _, page := range r.Pages {
		t := page.ToTemplate()
		r.Templates[page.GetTemplateTitle()] = t
	}
}
