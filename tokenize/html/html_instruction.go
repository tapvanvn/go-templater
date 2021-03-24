package html

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotokenize"
	"github.com/tapvanvn/gotokenize/xml"
)

type HTMLInstructionMeaning struct {
	xml.XMLHightMeaning
}

func CreateHTMLInstructionMeaning() HTMLInstructionMeaning {
	return HTMLInstructionMeaning{
		XMLHightMeaning: xml.CreateXMLMeaning(),
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
		if token.Type == xml.TokenXMLElement {

			if err := meaning.buildElement(token, context); err != nil {
				//TODO report error
				fmt.Println(err.Error())
				continue
			}
		}

		tmpStream.AddToken(*token)
	}
	meaning.SetStream(tmpStream)
}

func (meaning *HTMLInstructionMeaning) buildHead(token *gotokenize.Token) {

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
			meaning.buildHead(head)
			tmpStream.AddToken(*head)
		}
		for {
			childToken := iter.Read()
			if childToken == nil {
				break
			}
			if err := meaning.buildElement(childToken, context); err != nil {
				fmt.Println(err.Error())
				continue
			}
			tmpStream.AddToken(*childToken)
		}
		token.Children = tmpStream
	}
	return nil
}

func (meaning *HTMLInstructionMeaning) buildInstructionTemplate(token *gotokenize.Token, context *gosmartstring.SSContext) error {

	token.Type = gosmartstring.TokenSSInstructionDo
	token.Content = "template"
	tmpStream := gotokenize.CreateStream()
	iter := token.Children.Iterator()
	head := iter.Read()
	if head == nil {
		return errors.New("syntax error")
	}
	output := gotokenize.Token{
		Type:    gosmartstring.TokenSSRegistryGlobal,
		Content: context.Root.IssueAddress(),
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
				address := uuid.New().String()
				idToken := gotokenize.Token{
					Type:    gosmartstring.TokenSSRegistry,
					Content: address,
				}
				id := strings.TrimSpace(valueToken.Children.ConcatStringContent())
				fmt.Println("regist id:", id, "address", address)
				context.Root.RegisterObject(address, gosmartstring.CreateString(id))
				tmpStream.AddToken(idToken)
				findID = true
			}

		} else if headToken.Content != "" {
			return errors.New("syntax error unknow attribute " + headToken.Content)
		}
		fmt.Println(headToken.Content, headToken.Children.ConcatStringContent())
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
		Type:    gosmartstring.TokenSSRegistryGlobal,
		Content: context.Root.IssueAddress(),
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
		if headToken.Content == "each" {
			elementToken.Content = headToken.Content
		} else if headToken.Content == "in" {
			token.Content = headToken.Content
		} else if headToken.Content != "" {
			return errors.New("syntax error unknown" + headToken.Content)
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
