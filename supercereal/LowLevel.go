package supercereal

import (
	"fmt"
)

type JSONStream struct {
	cb func(json []byte)

	front  int
	buffer []byte

	currentScopeIsEmpty bool
	lastWasKey          bool
}

func NewJSONStream() *JSONStream {
	js := &JSONStream{buffer: make([]byte, 1024*1024)} // 1mb buffer by default
	js.Reset()
	return js
}

func (p *JSONStream) OnJSON(cb func(json []byte)) {
	p.cb = cb
}

func (p *JSONStream) OpenArray() {
	p.currentScopeIsEmpty = true
	p.buffer[p.front] = '['
	p.front++
}

// Root or Value depends on if lastWasKey
func (p *JSONStream) OpenObject() {
	// put comma if this scope is not empty and last was not key
	if !p.currentScopeIsEmpty && !p.lastWasKey {
		p.buffer[p.front] = ','
		p.front++
	}

	p.currentScopeIsEmpty = true
	p.buffer[p.front] = '{'
	p.front++
}

func (p *JSONStream) CloseArray() {
	p.buffer[p.front] = ']'
	p.front++

	p.lastWasKey = false
	p.currentScopeIsEmpty = false
}

func (p *JSONStream) CloseObject() {
	p.buffer[p.front] = '}'
	p.front++

	p.lastWasKey = false
	p.currentScopeIsEmpty = false
}

// Key only
func (p *JSONStream) PutKey(key []byte) {

	p.lastWasKey = true

	// put comma if this scope is not empty!
	if !p.currentScopeIsEmpty {
		p.buffer[p.front] = ','
		p.front++
	}
	p.currentScopeIsEmpty = false

	p.buffer[p.front] = '"'
	p.front++
	copy(p.buffer[p.front:], key)
	p.front += len(key)
	p.buffer[p.front] = '"'
	p.front++
	p.buffer[p.front] = ':'
	p.front++
}

// Value only
func (p *JSONStream) PutInt(value int) {
	p.lastWasKey = false

	byteRep := []byte(fmt.Sprintf("%d", value))
	copy(p.buffer[p.front:], byteRep)
	p.front += len(byteRep)
}

// Value only
func (p *JSONStream) PutNull() {
	p.lastWasKey = false

	byteRep := []byte("null")
	copy(p.buffer[p.front:], byteRep)
	p.front += len(byteRep)
}

// Value only
func (p *JSONStream) PutBoolean(value bool) {
	p.lastWasKey = false

	byteRep := "false"
	if value {
		byteRep = "true"
	}
	copy(p.buffer[p.front:], byteRep)
	p.front += len(byteRep)
}

// Value only
func (p *JSONStream) PutString(value []byte) {
	p.lastWasKey = false

	p.buffer[p.front] = '"'
	p.front++
	copy(p.buffer[p.front:], value)
	p.front += len(value)
	p.buffer[p.front] = '"'
	p.front++
}

func (p *JSONStream) Reset() {
	p.front = 0
	p.currentScopeIsEmpty = true
	p.lastWasKey = false
	//p.OpenObject()
}

func (p *JSONStream) End() {
	//p.CloseObject()
	p.cb(p.buffer[:p.front])
}