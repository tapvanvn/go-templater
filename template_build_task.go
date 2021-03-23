package gotemplater

import (
	"github.com/tapvanvn/gotemplater/worker"
	"github.com/tapvanvn/gotokenize"
)

type TemplateBuildTask struct {
	template *Template
}

func (task *TemplateBuildTask) Process(tool interface{}) {
	if templateTool, ok := tool.(*worker.TemplateTool); ok {
		templateTool.HTML.Prepare(&task.template.Stream)
		tmpStream := gotokenize.CreateStream()
		for {
			token := templateTool.HTML.Next()
			if token == nil {
				break
			}
			tmpStream.AddToken(*token)

			//Todo find instruction
		}
		tmpStream.Debug(0, nil)
		task.template.Stream = tmpStream
		task.template.IsReady = true
	}
}

func (task *TemplateBuildTask) ToolLabel() string {
	return "template_tool"
}
