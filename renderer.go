package gotemplater

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater/tokenize/html"
	"github.com/tapvanvn/gotokenize"
)

type Renderer struct {
	id uuid.UUID
}

func CreateRenderer() Renderer {
	return Renderer{
		id: uuid.New(),
	}
}

func (r *Renderer) Compile(stream *gotokenize.TokenStream, context *gosmartstring.SSContext) (string, error) {

	iter := stream.Iterator()
	//stream.Debug(0, nil)
	return r.compileStream(&iter, context)
}
func (r *Renderer) compileStream(iter *gotokenize.Iterator, context *gosmartstring.SSContext) (string, error) {

	content := ""
	for {
		token := iter.Read()
		if token == nil {
			break
		}

		if token.Type == gosmartstring.TokenSSInstructionDo || token.Type == gosmartstring.TokenSSInstructionExport {

			buildContent, err := r.compileInstructionDo(token, context)

			if err != nil {

				return content, err
			}
			content += buildContent

		} else if token.Type == gosmartstring.TokenSSInstructionEach {

			buildContent, err := r.compileInstructionEach(token, context)

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
		} else {
			if token.Type == gosmartstring.TokenSSLNormalstring {
				content += token.Content
			}

		}
	}
	return content, nil
}
func (r *Renderer) compileInstructionEach(token *gotokenize.Token, context *gosmartstring.SSContext) (string, error) {

	iter := token.Children.Iterator()
	addressToken := iter.Read()
	_ = iter.Read()
	content := ""
	obj := context.GetRegistry(addressToken.Content)
	if obj != nil && obj.Object != nil {
		if addressStack, ok := obj.Object.(*gosmartstring.SSAddressStack); !ok {
			return "", errors.New("runtime error " + obj.Object.GetType())
		} else {
			context.SetStackRegistry(addressStack)

			i := 0
			offset := iter.Offset
			stackNum := addressStack.GetStackNum()

			for {
				if i >= stackNum {
					break
				}
				addressStack.SetStack(i)
				//context.DebugCurrentStack()
				iter.Seek(offset)
				renderContent, err := r.compileStream(&iter, context)
				if err != nil {
					context.SetStackRegistry(nil)
					return content, nil
				}
				content += renderContent
				i++
			}
			context.SetStackRegistry(nil)
		}
	}
	return content, nil
}

func (r *Renderer) compileInstructionDo(token *gotokenize.Token, context *gosmartstring.SSContext) (string, error) {

	iter := token.Children.Iterator()
	addressToken := iter.Get()
	if addressToken.Type == gosmartstring.TokenSSRegistryIgnore {
		return "", nil
	}
	//fmt.Println("do:", token.Content, "address:", addressToken.Content)

	obj := context.GetRegistry(addressToken.Content)

	if obj != nil && obj.Object != nil {

		//fmt.Println("found", obj.Object.GetType())

		if obj.Object.CanExport() {

			return string(obj.Object.Export(context)), nil
		}
	} /*else {
		fmt.Println("not found", addressToken.Content)
		//context.PrintDebug(0)
	}*/
	return "", nil
}
