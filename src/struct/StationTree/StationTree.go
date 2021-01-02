package StationTree

import (
	"regexp"
	"strconv"

	"../Word"
)

type StationTree struct {
	Len          int
	CurrentVowel string
	WordList     *[]Word.Word
	ChildTree    map[string]*StationTree
}

func NewStationTree(length int, vowel string) *StationTree {
	var tree = new(StationTree)
	if matched, _ := regexp.MatchString(`^[yt]*$`, vowel); matched {
		//println("word: ", vowel, " is 'y,t' or empty")
		tree.Len = length
	} else {
		tree.Len = length + 1
	}
	tree.CurrentVowel = vowel
	tree.ChildTree = map[string]*StationTree{}
	return tree
}

func (tree *StationTree) addChildTree(vowel string, nextTree *StationTree) {
	(*tree).ChildTree[vowel] = nextTree
}

func (tree *StationTree) addWordList(word Word.Word) {
	var WordList *[]Word.Word = new([]Word.Word)
	*WordList = append(*WordList, word)
	(*tree).WordList = WordList
}

func (tree *StationTree) getChildTree(vowel string) (*StationTree, bool) {
	if childTree, ok := (*tree).ChildTree[vowel]; ok {
		return childTree, true
	}
	return nil, false
}

func (RootTree *StationTree) GrowTree(record []string) {
	word := Word.Word{
		Kanji: string(record[0]),
		Hira:  string(record[1]),
		Vowel: string(record[2]),
	}
	wordLen, _ := strconv.Atoi(record[3])
	var currentTree = RootTree
	for l := 0; l < wordLen; l++ {
		var currentVowel = word.Vowel[l : l+1]
		var totalVowel = word.Vowel[:l+1]
		var childTree, ok = currentTree.getChildTree(currentVowel)
		if ok {
			//fmt.Printf("%v-tree is exist. move it. last word is %v.\n", totalVowel, currentVowel)
			currentTree = childTree
		} else {
			//fmt.Printf("%v-tree is not exist. create and move it. last word is %v.\n", totalVowel, currentVowel)
			var nextTree = NewStationTree(l, totalVowel)
			currentTree.addChildTree(currentVowel, nextTree)
			currentTree = nextTree
		}
	}
	currentTree.addWordList(word)
}
