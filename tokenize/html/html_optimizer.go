package html

import (
	"strings"

	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotokenize/v2"
	"github.com/tapvanvn/gotokenize/v2/xml"
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

func (meaning *HTMLOptmizerMeaning) Prepare(proc *gotokenize.MeaningProcess, context *gosmartstring.SSContext) {
	meaning.HTMLInstructionMeaning.Prepare(proc, context)
	//meaning.HTMLInstructionMeaning.GetStream().Debug(0, nil)
	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	for {
		token := meaning.HTMLInstructionMeaning.Next(proc)
		if token == nil {
			break
		}
		meaning.optimizeToken(token, &tmpStream)
	}

	content := ""
	tmpStream2 := gotokenize.CreateStream(meaning.GetMeaningLevel())
	iter2 := tmpStream.Iterator()
	for {
		token := iter2.Read()
		if token == nil {
			break
		}
		if token.Type == TokenOptimized {
			content += token.Content
		} else {
			if content != "" {
				tmpStream2.AddToken(gotokenize.Token{
					Type:    TokenOptimized,
					Content: content,
				})
				content = ""
			}
			tmpStream2.AddToken(*token)
		}
	}
	if content != "" {
		tmpStream2.AddToken(gotokenize.Token{
			Type:    TokenOptimized,
			Content: content,
		})
		content = ""
	}

	proc.SetStream(proc.Context.AncestorTokens, &tmpStream2)
}
func (meaning *HTMLOptmizerMeaning) optimizeStream(iter *gotokenize.Iterator, outStream *gotokenize.TokenStream) {

	for {
		token := iter.Read()
		if token == nil {
			break
		}
		meaning.optimizeToken(token, outStream)
	}
}

//return true if all child token is optmized, fale if atleast one child token is instruction
func (meaning *HTMLOptmizerMeaning) optimizeToken(token *gotokenize.Token, outStream *gotokenize.TokenStream) {

	if token.Type == xml.TokenXMLElement {
		outStream.AddToken(gotokenize.Token{
			Type:    TokenOptimized,
			Content: "<" + token.Content + " ",
		})
		iter := token.Children.Iterator()
		head := iter.Get()
		if head != nil && head.Type == xml.TokenXMLElementAttributes {
			headIter := head.Children.Iterator()
			meaning.optimizeStream(headIter, outStream)
			iter.Read()
		}
		outStream.AddToken(gotokenize.Token{
			Type:    TokenOptimized,
			Content: ">",
		})
		meaning.optimizeStream(iter, outStream)

		outStream.AddToken(gotokenize.Token{
			Type:    TokenOptimized,
			Content: "</" + token.Content + ">",
		})
	} else if token.Type == xml.TokenXMLEndElement {
		outStream.AddToken(gotokenize.Token{
			Type:    TokenOptimized,
			Content: "<" + token.Content + " ",
		})
		iter := token.Children.Iterator()
		head := iter.Get()
		if head != nil && head.Type == xml.TokenXMLElementAttributes {
			headIter := head.Children.Iterator()
			meaning.optimizeStream(headIter, outStream)
			iter.Read()
		}
		if strings.Index(HTMLSingleTagName, ","+strings.ToLower(token.Content)+",") != -1 {
			outStream.AddToken(gotokenize.Token{
				Type:    TokenOptimized,
				Content: ">",
			})
		} else {
			outStream.AddToken(gotokenize.Token{
				Type:    TokenOptimized,
				Content: "/>",
			})
		}

	} else if token.Type == xml.TokenXMLAttribute {

		childIter := token.Children.Iterator()
		key := childIter.Get()
		val := childIter.GetBy(1)
		if key != nil && val != nil {
			outStream.AddToken(gotokenize.Token{
				Type:    TokenOptimized,
				Content: " " + key.Content + "=",
			})
		} else if val == nil {
			outStream.AddToken(gotokenize.Token{
				Type:    TokenOptimized,
				Content: key.Content,
			})
		}
		if val != nil {
			if val.Type == 0 {
				outStream.AddToken(gotokenize.Token{
					Type:    TokenOptimized,
					Content: "\"" + val.Content + "\"",
				})
			} else if val.Type == xml.TokenXMLString {
				outStream.AddToken(gotokenize.Token{
					Type:    TokenOptimized,
					Content: val.Content + val.Children.ConcatStringContent() + val.Content,
				})
			} else {

				outStream.AddToken(gotokenize.Token{
					Type:    TokenOptimized,
					Content: "\"",
				})
				outStream.AddToken(*val)
				outStream.AddToken(gotokenize.Token{
					Type:    TokenOptimized,
					Content: "\"",
				})
			}
		}
	} else if token.Type == 0 || token.Type == xml.TokenXMLSpace || token.Type == xml.TokenXMLOperator {
		outStream.AddToken(gotokenize.Token{
			Type:    TokenOptimized,
			Content: token.Content,
		})
	} else if token.Type == xml.TokenXMLString {
		outStream.AddToken(gotokenize.Token{
			Type:    TokenOptimized,
			Content: token.Content + token.Children.ConcatStringContent() + token.Content,
		})
	} else if token.Type == gosmartstring.TokenSSLSmarstring && token.Content != "" {
		outStream.AddToken(gotokenize.Token{
			Type:    TokenOptimized,
			Content: token.Content,
		})
		outStream.AddToken(*token)
		outStream.AddToken(gotokenize.Token{
			Type:    TokenOptimized,
			Content: token.Content,
		})

	} else {

		outStream.AddToken(*token)
	}
}
