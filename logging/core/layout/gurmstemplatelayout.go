package layout

import (
	"bytes"
	"errors"
	"fmt"
	"loggergo/infra/bufferpool"
	"loggergo/infra/cluster/nodetype"
	"loggergo/infra/infraerror"
	"loggergo/infra/lang"
	"loggergo/infra/timezone"
	"loggergo/logging/core/model/loglevel"
	"loggergo/mathsupport"
	"strings"
	"time"
)

var ESTIMATED_PATTERN_TEXT_LENGTH = 128
var LEVELS [][]byte
var NULL = []byte{'n', 'u', 'l', 'l'}
var COLON_SEPARATOR = []byte{' ', ':', ' '}
var TRACE_ID_LENGTH = 19
var STRUCT_NAME_LENGTH = 40

var NODE_TYPE_AI_SERVING = int('A')
var NODE_TYPE_GATEWAY = int('G')
var NODE_TYPE_SERVICE = int('S')
var NODE_TYPE_UNKNOWN = int('U')

func init() {
	var levels = []loglevel.LogLevel{0, 1, 2, 3, 4, 5}
	var levelCount = len(levels)
	LEVELS = make([][]byte, levelCount)
	maxLength := 0
	for _, level := range levels {
		maxLength = mathsupport.Max(len(level.String()), maxLength)
	}
	for i := 0; i < levelCount; i++ {
		level := lang.PadStart(levels[i].String(), maxLength, ' ')
		level = strings.ToUpper(level)
		var err error
		LEVELS[i], err = lang.GetBytes(level)
		if err != nil {
			println(err)
		}
	}

}

type GurmsTemplateLayout struct {
	nodeType int
	nodeId   []byte
}

func NewGurmsTemplateLayout(nodeType nodetype.NodeType, nodeId string) *GurmsTemplateLayout {
	var typ int
	noteId := nodeId
	switch nodeType {
	case nodetype.AI_SERVING:
		typ = NODE_TYPE_AI_SERVING
	case nodetype.GATEWAY:
		typ = NODE_TYPE_GATEWAY
	case nodetype.SERVICE:
		typ = NODE_TYPE_SERVICE
	default:
		typ = NODE_TYPE_UNKNOWN
	}
	return &GurmsTemplateLayout{
		nodeType: typ,
		nodeId:   []byte(noteId),
	}
}

func appendError(err error, buffer *bytes.Buffer) {
	buffer.WriteByte('\n')

	var errorMessage string

	level := infraerror.CountCauses(err)

	for err != nil {
		errorString := fmt.Sprintf("%s", err)

		err = errors.Unwrap(err)

		errorString2 := fmt.Sprintf("%s", err)
		length := len(errorString) - len(errorString2)
		if level == 0 {
			length = len(errorString)
		}
		errorString = errorString[0:length]

		errorMessage += fmt.Sprintf("level %d: %s", level, errorString)
		if err != nil {
			errorMessage += " | "
		}
		level--
	}

	buffer.Write([]byte(errorMessage))
}

func FormatBasic(layoutAL *GurmsTemplateLayout,
	structName []byte,
	level loglevel.LogLevel,
	msg *bytes.Buffer) *bytes.Buffer {
	buffer := bufferpool.BufferPool.Get().(*bytes.Buffer)

	return formatBasic0(layoutAL, buffer, structName, level, msg)
}

func Format(layoutAL *GurmsTemplateLayout,
	shouldParse bool,
	structName []byte,
	level loglevel.LogLevel,
	msg string,
	args []interface{},
	err error) *bytes.Buffer {
	estimatedErrorLength := 0
	if err != nil {
		causes := infraerror.CountCauses(err)
		if causes == 0 {
			estimatedErrorLength = 64
		} else {
			estimatedErrorLength = causes * 1024
		}
	}
	var estimatedLength int
	if msg == "" {
		estimatedLength = 0
	} else {
		estimatedLength = len(msg) + ESTIMATED_PATTERN_TEXT_LENGTH + estimatedErrorLength
	}
	if args != nil && shouldParse {
		estimatedLength += len(args) * 16
	}
	buffer := bufferpool.NewBufferWithLength(estimatedLength)

	return format0(buffer, layoutAL, shouldParse, structName, level, msg, args, err)
}

