package wordStore

import (
	"GoToStation/src/word"
)

type wordStore struct {
	Word word.Word
	Len  int
}

type wordStoreList []wordStore

func (ws *wordStore) GetWord() (w *word.Word) {
	return &(*ws).Word
}

func newWordStoreList() (wsl *wordStoreList) {
	return &wordStoreList{}
}

func (wsl *wordStoreList) len() int {
	return len(*wsl)
}

func (wsl *wordStoreList) push(w *word.Word) {
	var WS = wordStore{Word: *w, Len: (*w).Len()}
	*wsl = append(*wsl, WS)
}

func (wsl *wordStoreList) pop() (w *word.Word, wordLen int) {
	if (*wsl).len() == 0 {
		return nil, 0
	}
	var last = (*wsl).len() - 1
	w = (*wsl)[last].GetWord()
	*wsl = (*wsl)[:last]
	return w, (*w).Len()
}
