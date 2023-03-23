package mapx

import (
	"container/list"
)

type Obj struct {
	key string
	val interface{}
}
type OrderMap struct {
	l *list.List
	m map[string]*list.Element
}

func (m *OrderMap) Put(key string, val interface{}) {
	if e, ok := m.m[key]; ok {
		e.Value.(*Obj).val = val
		return
	}
	e := m.l.PushBack(&Obj{key: key, val: val})
	m.m[key] = e
}

func (m *OrderMap) Get(key string) interface{} {
	if e, ok := m.m[key]; ok {
		return e.Value.(*Obj).val
	}

	return nil
}

func (m *OrderMap) Delete(key string) {
	if e, ok := m.m[key]; ok {
		m.l.Remove(e)
		delete(m.m, key)
	}
}

func (m *OrderMap) Range(f func(key string, val interface{})) {
	e := m.l.Front()
	for e != nil {
		obj := e.Value.(*Obj)
		f(obj.key, obj.val)
		e = e.Next()
	}
}
