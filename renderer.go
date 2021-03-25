package gotemplater

import (
	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater/tokenize/html"
	"github.com/tapvanvn/gotokenize"
)

type Renderer struct {
}

func (r *Renderer) Compile(stream *gotokenize.TokenStream, context *gosmartstring.SSContext) (string, error) {

	iter := stream.Iterator()
	content := ""
	for {
		token := iter.Read()
		if token == nil {
			break
		}

		if token.Type == gosmartstring.TokenSSInstructionDo {

			buildContent, err := r.compileInstructionDo(token, context)

			if err != nil {

				return content, err
			}
			content += buildContent

		} else if token.Type == html.TokenOptimized {

			content += token.Content

		} else if token.Children.Length() > 0 {

			buildContent, err := r.Compile(&token.Children, context)
			if err != nil {
				return content, err
			}
			content += buildContent
		}
	}
	return content, nil
}

func (r *Renderer) compileInstructionDo(token *gotokenize.Token, context *gosmartstring.SSContext) (string, error) {

	iter := token.Children.Iterator()
	addressToken := iter.Get()
	obj := context.GetRegistry(addressToken.Content)
	if obj != nil && obj.Object != nil && obj.Object.CanExport() {
		return string(obj.Object.Export(context)), nil
	}
	return "", nil
}
