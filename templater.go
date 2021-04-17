package gotemplater

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tapvanvn/gosmartstring"
	ss "github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater/utility"
	"github.com/tapvanvn/gotokenize"
	"github.com/tapvanvn/goworker"
)

//Templater manage
type templater struct {
	namespaces     map[string][]string //namespace to path
	loadedTemplate map[string]*Template
}

var Templater *templater = nil

//InitTemplater should call once at begining of app to init templater
func InitTemplater(numWorker int) error {
	if Templater == nil {

		gosmartstring.SSInsructionMove(5000)

		Templater = &templater{
			namespaces:     map[string][]string{},
			loadedTemplate: map[string]*Template{},
		}
		goworker.AddTool("template_build", &TemplateBuildBlacksmith{})

		if numWorker < 1 {
			goworker.OrganizeWorker(1)
		} else {
			goworker.OrganizeWorker(numWorker)
		}

		return nil
	}
	return errors.New("templater already init")
}
func GetTemplater() *templater {
	if Templater == nil {
		if err := InitTemplater(1); err != nil {
			panic(err)
		}
	}
	return Templater
}

//MARK: implement functions

func (tpt *templater) Debug() {

	fmt.Println("debug")
}

//AddNamespace add a namespace
func (tpt *templater) AddNamespace(namespace string, path string) error {

	segments := strings.Split(path, "/")
	resultSegments := []string{}
	for _, segment := range segments {
		if segment == "." {

		} else if segment == ".." {
			numSegment := len(resultSegments)
			if numSegment > 0 {
				resultSegments = resultSegments[0:numSegment]
			} else {
				return errors.New("path error")
			}
		} else if len(segment) > 0 {

			resultSegments = append(resultSegments, segment)
		}
	}
	tpt.namespaces[namespace] = resultSegments
	return nil
}

func (tpt *templater) GetPath(id string) ([]string, error) {
	relativePath := strings.TrimSpace(id)
	sep := strings.Index(relativePath, ":")

	namespace := ""
	if sep >= 0 {
		namespace = relativePath[0:sep]
		relativePath = relativePath[sep+1:]
	}
	nsPathSegments, ok := tpt.namespaces[namespace]
	if !ok {
		return nil, errors.New("namespace is not defined")
	}

	return utility.GetAbsolutePath(nsPathSegments, relativePath)
}

func (tpt *templater) Render(id string, context *gosmartstring.SSContext) (string, error) {

	instructionDo := ss.BuildDo("template",
		[]ss.IObject{ss.CreateString(id)}, context)

	stream := gotokenize.CreateStream()
	stream.AddToken(instructionDo)
	compiler := ss.SSCompiler{}
	err := compiler.Compile(&stream, context)
	if err != nil {
		fmt.Println(err.Error())
		context.PrintDebug(0)
	}

	renderer := CreateRenderer()
	return renderer.Compile(&stream, context)

}

func (tpt *templater) ClearAllCache() {

	tpt.loadedTemplate = map[string]*Template{}
}

func (tpt *templater) ClearCache(id string) {

	delete(tpt.loadedTemplate, id)
}

func (tpt *templater) GetTemplate(id string) *Template {

	if template, ok := tpt.loadedTemplate[id]; ok {
		//fmt.Println("template loaded")
		return template
	}
	template := CreateTemplate(id, TXT)

	tpt.loadedTemplate[id] = &template

	loadPath, err := tpt.GetPath(id)
	if err != nil {
		fmt.Println(err.Error())
		template.Error = err
		return &template
	}
	numSegment := len(loadPath)

	lastSegment := loadPath[numSegment-1]

	extPos := strings.LastIndex(lastSegment, ".")
	if extPos >= 0 {
		ext := strings.ToLower(lastSegment[extPos+1:])
		if ext == "html" || ext == "htm" {
			template.HostLanguage = HTML
		} else if ext == "js" {
			template.HostLanguage = JS
		} else if ext == "json" {
			template.HostLanguage = JSON
		}
	}
	template.Path = loadPath

	err = template.load()

	if err != nil {
		template.Error = err
		fmt.Println(err.Error())
		return &template
	}
	//call build template task
	task := &TemplateBuildTask{
		template: &template,
	}
	goworker.AddTask(task)

	return &template
}
