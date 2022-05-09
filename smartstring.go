package gotemplater

import (
	ss "github.com/tapvanvn/gosmartstring"
)

func SSFTemplate(context *ss.SSContext, input ss.IObject, params []ss.IObject) ss.IObject {

	if len(params) == 1 {

		if sstring, ok := params[0].(*ss.SSString); ok {

			id := sstring.Value

			templater := __templater
			template := templater.GetTemplate(id)
			if template.Error != nil {

				//context.PrintDebug(0)
				return ss.CreateSSError(0, template.Error.Error())
			}
			err := template.build(context)
			if err != nil {

				//context.PrintDebug(0)
				return ss.CreateSSError(0, err.Error())
			}
			return template
		} else {
			return ss.CreateSSError(0, "load template fail bad param")
		}
	}
	return ss.CreateSSError(0, "load template fail bad param (need only one string param)")
}

//CreateHTMLRuntime create html runtime
func CreateHTMLRuntime() *ss.SSRuntime {
	htmlRuntime := ss.CreateRuntime(nil)
	htmlRuntime.RegisterFunction("template", SSFTemplate)
	return htmlRuntime
}
