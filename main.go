package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/inahym196/FormatStationData/src/stationTree"
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

/*
	func sjis_to_utf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
	}

	func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world\n")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "error\n")
	}
	tmp, _ := sjis_to_utf8(r.Form.Get("text"))
	tmp = jaconv.ToHebon(tmp)
	}

	func HttpFunc() {
	port, _ := strconv.Atoi(os.Args[1])
	fmt.Printf("Starting server at Port %d", port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}

	func WordSearch(str string, jsonData *[]StationData) (int, Word.Word, error) {
	for len(str) > 1 {
		for _, obj := range (*jsonData)[len(str)-1].Word {
			if strings.EqualFold(obj.Vowel, str) {
				return len(str), obj, nil
			}
		}
		str = str[:len(str)-1]
	}
	return 0, Word.Word{}, errors.New("target is not match")
	}

*/

func ReadJson(filename string) (jsonData *stationTree.StationTree) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bytes, &jsonData); err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func main() {

	const stationLenMax int = 25

	var RootTree = ReadJson("src/datalist.json")
	//var child, _ = RootTree.GetChildTree("u")
	//fmt.Printf("roottree:%v\n", child)
	var inputData = "きょうのごはんはハンバーグ"
	var romas = word.StrToRomas(inputData)
	fmt.Printf("input: %s\n", romas)
	//var searchStr string
	var StoreData = newWordStoreList()
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
		var wordList = RootTree.SearchLeafWordList(subRomas, sCount)
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
	fmt.Printf("Search Complate.StoreData:%v\n\t", StoreData)

}
