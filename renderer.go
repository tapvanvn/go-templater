package gotemplater

import (
	"errors"
	"fmt"

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
	return r.compileStream(&iter, context)
}
func (r *Renderer) compileStream(iter *gotokenize.Iterator, context *gosmartstring.SSContext) (string, error) {

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
			for stactId, stack := range addressStack.Address {
				for address, translate := range stack {
					fmt.Println("stack:", stactId, address, translate)
				}
			}
			defer context.SetStackRegistry(nil)
			i := 0
			offset := iter.Offset
			stackNum := addressStack.GetStackNum()

			for {
				if i >= stackNum {
					break
				}
				addressStack.SetStack(i)
				context.DebugCurrentStack()
				iter.Seek(offset)
				renderContent, err := r.compileStream(&iter, context)
				if err != nil {
					return content, nil
				}
				content += renderContent
				i++
			}
		}
	}
	return content, nil
}

func (r *Renderer) compileInstructionDo(token *gotokenize.Token, context *gosmartstring.SSContext) (string, error) {

	iter := token.Children.Iterator()
	addressToken := iter.Get()

	fmt.Println("do:", token.Content, "address:", addressToken.Content)

	obj := context.GetRegistry(addressToken.Content)

	if obj != nil && obj.Object != nil {

		fmt.Println("found", obj.Object.GetType())

		if obj.Object.CanExport() {

			return string(obj.Object.Export(context)), nil
		}
	} else {
		fmt.Println("not found")
		context.PrintDebug(0)
	}
	return "", nil
}
