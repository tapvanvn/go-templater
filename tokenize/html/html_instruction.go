package html

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tapvanvn/gosmartstring"

	"github.com/tapvanvn/gotokenize/v2"
	"github.com/tapvanvn/gotokenize/v2/xml"
)

type HTMLInstructionMeaning struct {
	*gotokenize.AbstractMeaning
	SS *gosmartstring.SmarstringInstructionMeaning
}

func CreateHTMLInstructionMeaning() *HTMLInstructionMeaning {
	return &HTMLInstructionMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(xml.NewXMLHighMeaning()),
		SS:              gosmartstring.CreateSSInstructionMeaning(),
	}
}

func (meaning *HTMLInstructionMeaning) Prepare(proc *gotokenize.MeaningProcess) {

	meaning.AbstractMeaning.Prepare(proc)

	context := proc.Context.BindingData.(*gosmartstring.SSContext)

	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

	for {
		token := meaning.AbstractMeaning.Next(proc)

		if token == nil {

			break
		}
		if token.Type == xml.TokenXMLComment {
			//remove comment
			continue
		}
		if token.Type == xml.TokenXMLElement || token.Type == xml.TokenXMLEndElement {

			if err := meaning.buildElement(token, &tmpStream, context); err != nil {
				//TODO report error
				fmt.Println(err.Error())
				continue
			}
		} else if token.Type == xml.TokenXMLString {
			continue
		} else if token.Type == 0 {
			continue
		} else {
			continue
		}

		tmpStream.AddToken(*token)
	}

	proc.SetStream(proc.Context.AncestorTokens, &tmpStream)

	// fmt.Println("--after html_instruction prepare---")
	// proc.Stream.Debug(0, HTMLTokenNaming, HTMLDebugOption)
	// fmt.Println("--end after html_instruction prepare---")
}

func (meaning *HTMLInstructionMeaning) buildHead(token *gotokenize.Token, context *gosmartstring.SSContext) {

	iter := token.Children.Iterator()
	for {

		keyToken := iter.Read()
		if keyToken == nil {
			break
		}
		oper := iter.Get()
		if oper != nil && oper.Content == "=" {

			iter.Read()
			value := iter.Read()

			if value != nil && value.Type == xml.TokenXMLString {
				content := value.Children.ConcatStringContent()
				if strings.Index(content, "{{") > -1 {

					valueContent := value.Content
					if value.Type == xml.TokenXMLString {
						valueContent = value.Children.ConcatStringContent()
					}
					valueStream := gotokenize.CreateStream(0)

					valueStream.Tokenize(valueContent)

					proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &valueStream)
					proc.Context.BindingData = context
					//context.DebugLevel = 1

					meaning.SS.Prepare(proc)

					tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

					for {
						ssToken := meaning.SS.Next(proc)
						if ssToken == nil {
							break
						}
						tmpStream.AddToken(*ssToken)
					}

					//gotokenize.DebugMeaning(meaning.SS)

					value.Type = gosmartstring.TokenSSLSmartstring
					value.Children = tmpStream
				}
			}
		}
	}
}

func (meaning *HTMLInstructionMeaning) buildElement(token *gotokenize.Token, stream *gotokenize.TokenStream, context *gosmartstring.SSContext) error {
	if strings.Index(HTMLInstructionTagName, ","+token.Content+",") > -1 {
		//instruction
		switch strings.ToLower(token.Content) {
		case "for":
			return meaning.buildInstructionFor(token, context)
		case "template":
			stream.AddToken(gotokenize.Token{Type: gosmartstring.TokenSSInstructionReset})
			return meaning.buildInstructionTemplate(token, context)
		case "ssscript":
			stream.AddToken(gotokenize.Token{Type: gosmartstring.TokenSSInstructionReset})
			return meaning.buildInstructionSSScript(token, context)
		}
	} else {
		iter := token.Children.Iterator()
		head := iter.Read()
		tmpStream := gotokenize.CreateStream(0)
		if head != nil {
			meaning.buildHead(head, context)
			tmpStream.AddToken(*head)
		}
		for {
			childToken := iter.Read()
			if childToken == nil {
				break
			}
			if childToken.Type == xml.TokenXMLElement || childToken.Type == xml.TokenXMLEndElement {
				if err := meaning.buildElement(childToken, &tmpStream, context); err != nil {
					//fmt.Println(err.Error())
					//TODO: is this safe to continue
					continue
				}
			} else if childToken.Type == 0 || childToken.Type == xml.TokenXMLString {
				content := childToken.Content
				if childToken.Type == xml.TokenXMLString {
					content = childToken.Children.ConcatStringContent()
				}
				testPos := strings.Index(content, "{{")
				if testPos > -1 {

					endPos := strings.Index(content, "}}")
					if endPos == -1 {
						meaning.continueSmartstring(iter, &content)
					}
					valueStream := gotokenize.CreateStream(0)
					valueStream.Tokenize(content)
					//fmt.Println("ss:", content)
					//context.DebugLevel = 1
					proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &valueStream)
					proc.Context.BindingData = context
					meaning.SS.Prepare(proc)

					gatherStream := gotokenize.CreateStream(0)

					for {
						ssToken := meaning.SS.Next(proc)
						if ssToken == nil {
							break
						}
						gatherStream.AddToken(*ssToken)
					}
					if childToken.Type == 0 {
						childToken.Content = ""
					}

					childToken.Type = gosmartstring.TokenSSLSmartstring
					childToken.Children = gatherStream
				}

			}
			tmpStream.AddToken(*childToken)
		}
		token.Children = tmpStream
	}
	return nil
}

