package layout

import "bytes"

const DEPTH_LIMIT = 16
const WHITESPACE = ' '

var AT = []byte("at ")
var CAUSED_BY = []byte("caused by: ")
var COLON = []byte(": ")
var CYCLIC_EXCEPTION = []byte(">>(Cyclic Exception?)>>")
var NATIVE = []byte("native")
var SUPPRESSED = []byte("suppressed: ")
var UNKNOWN = []byte("unknown")

func pad(buffer bytes.Buffer, minLength int) {
	for i := 0; i < minLength; i++ {
		buffer.WriteByte(WHITESPACE)
	}
}

func padStart(buffer bytes.Buffer, bytes []byte, minLength int) {

}
