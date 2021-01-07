package roma

import (
	"fmt"

	//"github.com/inahym196/FormatStationData/src/word"
	"../src/word"
)

type Romas []string

/* ===== internal func ===== */
func InitRomas(strArr []string) (rs *Romas) {
	rs = &Romas{}
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
	for ; start < end; start++ {
		(*RS).Add((*rs).GetAt(start))
	}
	return RS
}

func (rs *Romas) String() (str string) {
	for _, m := range *rs {
		str += fmt.Sprintf("%v", m)
	}
	return str
}

func (rs *Romas) Len() int {
	return len(*rs)
}

func (rs *Romas) InsertBefore(w *word.Word) {
	romas = (*w).Romas
	*rs = append(romas, *rs)
}
