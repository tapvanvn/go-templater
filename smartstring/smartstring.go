package smartstring

import (
	ss "github.com/tapvanvn/gosmartstring"
)

var HTMLRuntime *ss.SSRuntime = createHTMLRuntime()

func createHTMLRuntime() *ss.SSRuntime {
	htmlRuntime := ss.CreateRuntime(nil)
	htmlRuntime.RegisterFunction("template", SSFTemplate)
	return htmlRuntime
}
