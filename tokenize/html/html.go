package html

import (
	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotokenize/v2"
	"github.com/tapvanvn/gotokenize/v2/xml"
)

var HTMLInstructionTagName string = ",for,template,"
var HTMLSingleTagName string = ",!doctype,meta,"

func HTMLTokenNaming(tokenType int) string {
	if name := gosmartstring.SSNaming(tokenType); name != "unknown" {
		return name
	}
	if name := xml.XMLNaming(tokenType); name != "unknown" {
		return name
	}
	if tokenType == TokenOptimized {
		return "optimized"
	}
	return "unknown"
}

var HTMLDebugOption = &gotokenize.DebugOption{

	ExtendTypeSize: 6,
}
