package stationTree

import (
	"GoToStation/src/roma"
	"GoToStation/src/word"
	"fmt"

	"github.com/inahym196/gojaconv/jaconv"
)

func Search(s string) (ws *wordStore.wordStore) {
	var romas = word.StrToRomas(s)
	return &wordStore.wordStore{}
}

/*
func SearchStation(tree *stationTree.StationTree) (StoreData *wordStore) {

	var latterRomas, subRomas = romas, romas
	var sStart, sCount = 0, stationLenMax

	for latterRomas.Len() > 0 {
		if latterRomas.Len() > stationLenMax {
			subRomas = latterRomas.Slice(0, stationLenMax)
		} else if latterRomas.GetAt(0) == "-" {
			subRomas = latterRomas.Slice(1, stationLenMax)
		} else {
			subRomas = latterRomas
		}
		fmt.Printf("Search:[%v],Count:[%v]\n==========\n", subRomas, sCount)
		var wordList = tree.SearchLeafWordList(subRomas, sCount)
		if wordList != nil {
			sCount = stationLenMax
			word := wordList.Eval()
			fmt.Printf("matched. select word: %v\n", word)
			wordLen := word.Len()
			StoreData.push(word)
			fmt.Printf("StoreData: %v\t", StoreData)
			fmt.Printf("sCount:%v\n", sCount)
			fmt.Printf("latterRomas: %v [%v:%v]=>", latterRomas, wordLen, wordLen+sCount)
			latterRomas = latterRomas.Slice(wordLen, wordLen+sCount)
			fmt.Printf("%v\n", latterRomas)
			sStart += wordLen
		} else if popWord, wordLen := (*StoreData).pop(); wordLen != 0 {
			fmt.Printf("Popped.\nlatterRomas: %v =>", latterRomas)
			latterRomas.InsertBefore(popWord.GetRomas())
			fmt.Printf("%v\n", latterRomas)
			sStart -= wordLen
			sCount = wordLen - 1
		} else {
			//fmt.Printf("word[%v] couldn't find even thought popped it.", latterRomas)
			break
		}
	}
}
*/

type StationTree struct {
	Vowel     string                  `json:"Vowel"`
	WordList  *word.WordList          `json:"WordList"`
	ChildTree map[string]*StationTree `json:"ChildTree"`
}

/* ===== internal func ===== */
func (tree *StationTree) addChildTree(s string, nextTree *StationTree) {
	(*tree).ChildTree[s] = nextTree
}

func (tree *StationTree) addWordList(w *word.Word) {
	(*tree).WordList.Add(w)
}

func (tree *StationTree) GetChildTree(s string) (*StationTree, bool) {
	if childTree, ok := (*tree).ChildTree[s]; ok {
		return childTree, true
	}
	return nil, false
}

func (tree *StationTree) getWordList() (wl *word.WordList, l int) {
	var WordListLen = (*tree).WordList.Len()
	if WordListLen > 0 {
		return (*tree).WordList, WordListLen
	}
	return nil, 0
}

func (tree *StationTree) getDepth() int {
	return len((*tree).Vowel)
}

func NewStationTree(v string) *StationTree {
	var tree = new(StationTree)
	tree.Vowel = v
	tree.WordList = word.NewWordList()
	tree.ChildTree = map[string]*StationTree{}
	return tree
}

/* ===== public func ===== */
func (tree *StationTree) String() (str string) {

	str = "&tree::[ChildTree: ["
	for i := range (*tree).ChildTree {
		str += fmt.Sprintf("%v,", i)
	}
	str += "],"

	var WordLen = tree.WordList.Len()
	if firstWord := tree.WordList.FirstWord(); firstWord != nil {
		str += fmt.Sprintf("WordList[%v] = %v,...", WordLen, firstWord)
	} else {
		str += fmt.Sprintf("FirstWord = <nil>")
	}
	return str
}

func (tree *StationTree) makeLeaf(romas *roma.Romas) (leaf *StationTree) {
	var currentTree = tree
	var totalVowel = ""
	for i, romasLen := 0, romas.Len(); i < romasLen; i++ {
		var currentVowel = romas.GetVowel(i)
		totalVowel += currentVowel
		//var currentRomas = romas.Slice(0, l)
		if childTree, ok := currentTree.GetChildTree(currentVowel); ok {
			//fmt.Printf("currentRoma=%v,%v-tree !=nil. move it.vowel=%v \n", currentRoma, currentRomas, currentVowel)
			currentTree = childTree
		} else {
			//fmt.Printf("currentRoma=%v,%v-tree ==nil. create and move it. vowel=%v \n", currentRoma, currentRomas, currentVowel)
			var nextTree = NewStationTree(totalVowel)
			currentTree.addChildTree(currentVowel, nextTree)
			currentTree = nextTree
		}
	}
	return currentTree
}

func (tree *StationTree) GrowTree(record []string) {

	var word = word.NewWord(record[0], record[1])
	var romas = roma.InitRomas(jaconv.ToHebon(record[1]))
	var leaf = (*tree).makeLeaf(romas)
	leaf.addWordList(word)

}

func (tree *StationTree) SearchLeafWordList(romas *roma.Romas, limit int) (leafWordList *word.WordList) {

	if romas == nil {
		return nil
	} else if limit > romas.Len() || 0 > limit {
		limit = romas.Len()
	}
	for i := 0; i < limit; i++ {
		var searchVowel = romas.GetVowel(i)
		//fmt.Printf("currentTree: %v\tcurrentVowel: %v\n", tree, searchVowel)
		if childTree, ok := tree.GetChildTree(searchVowel); ok {
			tree = childTree
			if wordList, wordListLen := tree.getWordList(); wordListLen > 0 {
				leafWordList = wordList
			}
		} else {
			//fmt.Printf("SearchLeafWordListFunction::childTree is not exist. \n")
			return leafWordList
		}
	}
	return leafWordList
}
