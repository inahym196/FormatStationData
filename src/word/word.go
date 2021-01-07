package word

type WordList []Word
type Word struct {
	Kanji string
	Hira  string
}

/* ===== internal func ===== */

/* ===== public func ===== */
func NewWord(kanji, hira string) (w *Word) {
	return &Word{kanji, hira}
}

func NewWordList() (wl *WordList) {
	return &WordList{}
}

func (wl *WordList) Add(w *Word) {
	*wl = append(*wl, *w)
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
