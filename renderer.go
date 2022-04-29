package gotemplater

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater/tokenize/html"

	"github.com/tapvanvn/gotokenize/v2"
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

	return r.compileStream(iter, context)
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
	//#0 : address
	//#1 :
	//#2 : element name

	//fmt.Println("--begin render each --")
	//token.Debug(0, html.HTMLTokenNaming, html.HTMLDebugOption)
	//fmt.Println("--end render each --")

	iter := token.Children.Iterator()
	addressToken := iter.Read() //0
	_ = iter.Read()             // 1
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
				if i >= stackNum-1 {
					break
				}
				addressStack.SetStack(i)

				iter.Seek(offset)
				renderContent, err := r.compileStream(iter, context)
				if err != nil {
					fmt.Println("render err:", err.Error())
					context.SetStackRegistry(nil)
					return content, nil
				}
				content += renderContent
				i++
			}
			context.SetStackRegistry(nil)
		}
	} else {
		fmt.Println("loop with empty array")
	}
	return content, nil
}

func (r *Renderer) compileInstructionDo(token *gotokenize.Token, context *gosmartstring.SSContext) (string, error) {

	iter := token.Children.Iterator()
	addressToken := iter.Get()
	if addressToken.Type == gosmartstring.TokenSSRegistryIgnore {
		return "", nil
	}

	obj := context.GetRegistry(addressToken.Content)

	if obj != nil && obj.Object != nil {

		if obj.Object.CanExport() {

			return string(obj.Object.Export(context)), nil
		}
	}
	return "", nil
}
