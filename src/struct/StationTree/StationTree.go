package StationTree

import (
	"fmt"
	"strconv"

	"../mora"
	"../word"
)

type StationTree struct {
	//Len          int
	//CurrentVowel string
	Moras     mora.Moras
	WordList  *word.WordList
	ChildTree map[mora.Mora]*StationTree
}

func NewStationTree(ms *mora.Moras) *StationTree {
	var tree = new(StationTree)
	if ms == nil {
		ms = mora.NewMoras()
	}
	tree.Moras = *ms
	tree.WordList = word.NewWordList()
	tree.ChildTree = map[mora.Mora]*StationTree{}
	return tree
}

func (tree *StationTree) addChildTree(mora *mora.Mora, nextTree *StationTree) {
	(*tree).ChildTree[*mora] = nextTree
}

func (tree *StationTree) addWordList(w *word.Word) {
	(*tree).WordList.Add(w)
}

func (tree *StationTree) getChildTree(mora *mora.Mora) (*StationTree, bool) {
	if childTree, ok := (*tree).ChildTree[*mora]; ok {
		return childTree, true
	}
	return nil, false
}

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

func (RootTree *StationTree) GrowTree(record []string) {

	var word = word.NewWord(record[0], record[1])
	var moras = mora.ToMoras(record[2])
	var wordLen, _ = strconv.Atoi(record[3])
	var currentTree = RootTree
	for l := 0; l < wordLen; l++ {
		var currentMora = (*moras).Get(l)
		var moraLen = currentMora.Len()
		l += moraLen - 1
		var currentMoras = (*moras)[:l+1]
		if childTree, ok := currentTree.getChildTree(currentMora); ok {
			//fmt.Printf("current mora is %v. %v-tree is exist. move it. \n", *currentMora, currentMoras.String())
			currentTree = childTree
		} else {
			//fmt.Printf("current mora is %v. %v-tree is not exist. create and move it. \n", *currentMora, currentMoras.String())
			var nextTree = NewStationTree(&currentMoras)
			currentTree.addChildTree(currentMora, nextTree)
			currentTree = nextTree
		}
	}
	currentTree.addWordList(word)

}
