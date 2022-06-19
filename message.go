package cfg

import (
	"reflect"
)

type message struct {
	reflect.Type
}

func newMessage(msg interface{}) *message {
	m := &message{
		Type: reflect.ValueOf(msg).Type(),
	}
	return m
}

// new は reflect.New のラッパー
//
// reflect.New は型のインスタンスを生成後にそのポインタを返すが、
// インスタンス自体では無いことに注意する
func (m *message) new() interface{} {
	return reflect.New(m.Type).Interface()
}
