package gotemplater

import (
	"errors"
	"io/ioutil"
	"os"
)

type Template struct {
	ID      string //each teample has an id, the value is absolute path of that template file
	Path    string
	Error   error
	IsReady bool
}

func (template *Template) Render(context *Context) (string, error) {
	if !template.IsReady {
		if template.Error != nil {
			return "", template.Error
		}
		return "", errors.New("template error")
	}
	//MARK: init route
	file, err := os.Open(template.Path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	return string(bytes), nil
}
