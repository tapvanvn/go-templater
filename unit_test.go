package gotemplater_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ss "github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater"
	"github.com/tapvanvn/gotemplater/smartstring"
	html "github.com/tapvanvn/gotemplater/tokenize/html"
	"github.com/tapvanvn/gotokenize"
	"github.com/tapvanvn/gotokenize/xml"
)

func TestNamespace(t *testing.T) {

	rootPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	gotemplater.InitTemplater(1)
	templater := gotemplater.GetTemplater()

	templater.AddNamespace("test", rootPath+"/test")
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

	stream := gotokenize.CreateStream()
	stream.Tokenize(string(bytes))

	meaning := html.CreateHTMLInstructionMeaning()
	meaning.Prepare(&stream)

	token := meaning.Next()

	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, "[", xml.XMLNaming(token.Type), "]", token.Content)
		if token.Children.Length() > 0 {
			token.Children.Debug(1, xml.XMLNaming)
		}
		token = meaning.Next()
	}

	fmt.Println(strings.Join(path, "/"))
}

func TestInstructionTemplate(t *testing.T) {

	rootPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	gotemplater.InitTemplater(1)

	templater := gotemplater.GetTemplater()
	templater.AddNamespace("test", rootPath+"/test")

	context := ss.CreateContext(smartstring.HTMLRuntime)

	instructionDo := ss.BuildInstructionDo("template",
		[]ss.IObject{ss.CreateString("test:html/index.html")}, context)

	compiler := ss.SSCompiler{}
	compiler.Compile(&instructionDo, context)

	time.Sleep(time.Second * 2)
}
