package roma

import (
	"fmt"
)

type Romas []string

/* ===== internal func ===== */
func InitRomas(strArr []string) (ms *Romas) {
	ms = &Romas{}
	*ms = strArr
	//fmt.Printf("hebon is %v\n", strArr)
	return ms
}

/* ===== public func ===== */
func (ms *Romas) Add(s string) {
	*ms = append(*ms, s)
}

func (ms *Romas) Get(l int) (s string) {
	return (*ms)[l]
}

func (ms *Romas) Slice(start, end int) (s string) {
	for ; start < end; start++ {
		s += (*ms).Get(start)
	}
	return s
}

func (ms *Romas) String() (str string) {
	for _, m := range *ms {
		str += fmt.Sprintf("%v", m)
	}
	return str
}

func (ms *Romas) Len() int {
	return len(*ms)
}
