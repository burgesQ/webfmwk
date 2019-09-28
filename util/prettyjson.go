package util

import (
    "io"
	"sync"
    "bufio"
    "github.com/json-iterator/go"
)

type PrettyJson struct {
    input  io.Reader
    output *bufio.Writer
    iter   *jsoniter.Iterator
    wg     sync.WaitGroup
    err    error
    comp   bool
}

func NewPrettyJson(in io.Reader, out io.Writer) *PrettyJson {

    return &PrettyJson{
        in,
        bufio.NewWriterSize(out, 1024),
        jsoniter.Parse(jsoniter.ConfigFastest, in, 512*1024),
        sync.WaitGroup{},
        nil,
        false,
    }
}

func (p *PrettyJson) SetCompactMode() *PrettyJson {
    p.comp = true
    return p
}

func (p *PrettyJson) printSimpleValue() {

    bts := p.iter.SkipAndReturnBytes()
    p.output.Write(bts)
}

func (p *PrettyJson) parseElmt(prefix string) bool {
    
    t := p.iter.WhatIsNext()

    switch t {

    case jsoniter.ArrayValue:

        p.output.WriteByte('[')
        count := 0
        p.iter.ReadArrayCB(func(*jsoniter.Iterator) bool {
            if count > 0 {
                p.output.WriteByte(',')
            }
            if !p.comp {
                p.output.WriteByte('\n')
                p.output.WriteString(prefix)
                p.output.WriteString("  ")
            }

            count++
            return p.parseElmt(prefix + "  ")
        })

        if !p.comp && count > 0 {
            p.output.WriteByte('\n')
            p.output.WriteString(prefix)
        }

        p.output.WriteByte(']')

    case jsoniter.ObjectValue:

        p.output.WriteByte('{')
        count := 0
        p.iter.ReadMapCB(func(Iter *jsoniter.Iterator, field string) bool {
            if count > 0 {
                p.output.WriteByte(',')
            }

            if !p.comp {
                p.output.WriteByte('\n')
                p.output.WriteString(prefix)
                p.output.WriteString("  ")
            }

            p.output.WriteByte('"')
            p.output.WriteString(field)
            p.output.WriteString("\":")

            if !p.comp {
                p.output.WriteByte(' ')
            }

            count++
            return p.parseElmt(prefix + "  ")
        })

        if !p.comp && count > 0 {
            p.output.WriteByte('\n')
            p.output.WriteString(prefix)
        }

        p.output.WriteByte('}')

    case jsoniter.InvalidValue:
        return false
    default:
        p.printSimpleValue()
    }

    return true
}

func (p *PrettyJson) run() {

    defer p.output.Flush()

    for  {
        if !p.parseElmt("") {

            if p.iter.Error != io.EOF {
                p.err = p.iter.Error
            }

            return
        }
    }
}

func (p *PrettyJson) Start() {
    p.wg.Add(1)
    go func () {
        defer p.wg.Done()
        p.run()

        closer, ok := p.input.(io.Closer)
        if ok { closer.Close() }
    }()
}

func (p *PrettyJson) Close() error {
    p.wg.Wait()
    return p.err
}