func format0(buffer *bytes.Buffer,
	layoutAL *GurmsTemplateLayout,
	shouldParse bool,
	structName []byte,
	level loglevel.LogLevel,
	msg string,
	args []interface{},
	err error) *bytes.Buffer {
	timestamp := timezone.ToBytes(time.Now())

	buffer.Write(timestamp)
	buffer.WriteByte(WHITESPACE)
	buffer.Write(LEVELS[level])
	buffer.WriteByte(WHITESPACE)
	buffer.WriteByte(byte(layoutAL.nodeType))
	buffer.WriteByte(WHITESPACE)
	buffer.Write(layoutAL.nodeId)
	buffer.WriteByte(WHITESPACE)

	if structName != nil {
		buffer.WriteByte(WHITESPACE)
		buffer.Write(structName)
	}
	buffer.Write(COLON_SEPARATOR)

	if msg != "" {
		appendMessage(shouldParse, msg, args, buffer)
	}

	if err != nil {
		appendError(err, buffer)
	}

	buffer.WriteByte('\n')
	return buffer
}

func formatBasic0(layoutAL *GurmsTemplateLayout,
	buffer *bytes.Buffer,
	structName []byte,
	level loglevel.LogLevel,
	msg *bytes.Buffer) *bytes.Buffer {
	timestamp := timezone.ToBytes(time.Now())

	buffer.Write(timestamp)
	buffer.WriteByte(WHITESPACE)
	buffer.Write(LEVELS[level])
	buffer.WriteByte(WHITESPACE)
	buffer.WriteByte(byte(layoutAL.nodeType))
	buffer.WriteByte(WHITESPACE)
	buffer.Write(layoutAL.nodeId)

	if structName != nil {
		buffer.WriteByte(WHITESPACE)
		buffer.Write(structName)
	}
	buffer.Write(COLON_SEPARATOR)

	msg.WriteByte('\n')

	buffer.Write(msg.Bytes())

	msg.Reset()
	bufferpool.BufferPool.Put(msg)

	return buffer
}

func appendMessage(shouldParse bool,
	msg string,
	args []interface{},
	buffer *bytes.Buffer) {
	if !shouldParse {
		buffer.Write([]byte(msg))
		return
	}
	var argCount int
	if args == nil {
		argCount = 0
	} else {
		argCount = len(args)
	}
	if argCount == 0 {
		buffer.Write([]byte(msg))
		return
	}
	argIndex := 0
	bytes := []byte(msg)
	length := len(bytes)
	for i := 0; i < length; i++ {
		b := bytes[i]
		if b == '{' && i < length-1 && bytes[i+1] == '}' {
			if argIndex < argCount {
				arg := args[argIndex]
				argIndex++
				str := fmt.Sprintf("%v", arg)
				buffer.Write([]byte(str))
			} else {
				buffer.Write(nil)
			}
			i++
		} else {
			buffer.WriteByte(b)
		}
	}
}

func FormatStructName(name string) []byte {
	rawBytes := []byte(name)
	if len(rawBytes) == STRUCT_NAME_LENGTH {
		return rawBytes
	}
	if len(rawBytes) > STRUCT_NAME_LENGTH {
		parts := lang.TokenizeToStringArray(name, ".")
		structName := parts[len(parts)-1]
		structNameLength := len(structName)
		if structNameLength >= STRUCT_NAME_LENGTH {
			return []byte(structName[:STRUCT_NAME_LENGTH])
		}
		result := make([]byte, STRUCT_NAME_LENGTH)
		writeIndex := STRUCT_NAME_LENGTH
		for i := len(parts) - 1; i >= 0; i-- {
			part := []byte(parts[i])
			if i == len(parts)-1 {
				writeIndex -= len(part)
				copy(result[writeIndex:], part)
			} else if writeIndex >= 2 {
				writeIndex -= 2
				result[writeIndex] = part[0]
				result[writeIndex+1] = '.'
			} else {
				break
			}
		}
		for i := 0; i < writeIndex; i++ {
			result[i] = ' '
		}
		return result
	}
	return []byte(lang.PadStart(name, STRUCT_NAME_LENGTH, ' '))
}
