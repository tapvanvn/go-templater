package main

import (
	"fmt"
	"os"

	"github.com/tapvanvn/gosmartstring"
	ss "github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater"
)

func printUtf8(content string) {
	for _, c := range content {
		fmt.Printf("%c", c)
	}
}
func main() {
	rootPath, _ := os.Getwd()
	//rootPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	gotemplater.InitTemplater(1)

	templater := gotemplater.GetTemplater()
	templater.AddNamespace("test", rootPath+"/test")

	context := ss.CreateContext(gotemplater.CreateHTMLRuntime())
	array := gosmartstring.CreateSSArray()

	array.Stack = append(array.Stack, gosmartstring.CreateString("todo 1"))
	array.Stack = append(array.Stack, gosmartstring.CreateString("todo 2"))

	context.RegisterObject("todo_list", array)

	resultContent, err := templater.Render("test:html/index.html", context)

	if err != nil {

		fmt.Println(err.Error())
	}
	templater.ClearCache("test:html/index.html")
	templater.ClearCache("testabc")
	templater.ClearAllCache()

	printUtf8(resultContent)
	fmt.Println()
	//stream.Debug(0, nil)

}
