package word

import (
	"unicode"

	"github.com/inahym196/FormatStationData/src/roma"
	"github.com/inahym196/gojaconv/jaconv"
)

type WordList []Word
type Word struct {
	Kanji string     `json:"Kanji"`
	Hira  string     `json:"Hira"`
	Romas roma.Romas `json:"Romas"`
}

/* ===== internal func ===== */

/* ===== public func ===== */
func KanaToHira(str string) string {
	codeDiff := 0x30a1 - 0x3041
	src := []rune(str)
	dst := make([]rune, len(src))
	for i, r := range src {
		switch {
		case unicode.In(r, unicode.Katakana):
			dst[i] = r - rune(codeDiff)
		default:
			dst[i] = r
		}
	}
	return string(dst)
}

func StrToRomas(s string) *roma.Romas {
	return roma.InitRomas(jaconv.ToHebon(KanaToHira(s)))
}

func NewWord(kanji, hira string) (w *Word) {
	var romas = StrToRomas(hira)
	return &Word{kanji, hira, *romas}
}

func (w *Word) Len() int {
	return (*w).Romas.Len()
}

func (w *Word) GetRomas() *roma.Romas {
	return &(*w).Romas
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
