package gotemplater

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tapvanvn/gosmartstring"
	ss "github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater/utility"
	"github.com/tapvanvn/gotokenize"
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
		ID:           id,
		Path:         []string{},
		Error:        nil,
		IsReady:      false,
		HostLanguage: hostLanguage,
		Stream:       gotokenize.CreateStream(),
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
		return errors.New("template error")
	}

	path := "/" + strings.Join(template.Path, "/")

	//MARK: init route
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	stream := gotokenize.CreateStream()
	stream.Tokenize(string(bytes))

	template.Stream.Tokenize(string(bytes))
	return nil
}
