package word

import (
	"../mora"
)

type WordList []Word
type Word struct {
	Kanji string
	Hira  string
	Moras *mora.Moras
}

func NewWord(kanji, hira, vowel string) (w *Word) {
	return &Word{kanji, hira, mora.ToMoras(vowel)}
}

func (w *Word) GetMora(len int) (*mora.Mora, int) {
	var m = (*w).Moras.Get(len)
	return m, m.Len()
}

/*
func (w *Word) FormerMora(len int) *mora.Moras {
	var ms = (*w).Moras.Former(len)
	return ms
}
*/

func NewWordList() (wl *WordList) {
	return &WordList{}
}

func (wl *WordList) Add(w *Word) {
	*wl = append(*wl, *w)
}

// ===================================
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
