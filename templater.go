package gotemplater

import (
	"errors"
	"fmt"
	"strings"
)

//Templater manage
type templater struct {
	namespaces     map[string][]string //namespace to path
	loadedTemplate map[string]*Template
}

var Templater *templater = nil

//InitTemplater should call once at begining of app to init templater
func InitTemplater() *templater {
	if Templater == nil {
		Templater = &templater{
			namespaces:     map[string][]string{},
			loadedTemplate: map[string]*Template{},
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
	mergeSegments := append(nsPathSegments)

	segments := strings.Split(relativePath, "/")
	for _, segment := range segments {

		if segment == "." {

		} else if segment == ".." {
			numSegment := len(mergeSegments)
			if numSegment > 0 {
				mergeSegments = mergeSegments[0:numSegment]
			} else {
				return nil, errors.New("path error")
			}
		} else if len(segment) > 0 {
			mergeSegments = append(mergeSegments, segment)
		}
	}

	return mergeSegments, nil
}

func (tpt *templater) Render(id string, context *Context) (string, error) {

	template := tpt.GetTemplate(id)

	if template.IsReady {

		content, err := template.Render(context)
		if err != nil {
			return "", errors.New("render error")
		}
		return content, nil
	}
	return "", errors.New("renderer not ready")
}

func (tpt *templater) GetTemplate(id string) *Template {

	if template, ok := tpt.loadedTemplate[id]; ok {
		return template
	}
	template := &Template{
		ID:      id,
		IsReady: false,
		Error:   nil,
	}
	tpt.loadedTemplate[id] = template

	loadPath, err := tpt.GetPath(id)
	if err != nil {
		template.Error = err
		return template
	}
	template.Path = "/" + strings.Join(loadPath, "/")
	template.IsReady = true
	return template
}
