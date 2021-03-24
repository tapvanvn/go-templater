package html

import (
	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotokenize"
)

const (
	TokenOptimized = 2000
)

type HTMLOptmizerMeaning struct {
	HTMLInstructionMeaning
}

func CreateHTMLOptmizer() HTMLOptmizerMeaning {
	return HTMLOptmizerMeaning{
		HTMLInstructionMeaning: CreateHTMLInstructionMeaning(),
	}
}

func (meaning *HTMLOptmizerMeaning) Prepare(stream *gotokenize.TokenStream, context *gosmartstring.SSContext) {
	meaning.HTMLInstructionMeaning.Prepare(stream, context)

}
