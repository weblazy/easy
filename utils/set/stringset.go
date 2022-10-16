package set

import (
	"strings"
)

type StringSet struct {
	store map[string]struct{}
}

func NewStringSet() *StringSet {
	return &StringSet{store: make(map[string]struct{})}
}

func (ss *StringSet) Add(val string) {
	ss.store[val] = struct{}{}
}

func (ss *StringSet) BatchAdd(vals ...string) {
	for _, vv := range vals {
		ss.Add(vv)
	}
}

func (ss *StringSet) Has(val string) bool {
	_, ok := ss.store[val]
	return ok
}

func (ss *StringSet) Delete(val string) {
	delete(ss.store, val)
}

func (ss *StringSet) Size() int {
	return len(ss.store)
}

func (ss *StringSet) ToArray() []string {
	arr := make([]string, 0, len(ss.store))

	for v := range ss.store {
		arr = append(arr, v)
	}

	return arr
}

func (ss *StringSet) String() string {
	return strings.Join(ss.ToArray(), ",")
}
