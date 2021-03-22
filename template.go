package gotemplater

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	ss "github.com/tapvanvn/gosmartstring"
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

	stream       gotokenize.TokenStream
	instructions []*gotokenize.Token
}

func CreateTemplate(id string, hostLanguage LanguageType) Template {
	return Template{
		IObject:      ss.SSObject{},
		ID:           id,
		IsReady:      false,
		Error:        nil,
		HostLanguage: hostLanguage,
		stream:       gotokenize.CreateStream(),
		instructions: []*gotokenize.Token{},
	}
}

//GetRelativePath get
func (template *Template) GetRelativePath(path string) ([]string, error) {

	return GetAbsolutePath(template.Path, path)
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

	template.stream.Tokenize(string(bytes))
	return nil
}
