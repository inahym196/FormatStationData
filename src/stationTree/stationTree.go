package stationTree

import (
	"fmt"

	"github.com/inahym196/FormatStationData/src/roma"
	"github.com/inahym196/FormatStationData/src/word"
	"github.com/inahym196/gojaconv/jaconv"
)

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

func (tree *StationTree) MakeLeaf(romas *roma.Romas) (leaf *StationTree) {
	var currentTree = tree
	var totalVowel = ""
	for i, romasLen := 0, romas.Len(); i < romasLen; i++ {
		var currentVowel = romas.GetVowel(i)
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
	var leaf = (*tree).MakeLeaf(romas)
	leaf.addWordList(word)

}

func (tree *StationTree) SearchLeaf(romas *roma.Romas) (leaf *StationTree, reached bool) {
	//var totalVowel = ""
	for i, romasLen := 0, romas.Len(); i < romasLen; i++ {
		var currentVowel = romas.GetVowel(i)
		if childTree, ok := tree.getChildTree(currentVowel); ok {
			tree = childTree
		} else {
			fmt.Printf("childTree is not exist. \n")
			return tree, false
		}
	}
	//tree.Debug()
	return tree, true
}
