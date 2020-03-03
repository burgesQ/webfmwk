package pretty

import (
	"bufio"
	"bytes"
	"io"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

const (
	_bodyOpen     = '{'
	_bodyClose    = '}'
	_bracketOpen  = '['
	_bracketClose = ']'
	_jump         = "\n"
	_space        = " "
	_doubleSpace  = "  "
	_coma         = ","
	_amp          = `"`
	_next         = `":`
)

type (

	// JSON implement blablabla
	JSON struct {
		input  io.Reader
		output *bufio.Writer
		iter   *jsoniter.Iterator
		wg     sync.WaitGroup
		err    error
		comp   bool
	}

	// JsonError struct {
	// 	e       error
	// 	context string
	// }
)

// func (e JsonError) Error() string {
// 	return e.context + ": " + e.e.Error()
// }

// func NewJSONErr(e error) {
// 	if e != nil {
// 		panic(JsonError{e, "cannot write to output"})
// 	}
// }

// Use the pretty json utilitary to create well, pretty json ? :nerd_face:
func SimplePrettyJSON(r io.Reader, pretty bool) (string, error) {
	var (
		o  = new(bytes.Buffer)
		pj = NewPrettyJSON(r, o)
	)

	if !pretty {
		pj.SetCompactMode()
	}

	pj.Start()

	if err := pj.Close(); err != nil {
		return "", err
	}

	return o.String(), nil
}

// NewPrettyJSON return a instacied JSON
func NewPrettyJSON(in io.Reader, out io.Writer) JSON {
	return JSON{
		input:  in,
		output: bufio.NewWriterSize(out, 1024),
		iter:   jsoniter.Parse(jsoniter.ConfigFastest, in, 512*1024),
		wg:     sync.WaitGroup{},
		err:    nil,
		comp:   false,
	}
}

// Start launch the JSON parsing
func (pj *JSON) Start() {
	pj.wg.Add(1)
	go pj.start()
}

// Close stop the JSON parsing
func (pj *JSON) Close() error {
	pj.wg.Wait()
	return pj.err
}

// SetCompactMode render a compact JSON
func (pj *JSON) SetCompactMode() {
	pj.comp = true
}

func (pj *JSON) start() {
	defer pj.wg.Done()
	pj.run()

	closer, ok := pj.input.(io.Closer)
	if ok {
		pj.err = closer.Close()
	}
}

func (pj *JSON) printSimpleValue() {
	bts := pj.iter.SkipAndReturnBytes()
	_, pj.err = pj.output.Write(bts)
}

func (pj *JSON) newArray(prefix string) {
	var count int

	pj.err = pj.output.WriteByte(_bracketOpen)
	pj.iter.ReadArrayCB(func(*jsoniter.Iterator) bool {
		elem := ""

		if count > 0 {
			elem += _coma
		}

		if !pj.comp {
			elem += _jump + prefix + _doubleSpace
		}

		_, pj.err = pj.output.WriteString(elem)
		count++

		return pj.parseElmt(prefix + _doubleSpace)
	})

	if !pj.comp && count > 0 {
		_, pj.err = pj.output.WriteString(_jump + prefix)
	}

	pj.err = pj.output.WriteByte(_bracketClose)
}

func (pj *JSON) newObject(prefix string) {
	var count int

	pj.err = pj.output.WriteByte(_bodyOpen)
	pj.iter.ReadMapCB(func(Iter *jsoniter.Iterator, field string) bool {
		elem := _amp + field + _next

		if !pj.comp {
			elem = _jump + prefix + _doubleSpace + elem + _space
		}

		if count > 0 {
			elem = _coma + elem
		}

		_, pj.err = pj.output.WriteString(elem)
		count++

		return pj.parseElmt(prefix + _doubleSpace)
	})

	if !pj.comp && count > 0 {
		_, pj.err = pj.output.WriteString(_jump + prefix)
	}

	pj.err = pj.output.WriteByte(_bodyClose)
}

func (pj *JSON) parseElmt(prefix string) bool {
	var t = pj.iter.WhatIsNext()

	switch t {
	case jsoniter.ArrayValue:
		pj.newArray(prefix)

	case jsoniter.ObjectValue:
		pj.newObject(prefix)

	case jsoniter.InvalidValue:
		return false
	default:
		pj.printSimpleValue()
	}

	return true
}

func (pj *JSON) run() {
	defer pj.output.Flush()

	for {
		if !pj.parseElmt("") {
			if pj.iter.Error != io.EOF {
				pj.err = pj.iter.Error
			}

			return
		}
	}
}
