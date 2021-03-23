package gotemplater

import (
	"fmt"

	ss "github.com/tapvanvn/gosmartstring"
)

func SSFTemplate(context *ss.SSContext, input ss.IObject, params []ss.IObject) ss.IObject {

	if len(params) == 1 {

		if sstring, ok := params[0].(ss.SSString); ok {

			id := sstring.Value

			fmt.Print("id", id)
			templater := GetTemplater()
			templater.GetTemplate(id)
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
