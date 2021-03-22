package html

import (
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

func (meaning *HTMLInstructionMeaning) Prepare(stream *gotokenize.TokenStream) {

	meaning.XMLHightMeaning.Prepare(stream)

	tmpStream := gotokenize.CreateStream()

	for {
		token := meaning.XMLHightMeaning.Next()

		if token == nil {

			break
		}
		if token.Type == xml.TokenXMLComment {

			continue
		}
		if token.Type == xml.TokenXMLElement {

			meaning.buildElement(token)
		}

		tmpStream.AddToken(*token)
	}
	meaning.SetStream(tmpStream)
}

func (meaning *HTMLInstructionMeaning) buildHead(token *gotokenize.Token) {

}
func (meaning *HTMLInstructionMeaning) buildElement(token *gotokenize.Token) {

}
