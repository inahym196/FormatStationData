package mora

import (
	"fmt"
	"regexp"
)

type Mora string
type Moras []Mora

/* ===== internal func ===== */
func toMora(s string) *Mora {
	var m = Mora(s)
	return &m
}

func (ms *Moras) Add(m *Mora) {
	*ms = append(*ms, *m)
}

func NewMoras() (ms *Moras) {
	return &Moras{}
}

/* ===== public func ===== */
func (m *Mora) Len() int {
	return len(*m)
}

func (ms *Moras) Get(l int) (m *Mora) {
	return &(*ms)[l]
}

func (ms Moras) String() (str string) {
	for _, m := range ms {
		str += fmt.Sprintf("%v", m)
	}
	return str
}

func ToMoras(s string) *Moras {
	var m *Mora
	var ms = NewMoras()
	//fmt.Printf("str is %v.\n", s)
	for i, l := 0, len(s); i < l; i++ {
		//fmt.Printf("i is %v. l is %v.\n", i, l)
		if matched, _ := regexp.MatchString(`^ty`, s); matched {
			m = toMora(s[i : i+3])
		} else if matched, _ = regexp.MatchString(`^[ty]`, s); matched {
			m = toMora(s[i : i+2])
		} else {
			//fmt.Printf("else. s[%v:%v]=%v\n", i, i+1, s[i:i+1])
			m = toMora(s[i : i+1])
		}
		ms.Add(m)
	}
	return ms
}
