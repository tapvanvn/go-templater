package gotemplater

import (
	ss "github.com/tapvanvn/gosmartstring"
)

func SSFTemplate(context *ss.SSContext, input ss.IObject, params []ss.IObject) ss.IObject {

	if len(params) == 1 {

		if sstring, ok := params[0].(*ss.SSString); ok {

			id := sstring.Value

			templater := GetTemplater()
			template := templater.GetTemplate(id)

			err := template.build(context)
			if err != nil {

				context.PrintDebug(0)

			}
			return template
		}
	}
	return nil
}

//CreateHTMLRuntime create html runtime
func CreateHTMLRuntime() *ss.SSRuntime {
	htmlRuntime := ss.CreateRuntime(nil)
	htmlRuntime.RegisterFunction("template", SSFTemplate)
	return htmlRuntime
}
