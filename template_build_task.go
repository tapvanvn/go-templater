package gotemplater

import (
	"github.com/tapvanvn/gotokenize"
)

type TemplateBuildTask struct {
	template *Template
}

func (task *TemplateBuildTask) Process(tool interface{}) {
	if templateTool, ok := tool.(*TemplateTool); ok {
		templateTool.HTML.Prepare(&task.template.Stream, task.template.Context)
		tmpStream := gotokenize.CreateStream()
		for {
			token := templateTool.HTML.Next()
			if token == nil {
				break
			}
			tmpStream.AddToken(*token)

			//Todo find instruction
		}
		//tmpStream.Debug(0, nil)
		task.template.Stream = tmpStream
		task.template.IsReady = true
	}
}

func (task *TemplateBuildTask) ToolLabel() string {
	return "template_tool"
}
