package gotemplater

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/tapvanvn/gosmartstring"
	ss "github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater/utility"
	"github.com/tapvanvn/gotokenize/v2"
)

type LanguageType = int

const (
	HTML LanguageType = iota
	JS
	CSS
	JSON
	TXT
)

type Template struct {
	ss.IObject

	ID           string //each teample has an id, the value is absolute path of that template file
	Path         []string
	uid          uuid.UUID
	Error        error
	IsReady      bool
	HostLanguage LanguageType
	Stream       gotokenize.TokenStream
	Context      *gosmartstring.SSContext
	instructions []*gotokenize.Token
}

func CreateTemplate(id string, hostLanguage LanguageType) Template {
	return Template{
		IObject:      ss.SSObject{},
		uid:          uuid.New(),
		ID:           id,
		Path:         []string{},
		Error:        nil,
		IsReady:      false,
		HostLanguage: hostLanguage,
		Stream:       gotokenize.CreateStream(0),
		Context:      gosmartstring.CreateContext(CreateHTMLRuntime()),
		instructions: []*gotokenize.Token{},
	}
}

//GetRelativePath get
func (template *Template) GetRelativePath(path string) ([]string, error) {

	return utility.GetAbsolutePath(template.Path, path)
}

func (template *Template) load() error {

	if !template.IsReady {

		if template.Error != nil {

			return template.Error
		}
	}

	path := "/" + strings.Join(template.Path, "/")

	//MARK: init route
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(string(bytes))

	template.Stream.Tokenize(string(bytes))
	return nil
}

func (template *Template) build(context *gosmartstring.SSContext) error {

	compiler := ss.SSCompiler{}

	template.Context.BindingTo(context)

	err := compiler.Compile(&template.Stream, template.Context)

	//fmt.Println("--DEBUG CONTEXT BUILD--", template.ID, template.uid)
	//template.Context.PrintDebug(0)
	template.Context.BindingTo(nil)
	return err
}

func (template Template) CanExport() bool {
	return true
}

func (template Template) Export(context *gosmartstring.SSContext) []byte {

	template.Context.BindingTo(context)
	//template.Context.PrintDebug(0)
	var content = ""
	renderer := CreateRenderer()
	content, err := renderer.Compile(&template.Stream, template.Context)
	if err != nil {
		//TODO: report error
		fmt.Println(err.Error())
	}
	//fmt.Println("--DEBUG CONTEXT RENDER-- ", template.ID, template.uid)
	//template.Context.PrintDebug(0)
	template.Context.BindingTo(nil)
	return []byte(content)
}

func (obj Template) GetType() string {
	return "template" + obj.ID
}