func (meaning *HTMLInstructionMeaning) continueSmartstring(iter *gotokenize.Iterator, content *string) {
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Type == xml.TokenXMLSpace || childToken.Type == gotokenize.TokenSpace {
			continue
		} else if childToken.Type == xml.TokenXMLString {
			*content += childToken.Content + childToken.Children.ConcatStringContent() + childToken.Content
		} else if childToken.Type == gotokenize.TokenWord {
			testPos := strings.Index(childToken.Content, "}}")
			if testPos > -1 {
				*content += childToken.Content
				break
			}
		} else {

			break
		}
	}
	fmt.Println("found:", *content)
}

func (meaning *HTMLInstructionMeaning) buildInstructionTemplate(token *gotokenize.Token, context *gosmartstring.SSContext) error {

	//fmt.Println("build ins template with context:", context.ID())
	token.Type = gosmartstring.TokenSSInstructionDo
	token.Content = "template"
	tmpStream := gotokenize.CreateStream(0)
	iter := token.Children.Iterator()
	head := iter.Read()
	if head == nil || head.Type != xml.TokenXMLElementAttributes {
		return errors.New("syntax error")
	}
	output := gotokenize.Token{
		Type:    gosmartstring.TokenSSRegistry,
		Content: context.IssueAddress(),
	}
	tmpStream.AddToken(output)

	attIter := head.Children.Iterator()
	findID := false
	//head.Debug(10, xml.XMLNaming, nil)
	for {
		keyToken := attIter.Read()
		if keyToken == nil {
			break
		}
		oper := attIter.Get()
		if oper != nil && oper.Content == "=" {

			attIter.Read()
			valueToken := attIter.Read()

			if strings.ToLower(keyToken.Content) == "id" {

				address := context.IssueAddress()
				idToken := gotokenize.Token{
					Type:    gosmartstring.TokenSSRegistry,
					Content: address,
				}
				id := strings.TrimSpace(valueToken.Children.ConcatStringContent())
				//fmt.Println("register template id:", gotokenize.ColorName(id), "at address", gotokenize.ColorContent(address))
				context.RegisterObject(address, gosmartstring.CreateString(id))
				tmpStream.AddToken(idToken)
				findID = true

			}
		}
	}
	token.Children = tmpStream

	if !findID {
		return errors.New("template syntax error no id found")
	}
	return nil
}

func (meaning *HTMLInstructionMeaning) buildInstructionFor(token *gotokenize.Token, context *gosmartstring.SSContext) error {
	token.Type = gosmartstring.TokenSSInstructionEach

	//for content is loop pack name
	token.Content = ""
	tmpChildren := gotokenize.CreateStream(0)

	iter := token.Children.Iterator()
	head := iter.Read()
	if head == nil || head.Type != xml.TokenXMLElementAttributes {

		return errors.New("syntax error")
	}

	outputToken := gotokenize.Token{
		Type:    gosmartstring.TokenSSRegistry,
		Content: context.IssueAddress(),
	}
	elementToken := gotokenize.Token{
		Type: gosmartstring.TokenSSRegistry,
	}
	attIter := head.Children.Iterator()
	for {

		keyToken := attIter.Read()
		if keyToken == nil {
			break
		}
		oper := attIter.Get()
		if oper != nil && oper.Content == "=" {

			attIter.Read()
			valueToken := attIter.Read()

			if keyToken.Content == "each" {

				elementToken.Content = valueToken.Children.ConcatStringContent()

			} else if keyToken.Content == "in" {
				token.Content = valueToken.Children.ConcatStringContent()
			}
		}
	}
	tmpChildren.AddToken(outputToken)
	tmpChildren.AddToken(elementToken)

	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if err := meaning.buildElement(childToken, &tmpChildren, context); err != nil {
			//fmt.Println("err", err.Error())
			//TODO: is this safe to continue
			continue
		}
		tmpChildren.AddToken(*childToken)
	}

	token.Children = tmpChildren
	//token.Debug(0, HTMLTokenNaming, &gotokenize.DebugOption{
	//	ExtendTypeSize: 6,
	//})
	return nil
}

func (meaning *HTMLInstructionMeaning) buildInstructionSSScript(token *gotokenize.Token, context *gosmartstring.SSContext) error {
	fmt.Println("instruction ssscript")
	iter := token.Children.Iterator()

	content := "{{"
	for {
		childToken := iter.Read()

		if childToken == nil {
			break
		}
		if childToken.Type == xml.TokenXMLElementAttributes {
			continue
		}
		if childToken.Type == xml.TokenXMLString {
			content += childToken.Content + childToken.Children.ConcatStringContent() + childToken.Content
			continue
		}
		content += string(childToken.Content)
	}
	content += "}}"
	fmt.Println("sscontent:", content)

	valueStream := gotokenize.CreateStream(0)

	valueStream.Tokenize(content)

	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &valueStream)
	proc.Context.BindingData = context

	meaning.SS.Prepare(proc)

	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

	for {
		ssToken := meaning.SS.Next(proc)
		if ssToken == nil {
			break
		}
		tmpStream.AddToken(*ssToken)
	}

	token.Type = gosmartstring.TokenSSLSmartstring
	token.Children = tmpStream
	token.Content = ""
	// fmt.Println("--ssscript")
	// token.Debug(0, HTMLTokenNaming, HTMLDebugOption)
	// fmt.Println("--end ssscript")
	return nil
}
