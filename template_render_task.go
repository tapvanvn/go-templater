package gotemplater

type TemplateRenderTask struct {
	template *Template
}

func (task *TemplateRenderTask) Process(tool interface{}) {
	/*if templateTool, ok := tool.(*TemplateBuildTool); ok {
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
	}*/
}

func (task *TemplateRenderTask) ToolLabel() string {
	return "template_render"
}
