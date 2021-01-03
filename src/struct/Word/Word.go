package Word

import (
	"regexp"
)

type Word struct {
	Kanji string
	Hira  string
	Vowel string
}

type WordList []Word

func (w *Word) GetSyllable(l int) (syllable string, length int) {
	var latterWord = (*w).Vowel[l:]
	//fmt.Printf("Word is %v.\nlatter Word is %v.", w, latterWord)
	if matched, _ := regexp.MatchString(`^ty`, latterWord); matched {
		syllable = latterWord[:3]
	} else if matched, _ = regexp.MatchString(`^[ty]`, latterWord); matched {
		syllable = latterWord[:2]
	} else {
		syllable = latterWord[:1]
	}
	//fmt.Printf("syllable is %v.\nLen is %v.\n", syllable, len(syllable))
	return syllable, len(syllable)
}
