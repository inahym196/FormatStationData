package word

import "github.com/inahym196/FormatStationData/src/roma"

type WordList []Word
type Word struct {
	Kanji string     `json:"Kanji"`
	Hira  string     `json:"Hira"`
	Romas roma.Romas `json:"Romas"`
}

/* ===== internal func ===== */

/* ===== public func ===== */
func NewWord(kanji, hira string, rs roma.Romas) (w *Word) {
	return &Word{kanji, hira, rs}
}

func (w *Word) Len() int {
	return (*w).Romas.Len()
}

func NewWordList() (wl *WordList) {
	return &WordList{}
}

func (wl *WordList) Add(w *Word) {
	*wl = append(*wl, *w)
}

func (wl *WordList) Eval() *Word {
	return &(*wl)[0]
}

func (w *WordList) Len() int {
	if w != nil {
		return len(*w)
	}
	return 0
}

func (w *WordList) FirstWord() (firstWord *Word) {
	if w.Len() != 0 {
		firstWord = &Word{}
		*firstWord = (*w)[0]
		return firstWord
	}
	return nil
}
