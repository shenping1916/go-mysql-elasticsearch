// https://studygolang.com/articles/12623
package base

import (
	"fmt"
	"strings"
	"encoding/json"
)

type Smap []*SortMapNode

type SortMapNode struct {
	Key string
	Val interface{}
}

func (c *Smap) Put(key string, val interface{}) {
	index, _, ok := c.get(key)
	if ok {
		(*c)[index].Val = val
	} else {
		node := &SortMapNode{Key: key, Val: val}
		*c = append(*c, node)
	}
}

func (c *Smap) Get(key string) (interface{}, bool) {
	_, val, ok := c.get(key)
	return val, ok
}

func (c *Smap) get(key string) (int, interface{}, bool) {
	for index, node := range *c {
		if node.Key == key {
			return index, node.Val, true
		}
	}
	return -1, nil, false
}

func ToSortedMapJson(smap *Smap) string {
	s := "{"
	for _, node := range *smap {
		v := node.Val
		isSamp := false
		str := ""
		switch v.(type){
		case *Smap:
			isSamp = true
			str = ToSortedMapJson(v.(*Smap))
		}

		if(!isSamp){
			b, _ := json.Marshal(node.Val)
			str = string(b)
		}

		s = fmt.Sprintf("%s\"%s\":%s,", s, node.Key, str)
	}
	s = strings.TrimRight(s, ",")
	s = fmt.Sprintf("%s}", s)
	return s
}
