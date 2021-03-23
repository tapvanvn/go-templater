package smartstring

import (
	"fmt"

	ss "github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater"
	"github.com/tapvanvn/gotokenize"
)

var __cacheTemplates map[string]*ssTemplate = map[string]*ssTemplate{}

func SSFTemplate(context *ss.SSContext, input ss.IObject, params []ss.IObject) ss.IObject {

	if len(params) == 1 {

		if sstring, ok := params[0].(ss.SSString); ok {

			id := sstring.Value
			if template, ok := __cacheTemplates[id]; ok {
				return template
			}
			fmt.Print("id", id)
			templater := gotemplater.GetTemplater()
			templater.GetTemplate(id)
		}
	}
	return nil
}

type ssTemplate struct {
	ss.SSObject
	stream       gotokenize.TokenStream
	instructions []*gotokenize.Token
}

func CreateSSTemplate(stream *gotokenize.TokenStream) *ssTemplate {
	template := &ssTemplate{
		SSObject:     ss.SSObject{},
		stream:       gotokenize.CreateStream(),
		instructions: []*gotokenize.Token{},
	}

	return template
}
