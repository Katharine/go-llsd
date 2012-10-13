package llsd

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"container/list"
	"encoding/ascii85"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

var InvalidUnmarshalError = errors.New("Attempt to unmarshal into invalid object.")

type XMLLLSDParser struct {
	decoder *xml.Decoder
}

type Map map[string]interface{}
type Array []interface{}

func UnmarshalXML(data []byte, v interface{}) error {
	x := XMLLLSDParser{}
	return x.Unmarshal(data, v)
}

func (x *XMLLLSDParser) Unmarshal(data []byte, v interface{}) (err error) {
	// Convert non-fatal panics into returned errors.
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	buffer := bytes.NewBuffer(data)
	x.decoder = xml.NewDecoder(buffer)

	result := x.parseLLSD()
	switch t := result.(type) {
	case Map, Array:
		rv.Elem().Set(reflect.ValueOf(t))
	case Undef:
		switch rv.Elem().Kind() {
		case reflect.Slice:
			rv.Elem().Set(reflect.ValueOf(Array{}))
		case reflect.Map:
			rv.Elem().Set(reflect.ValueOf(Map{}))
		case reflect.Interface:
			rv.Elem().Set(reflect.ValueOf(t))
		default:
			panic("Can't set reasonably type!")
		}
	}
	return nil
}

func (x *XMLLLSDParser) parseLLSD() interface{} {
	// Get the root LLSD
root_search:
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "llsd" {
				break root_search
			} else {
				x.decoder.Skip()
			}
			break
		}
	}

	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.EndElement:
			if t.Name.Local == "llsd" {
				return Undef{}
			} else {
				panic("Logic error: </" + t.Name.Local + "> in parseLLSD()")
			}
		case xml.StartElement:
			switch t.Name.Local {
			case "map":
				return x.parseMap()
			case "array":
				return x.parseArray()
			case "undef":
				return x.parseUndef()
			default:
				panic("Unexpected element type " + t.Name.Local)
			}
		default:
			//panic(fmt.Sprintf("Unexpected XML token: %v", t))
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseBoolean() bool {
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			panic("Unexpected start element in <boolean>")
		case xml.CharData:
			defer x.decoder.Skip() // Skip to the end of this element.
			s := string(t)
			b, e := strconv.ParseBool(s)
			if e != nil {
				b = false
			}
			return b
		case xml.EndElement:
			return false
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseInteger() int {
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			panic("Unexpected start element in <integer>")
		case xml.CharData:
			defer x.decoder.Skip()
			s := string(t)
			i, e := strconv.ParseInt(s, 10, 32)
			if e != nil {
				i = 0
			}
			return int(i)
		case xml.EndElement:
			return 0
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseReal() float64 {
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			panic("Unexpected start element in <real>")
		case xml.CharData:
			defer x.decoder.Skip()
			s := string(t)
			r, e := strconv.ParseFloat(s, 64)
			if e != nil {
				r = 0
			}
			return r
		case xml.EndElement:
			return 0.0
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseString() string {
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			panic("Unexpected start element in <string>")
		case xml.CharData:
			defer x.decoder.Skip()
			s := string(t.Copy())
			return s
		case xml.EndElement:
			return ""
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseUUID() uuid.UUID {
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			panic("Unexpected start element in <uuid>")
		case xml.CharData:
			defer x.decoder.Skip()
			return uuid.Parse(string(t))
		case xml.EndElement:
			return uuid.NIL
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseDate() time.Time {
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			panic("Unexpected start element in <date>")
		case xml.CharData:
			defer x.decoder.Skip()
			d, err := time.Parse("2006-01-02T15:04:05Z", string(t))
			if err != nil {
				return d
			} else {
				return time.Unix(0, 0)
			}
		case xml.EndElement:
			return time.Unix(0, 0)
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseURI() string {
	return x.parseString()
}

func (x *XMLLLSDParser) parseBinary(encoding string) []byte {
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			panic("Unexpected start element in <binary>")
		case xml.CharData:
			defer x.decoder.Skip()
			switch encoding {
			case "base64", "":
				b64, _ := base64.StdEncoding.DecodeString(string(t))
				return b64
			case "base85":
				b85 := make([]byte, ascii85.MaxEncodedLen(len(t)))
				n, _, _ := ascii85.Decode(b85, t, true)
				return b85[:n]
			case "base16":
				output := make([]byte, len(t)/2)
				for i := 0; i < len(t); i += 2 {
					n, _ := strconv.ParseInt(string(t[i:i+2]), 16, 8)
					output[i/2] = byte(n)
				}
				return output
			}
		case xml.EndElement:
			return make([]byte, 0)
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseMap() Map {
	m := Map{}
	current_key := ""
	has_key := false
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			if !has_key {
				if t.Name.Local == "key" {
					current_key = x.parseString()
					has_key = true
				} else {
					panic("Execpted <key>, got <" + t.Name.Local + ">")
				}
			} else {
				m[current_key] = x.parseSomething(t)
				has_key = false
			}
		case xml.EndElement:
			if t.Name.Local == "map" {
				return m
			}
			panic("Closing element other than </map>! (" + t.Name.Local + ") - this is impossible.")
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseArray() Array {
	// Use a linked list to save on repeatedly reallocating our array
	l := list.New()
	for {
		token, err := x.decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			l.PushBack(x.parseSomething(t))
		case xml.EndElement:
			if t.Name.Local == "array" {
				// Turn our linked list into a real array
				a := make(Array, 0, l.Len())
				for e := l.Front(); e != nil; e = e.Next() {
					a = append(a, e.Value)
				}
				return a
			}
			panic("Closing element other than </map>! (" + t.Name.Local + ") - this is impossible.")
		}
	}
	panic("unreachable")
}

func (x *XMLLLSDParser) parseUndef() Undef {
	x.decoder.Skip()
	return Undef{}
}

func (x *XMLLLSDParser) parseSomething(t xml.StartElement) interface{} {
	switch t.Name.Local {
	case "boolean":
		return x.parseBoolean()
	case "integer":
		return x.parseInteger()
	case "real":
		return x.parseReal()
	case "string":
		return x.parseString()
	case "uuid":
		return x.parseUUID()
	case "date":
		return x.parseDate()
	case "uri":
		return x.parseURI()
	case "binary":
		encoding := ""
		for _, v := range t.Attr {
			if v.Name.Local == "encoding" {
				encoding = v.Value
				break
			}
		}
		return x.parseBinary(encoding)
	case "map":
		return x.parseMap()
	case "array":
		return x.parseArray()
	case "undef":
		return x.parseUndef()
	default:
		panic("Unexpected element <" + t.Name.Local + ">")
	}
	panic("unreachable")
}
