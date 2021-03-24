package gotemplater

import (
	"fmt"
	"time"

	ss "github.com/tapvanvn/gosmartstring"
)

func SSFTemplate(context *ss.SSContext, input ss.IObject, params []ss.IObject) ss.IObject {
	fmt.Println("buidl template:", len(params))
	if len(params) == 1 {

		if sstring, ok := params[0].(*ss.SSString); ok {

			id := sstring.Value

			fmt.Print("build template id:", id)
			templater := GetTemplater()
			template := templater.GetTemplate(id)

			for {

				if template.IsReady {

					break
				}

				time.Sleep(time.Nanosecond * 10)
			}
			template.build(context)
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
