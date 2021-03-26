package html

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotokenize"
	"github.com/tapvanvn/gotokenize/xml"
)

type HTMLInstructionMeaning struct {
	xml.XMLHightMeaning
	SS gosmartstring.SmarstringInstructionMeaning
}

func CreateHTMLInstructionMeaning() HTMLInstructionMeaning {
	return HTMLInstructionMeaning{
		XMLHightMeaning: xml.CreateXMLMeaning(),
		SS:              gosmartstring.CreateSSInstructionMeaning(),
	}
}

func (meaning *HTMLInstructionMeaning) Prepare(stream *gotokenize.TokenStream, context *gosmartstring.SSContext) {

	meaning.XMLHightMeaning.Prepare(stream)

	tmpStream := gotokenize.CreateStream()

	for {
		token := meaning.XMLHightMeaning.Next()

		if token == nil {

			break
		}
		if token.Type == xml.TokenXMLComment {
			//remove comment
			continue
		}
		if token.Type == xml.TokenXMLElement || token.Type == xml.TokenXMLEndElement {

			if err := meaning.buildElement(token, context); err != nil {
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
	meaning.SetStream(tmpStream)
}

func (meaning *HTMLInstructionMeaning) buildHead(token *gotokenize.Token, context *gosmartstring.SSContext) {

	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Type != xml.TokenXMLAttribute {
			continue
		}
		childIter := childToken.Children.Iterator()
		_ = childIter.Get()
		value := childIter.GetAt(1)

		if value != nil && value.Type == xml.TokenXMLString {
			content := value.Children.ConcatStringContent()
			if strings.Index(content, "{{") > -1 {

				valueContent := value.Content
				if value.Type == xml.TokenXMLString {
					valueContent = value.Children.ConcatStringContent()
				}
				valueStream := gotokenize.CreateStream()
				valueStream.Tokenize(valueContent)

				meaning.SS.Prepare(&valueStream, context)

				tmpStream := gotokenize.CreateStream()

				for {
					ssToken := meaning.SS.Next()
					if ssToken == nil {
						break
					}
					tmpStream.AddToken(*ssToken)
				}

				value.Content = ""
				value.Type = gosmartstring.TokenSSLSmarstring
				value.Children = tmpStream
			}
		}
	}
}

func (meaning *HTMLInstructionMeaning) buildElement(token *gotokenize.Token, context *gosmartstring.SSContext) error {
	if strings.Index(HTMLInstructionTagName, ","+token.Content+",") > -1 {
		//instruction
		switch strings.ToLower(token.Content) {
		case "for":
			return meaning.buildInstructionFor(token, context)
		case "template":
			return meaning.buildInstructionTemplate(token, context)
		}
	} else {
		iter := token.Children.Iterator()
		head := iter.Read()
		tmpStream := gotokenize.CreateStream()
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
				if err := meaning.buildElement(childToken, context); err != nil {
					fmt.Println(err.Error())
					continue
				}
			} else if childToken.Type == 0 || childToken.Type == xml.TokenXMLString {
				content := childToken.Content
				if childToken.Type == xml.TokenXMLString {
					content = childToken.Children.ConcatStringContent()
				}

				if strings.Index(content, "{{") > -1 {

					valueStream := gotokenize.CreateStream()
					valueStream.Tokenize(content)

					meaning.SS.Prepare(&valueStream, context)

					gatherStream := gotokenize.CreateStream()

					for {
						ssToken := meaning.SS.Next()
						if ssToken == nil {
							break
						}
						gatherStream.AddToken(*ssToken)
					}

					childToken.Content = ""
					childToken.Type = gosmartstring.TokenSSLSmarstring
					childToken.Children = gatherStream
				}

			}
			tmpStream.AddToken(*childToken)
		}
		token.Children = tmpStream
	}
	return nil
}

func (meaning *HTMLInstructionMeaning) buildInstructionTemplate(token *gotokenize.Token, context *gosmartstring.SSContext) error {

	//fmt.Println("build ins template with context:", context.ID())
	token.Type = gosmartstring.TokenSSInstructionDo
	token.Content = "template"
	tmpStream := gotokenize.CreateStream()
	iter := token.Children.Iterator()
	head := iter.Read()
	if head == nil {
		return errors.New("syntax error")
	}
	output := gotokenize.Token{
		Type:    gosmartstring.TokenSSRegistry,
		Content: context.IssueAddress(),
	}
	tmpStream.AddToken(output)

	iter = head.Children.Iterator()
	findID := false
	for {
		headToken := iter.Read()
		if headToken == nil {
			break
		}

		if headToken.Type != xml.TokenXMLAttribute {

			continue
		}
		attIter := headToken.Children.Iterator()
		keyToken := attIter.Read()
		valueToken := attIter.Read()

		if keyToken != nil && valueToken != nil {
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

		} else if headToken.Content != "" {
			return errors.New("syntax error unknow attribute " + headToken.Content)
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
	tmpChildren := gotokenize.CreateStream()

	iter := token.Children.Iterator()
	head := iter.Read()
	if head == nil {
		return errors.New("syntax error")
	}
	headIter := head.Children.Iterator()
	outputToken := gotokenize.Token{
		Type:    gosmartstring.TokenSSRegistry,
		Content: context.IssueAddress(),
	}
	elementToken := gotokenize.Token{
		Type: gosmartstring.TokenSSRegistry,
	}
	for {
		headToken := headIter.Read()
		if headToken == nil {
			break
		}
		if headToken.Type != xml.TokenXMLAttribute {
			continue
		}
		attIter := headToken.Children.Iterator()
		keyToken := attIter.Read()
		valueToken := attIter.Read()

		if keyToken != nil && valueToken != nil {
			if keyToken.Content == "each" {
				elementToken.Content = valueToken.Children.ConcatStringContent()

			} else if keyToken.Content == "in" {
				token.Content = valueToken.Children.ConcatStringContent()
			} else if headToken.Content != "" {
				return errors.New("syntax error unknown" + headToken.Content)
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
		if err := meaning.buildElement(childToken, context); err != nil {
			continue
		}
		tmpChildren.AddToken(*childToken)
	}

	token.Children = tmpChildren
	return nil
}
