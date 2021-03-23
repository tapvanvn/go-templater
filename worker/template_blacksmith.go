package worker

import "github.com/tapvanvn/gotemplater/tokenize/html"

type TemplateBlacksmith struct {
}

//Make make tool
func (blacksmith *TemplateBlacksmith) Make() interface{} {

	tool := &TemplateTool{
		HTML: html.CreateHTMLInstructionMeaning(),
	}
	return tool
}
