package roma

import (
	"fmt"
)

type Romas []string

/* ===== internal func ===== */
func newRomas() (rs *Romas) {
	return &Romas{}
}

func InitRomas(strArr []string) (rs *Romas) {
	rs = newRomas()
	*rs = strArr
	//fmt.Printf("hebon is %v\n", strArr)
	return rs
}

/* ===== public func ===== */
func (rs *Romas) Add(s string) {
	*rs = append(*rs, s)
}

func (rs *Romas) GetAt(l int) (s string) {
	return (*rs)[l]
}

func (rs *Romas) GetVowel(l int) (s string) {
	s = (*rs)[l]
	return s[len(s)-1 : len(s)]
}

func (rs *Romas) Slice(start, end int) (RS *Romas) {
	RS = newRomas()
	//fmt.Printf("rsLen:%v,end:%v\t", rs.Len(), end)
	if rs.Len() < end {
		end = rs.Len()
	}
	//fmt.Printf("then, end:%v\n", end)
	//fmt.Printf("rs:%v\n", rs)
	for ; start < end; start++ {
		//fmt.Printf("\trs[%v]=%v\n", start, (*rs).GetAt(start))
		(*RS).Add((*rs).GetAt(start))
	}
	return RS
}

func (rs *Romas) String() (str string) {
	if (*rs).Len() == 0 {
		str = "[no data]"
	}
	for _, m := range *rs {
		str += fmt.Sprintf("%v ", m)
	}
	return str
}

func (rs *Romas) Len() int {
	if rs == nil {
		return 0
	}
	return len(*rs)
}

func (rs *Romas) InsertBefore(RS *Romas) {
	for _, s := range *rs {
		(*RS).Add(s)
	}
	*rs = *RS
}
