package gotemplater

import (
	"fmt"
	"time"

	ss "github.com/tapvanvn/gosmartstring"
)

func SSFTemplate(context *ss.SSContext, input ss.IObject, params []ss.IObject) ss.IObject {

	if len(params) == 1 {

		if sstring, ok := params[0].(*ss.SSString); ok {

			id := sstring.Value

			templater := GetTemplater()
			template := templater.GetTemplate(id)

			for {

				if template.IsReady {

					break
				}

				time.Sleep(time.Nanosecond * 10)
			}

			err := template.build(context)
			if err != nil {
				if id == "test:html/todo_item.html" {
					fmt.Println("--begin debug item error--")
					fmt.Println(err.Error())

					fmt.Println("---debug item context--")
					context.PrintDebug(0)

					fmt.Println("--end debug item error--")
				} else {
					context.PrintDebug(0)
				}
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
