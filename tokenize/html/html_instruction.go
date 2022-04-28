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
	//xml.XMLHightMeaning
	SS gosmartstring.SmarstringInstructionMeaning
}

func CreateHTMLInstructionMeaning() HTMLInstructionMeaning {
	return HTMLInstructionMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(xml.NewXMLHighMeaning()),
		SS:              gosmartstring.CreateSSInstructionMeaning(),
	}
}

func (meaning *HTMLInstructionMeaning) Prepare(proc *gotokenize.MeaningProcess, context *gosmartstring.SSContext) {

	meaning.AbstractMeaning.Prepare(proc)

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
	proc.SetStream(proc.Context.AncestorTokens, &tmpStream)
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
				valueStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
				valueStream.Tokenize(valueContent)
				proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &valueStream)
				meaning.SS.Prepare(proc, context)

				tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

				for {
					ssToken := meaning.SS.Next(proc)
					if ssToken == nil {
						break
					}
					tmpStream.AddToken(*ssToken)
				}

				value.Type = gosmartstring.TokenSSLSmarstring
				value.Children = tmpStream
				//tmpStream.Debug(0, nil)
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
		tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
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

					valueStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
					valueStream.Tokenize(content)
					proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &valueStream)
					meaning.SS.Prepare(proc, context)

					gatherStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

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
					childToken.Type = gosmartstring.TokenSSLSmarstring
					childToken.Children = gatherStream
					//gatherStream.Debug(0, nil)
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
	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
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
	tmpChildren := gotokenize.CreateStream(meaning.GetMeaningLevel())

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
