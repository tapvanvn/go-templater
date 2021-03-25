package gotemplater

import (
	"github.com/tapvanvn/gotemplater/tokenize/html"
)

type TemplateBuildBlacksmith struct {
}

//Make make tool
func (blacksmith *TemplateBuildBlacksmith) Make() interface{} {

	tool := &TemplateBuildTool{
		HTML: html.CreateHTMLOptmizer(),
	}
	return tool
}

type TemplateRenderBlacksmith struct {
}

//Make make tool
func (blacksmith *TemplateRenderBlacksmith) Make() interface{} {

	tool := &TemplateRenderTool{}
	return tool
}
