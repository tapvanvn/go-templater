package gotemplater

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
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
	ID           string //each teample has an id, the value is absolute path of that template file
	Path         []string
	Error        error
	IsReady      bool
	HostLanguage LanguageType
}

//GetRelativePath get
func (template *Template) GetRelativePath(path string) ([]string, error) {

	return GetAbsolutePath(template.Path, path)
}

func (template *Template) Render(context *Context) (string, error) {

	if !template.IsReady {

		if template.Error != nil {

			return "", template.Error
		}
		return "", errors.New("template error")
	}

	path := "/" + strings.Join(template.Path, "/")

	//MARK: init route
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	return string(bytes), nil
}
