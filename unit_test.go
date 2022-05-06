package gotemplater_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/tapvanvn/gosmartstring"
	ss "github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater"
	"github.com/tapvanvn/gotemplater/tokenize/html"
	"github.com/tapvanvn/gotokenize/v2/xml"

	"github.com/tapvanvn/gotokenize/v2"
)

func init() {

	rootPath, _ := os.Getwd()

	templater := gotemplater.GetTemplater()
	templater.AddNamespace("test", rootPath+"/test")
}

var compiler = &ss.SSCompiler{}

func createTestContext() *ss.SSContext {
	runtime := gotemplater.CreateHTMLRuntime()
	runtime.RegisterFunction("put", SSFPut)
	runtime.RegisterFunction("print", SSFPrint)

	context := ss.CreateContext(runtime)

	//single todo
	context.RegisterObject("todo", gosmartstring.CreateString("todo 1"))

	dic := ss.CreateSSStringMap()
	dic.Set("x", ss.CreateString("x_value"))
	context.RegisterObject("dic", dic)

	array := gosmartstring.CreateSSArray()

	array.Stack = append(array.Stack, gosmartstring.CreateString("todo 1"))
	array.Stack = append(array.Stack, gosmartstring.CreateString("todo 2"))

	context.RegisterObject("todo_list", array)

	return context
}

func SSFPut(context *ss.SSContext, input ss.IObject, params []ss.IObject) ss.IObject {
	//fmt.Println("call put")
	if len(params) == 1 {
		//fmt.Println("call put param")
		if name, ok := params[0].(*ss.SSString); ok {
			//fmt.Println("call put param2")
			formatedName := strings.TrimSpace(name.Value)
			if formatedName != "" {
				//fmt.Println("put to ", formatedName)

				hot := context.HotObject()
				context.RegisterObject(formatedName, hot)
				if hot == nil {
					fmt.Println("put nil to ", formatedName)
				}
			}
		}
	}
	return nil
}

func SSFPrint(context *ss.SSContext, input ss.IObject, params []ss.IObject) ss.IObject {
	fmt.Printf("call ssfprint %d\n", len(params))
	for i, param := range params {
		if str, ok := param.(*ss.SSString); ok {
			fmt.Printf("ssfprint-%d: %s\n", i, str.Value)
		} else {
			fmt.Printf("ssfprint-%d: %s\n", i, str.GetType())
		}
	}
	return nil
}

func TestNamespace(t *testing.T) {

	templater := gotemplater.GetTemplater()

	path, err := templater.GetPath("test:html/index.html")
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("/" + strings.Join(path, "/"))

	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		t.Fatal(err)
	}

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(string(bytes))

	meaning := html.CreateHTMLInstructionMeaning()

	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &stream)
	proc.Context.BindingData = ss.CreateContext(gotemplater.CreateHTMLRuntime())
	meaning.Prepare(proc)

	token := meaning.Next(proc)

	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, "[", xml.XMLNaming(token.Type), "]", token.Content)
		if token.Children.Length() > 0 {
			token.Children.Debug(1, xml.XMLNaming, nil)
		}
		token = meaning.Next(proc)
	}

	fmt.Println(strings.Join(path, "/"))
}

func TestHTMLTokenize(t *testing.T) {

	context := createTestContext()

	instructionDo := ss.BuildDo("template",
		[]ss.IObject{ss.CreateString("test:html/todo_item.html")}, context)

	stream := gotokenize.CreateStream(0)
	stream.AddToken(instructionDo)

	err := compiler.Compile(&stream, context)
	if err != nil {
		fmt.Println(err.Error())
		context.PrintDebug(0)
	}

	fmt.Println("-----FINISH------")
	renderer := gotemplater.CreateRenderer()
	resultContent, err := renderer.Compile(&stream, context)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resultContent)
}
func TestHTMLSSScript(t *testing.T) {

	context := createTestContext()

	instructionDo := ss.BuildDo("template",
		[]ss.IObject{ss.CreateString("test:html/script.html")}, context)

	stream := gotokenize.CreateStream(0)

	stream.AddToken(instructionDo)

	err := compiler.Compile(&stream, context)

	if err != nil {

		fmt.Println(err.Error())
		context.PrintDebug(0)
	}

	fmt.Println("-----FINISH------")
	renderer := gotemplater.CreateRenderer()
	resultContent, err := renderer.Compile(&stream, context)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resultContent)
}

func TestInstructionTemplate(t *testing.T) {

	context := createTestContext()

	instructionDo := ss.BuildDo("template",
		[]ss.IObject{ss.CreateString("test:html/index.html")}, context)

	stream := gotokenize.CreateStream(0)
	stream.AddToken(instructionDo)

	err := compiler.Compile(&stream, context)
	if err != nil {
		fmt.Println(err.Error())
		context.PrintDebug(0)
	}

	fmt.Println("-----FINISH------")
	renderer := gotemplater.CreateRenderer()
	resultContent, err := renderer.Compile(&stream, context)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resultContent)

}

func printUtf8(content string) {
	for _, c := range content {
		fmt.Printf("%c", c)
	}
}
func TestInstructionTemplate2(t *testing.T) {

	templater := gotemplater.GetTemplater()

	context := ss.CreateContext(gotemplater.CreateHTMLRuntime())

	array := gosmartstring.CreateSSArray()
	array.Stack = append(array.Stack, gosmartstring.CreateString("todo 1"))
	array.Stack = append(array.Stack, gosmartstring.CreateString("todo 2"))

	dic := gosmartstring.CreateSSStringMap()
	dic.Set("x", ss.CreateString("x_value"))

	context.RegisterObject("todo_list", array)
	context.RegisterObject("dic", dic)

	resultContent, err := templater.Render("test:html/index.html", context)

	if err != nil {

		fmt.Println(err.Error())
	}

	//templater.ClearCache("test:html/index.html")
	//templater.ClearCache("testabc")
	//templater.ClearAllCache()

	printUtf8(resultContent)
	fmt.Println()

}

func TestInstructionTemplate3(t *testing.T) {

	templater := gotemplater.GetTemplater()

	context := ss.CreateContext(gotemplater.CreateHTMLRuntime())
	context.RegisterObject("todo", ss.CreateString("test_todo"))

	resultContent, err := templater.Render("test:html/todo_item.html", context)

	if err != nil {

		fmt.Println(err.Error())
	}
	//templater.ClearCache("test:html/index.html")
	//templater.ClearCache("testabc")
	//templater.ClearAllCache()

	printUtf8(resultContent)
	fmt.Println()

}

func TestSmartstring(t *testing.T) {

	content := `abc {{context.bcd}} f`

	context := gosmartstring.CreateContext(gotemplater.CreateHTMLRuntime())

	strmap := gosmartstring.CreateSSStringMap()
	strmap.Set("bcd", gosmartstring.CreateString("value"))
	context.Root.RegisterObject("context", strmap)

	meaning := gosmartstring.CreateSSInstructionMeaning()
	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &stream)
	proc.Context.BindingData = context
	meaning.Prepare(proc)

	compileStream := gotokenize.CreateStream(0)
	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		compileStream.AddToken(*token)
	}
	if err := compiler.Compile(&compileStream, context); err != nil {
		t.Fatal(err)
	}
	context.PrintDebug(0)
	/*renderer := gotemplater.CreateRenderer()
	content, err := renderer.Compile(&compileStream, context)
	if err != nil {

		fmt.Println(err.Error())
	}
	compileStream.Debug(0, gosmartstring.SSNaming, nil)
	fmt.Println(content)*/

	gotokenize.DebugMeaning(meaning)
}
