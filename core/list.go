package core

import (
	"sort"
)

type UniqueList []interface{}
type StringStack struct {
	List         []string
	MaxSize      int
	IgnoreRepeat bool
}

func (ls *UniqueList) Add(val interface{}) {
	if ls == nil {
		panic("UniqueList object is nil")
	}

	for _, it := range *ls {
		if it == val {
			return
		}
	}
	*ls = append(*ls, val)
}

func (s *StringStack) Push(val string) {
	if s == nil {
		panic("UniqueList object is nil")
	}
	if s.List == nil {
		s.List = []string{}
	}
	//set default max size
	if s.MaxSize == 0 {
		s.MaxSize = 99999999
	}
	//ignore repeat
	if s.IgnoreRepeat {
		for _, it := range s.List {
			if it == val {
				return
			}
		}
	}

	l := len(s.List)
	if l < s.MaxSize {
		s.List = append([]string{val}, s.List...)
	} else {
		s.List = append([]string{val}, s.List[0:l-1]...)
	}

}

func (s *StringStack) Each(fn func(idx int, it string)) {
	for i, item := range s.List {
		fn(i, item)
	}
}

func (s *StringStack) EachByPage(idx int, size int, fn func(idx int, it string)) {
	if idx <= 0 {
		panic("idx must be an integer greater than 0")
	}
	end := idx * size
	start := (idx - 1) * size
	if end > len(s.List) {
		end = len(s.List)
	}
	if start > end {
		start = end
	}
	for i, item := range s.List[start:end] {
		fn(i, item)
	}
}

func (s *StringStack) ToBytes() []byte {
	return toBytes(s)
}

func (b Bytes) ToStack() *StringStack {
	var rs = &StringStack{}
	fromBytes(b, rs)
	return rs
}

func (ls *UniqueList) ToBytes() []byte {
	return toBytes(ls)
}

func (b Bytes) ToList() *UniqueList {
	var rs = &UniqueList{}
	fromBytes(b, rs)
	return rs
}

type PairMap struct {
	pairList  []string
	pairData  map[string]interface{}
	sortRules func(i, j interface{}) bool
}

func NewPairMap(sortRules func(i, j interface{}) bool) *PairMap {
	return &PairMap{
		pairList:  make([]string, 0),
		pairData:  make(map[string]interface{}),
		sortRules: sortRules,
	}
}

func (p *PairMap) Swap(i, j int) {
	p.pairList[i], p.pairList[j] = p.pairList[j], p.pairList[i]
}

func (p *PairMap) Len() int {
	return len(p.pairList)
}

func (p *PairMap) Less(i, j int) bool {
	key1 := p.pairList[i]
	key2 := p.pairList[j]
	return p.sortRules(p.pairData[key1], p.pairData[key2])
}

func (p *PairMap) Push(key string, data interface{}) {
	if p.Exsit(key) {
		p.Remove(key)
	}
	p.pairList = append(p.pairList, key)
	p.pairData[key] = data
}

func (p *PairMap) Remove(key string) {
	if !p.Exsit(key) {
		return
	}
	delete(p.pairData, key)
	for i, inKey := range p.pairList {
		if inKey == key {
			p.pairList = append(p.pairList[0:i], p.pairList[i+1:]...)
			return
		}
	}
}

func (p *PairMap) Sort() {
	sort.Sort(p)
}

func (p *PairMap) Exsit(key string) bool {
	_, ok := p.pairData[key]
	return ok
}

func (p *PairMap) GetLast(idx, size int) []interface{} {
	limit := idx * size
	var rs []interface{}
	for i, key := range p.pairList {
		if i >= limit && i < limit+size {
			data := p.pairData[key]
			rs = append(rs, data)
		}
		if i >= limit+size {
			break
		}
	}
	return rs
}

func (p *PairMap) Get(key string) interface{} {
	if !p.Exsit(key) {
		return nil
	}
	return p.pairData[key]
}
