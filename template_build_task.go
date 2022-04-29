package gotemplater

import (
	"github.com/tapvanvn/gotokenize/v2"
)

type TemplateBuildTask struct {
	template *Template
}

func (task *TemplateBuildTask) Process(tool interface{}) {

	if templateTool, ok := tool.(*TemplateBuildTool); ok {

		proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &task.template.Stream)
		proc.Context.BindingData = task.template.Context
		templateTool.HTML.Prepare(proc)

		tmpStream := gotokenize.CreateStream(0)
		for {
			token := templateTool.HTML.Next(proc)
			if token == nil {
				break
			}
			tmpStream.AddToken(*token)

			//Todo find instruction
		}
		//fmt.Println("--build--")
		//tmpStream.Debug(0, html.HTMLTokenNaming, html.HTMLDebugOption)
		//fmt.Println("--end build--")

		task.template.Stream = tmpStream
		task.template.IsReady = true
	}
}

func (task *TemplateBuildTask) ToolLabel() string {
	return "template_build"
}
