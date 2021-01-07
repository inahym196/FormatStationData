package stationTree

import (
	"fmt"

	"../roma"
	"../word"
	"github.com/inahym196/gojaconv/jaconv"
)

type StationTree struct {
	Vowel     string
	WordList  *word.WordList
	ChildTree map[string]*StationTree
}

/* ===== internal func ===== */
func (tree *StationTree) addChildTree(s string, nextTree *StationTree) {
	(*tree).ChildTree[s] = nextTree
}

func (tree *StationTree) addWordList(w *word.Word) {
	(*tree).WordList.Add(w)
}

func (tree *StationTree) getChildTree(s string) (*StationTree, bool) {
	if childTree, ok := (*tree).ChildTree[s]; ok {
		return childTree, true
	}
	return nil, false
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

	str = "&StationTree.StationTree::[ChildTree: ["
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

func (tree *StationTree) MakeLeafTree(romas *roma.Romas) (leaf *StationTree) {
	var currentTree = tree
	var romasLen = romas.Len()
	var totalVowel = ""
	for l := 0; l < romasLen; l++ {
		var currentRoma = romas.Get(l)
		var currentVowel = currentRoma[len(currentRoma)-1 : len(currentRoma)]
		totalVowel += currentVowel
		//var currentRomas = romas.Slice(0, l)
		if childTree, ok := currentTree.getChildTree(currentVowel); ok {
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
	var leaf = (*tree).MakeLeafTree(romas)
	leaf.addWordList(word)

}
