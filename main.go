package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"unicode"

	//"github.com/inahym196/FormatStationData/src/stationTree"
	"./src/stationTree"
	"./src/word"
	"github.com/inahym196/FormatStationData/src/roma"
	"github.com/inahym196/gojaconv/jaconv"
)

type wordStore struct {
	Word word.Word
	Len  int
}

func (ws *wordStore) len() int {
	return len(*ws)
}

func (ws *wordStore) push(w *word) {
	WS = wordStore{Word: *w, Len: *w.Len()}
	*ws = append(*ws, WS)
}

func (ws *wordStore) pop() (w *word, ok bool) {
	if *ws.Len() == 0 {
		return nil, false
	}
	var last = *ws.Len() - 1
	*w = (*ws)[last]
	*ws = (*ws)[:last]
	return w, true
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

func KanaToHira(str string) string {
	codeDiff := 0x30a1 - 0x3041
	src := []rune(str)
	dst := make([]rune, len(src))
	for i, r := range src {
		switch {
		case unicode.In(r, unicode.Katakana):
			dst[i] = r - rune(codeDiff)
		default:
			dst[i] = r
		}
	}
	return string(dst)
}

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

func strToRomas(s string) *roma.Romas {
	return roma.InitRomas(jaconv.ToHebon(KanaToHira(s)))
}

func main() {

	const stationLenMax int = 25

	var RootTree *stationTree.StationTree = ReadJson("src/datalist.json")
	var inputData = "オバマ"
	var romas = strToRomas(inputData)
	fmt.Printf("input: %s\n", romas)
	var searchStr string
	var StoreData *[]wordStore
	var latterRomas = romas
	if romas.Len() > stationLenMax {
		var subRomas = latterRomas.Slice(0, stationLenMax)
	}
	var sStart, sEnd = 0, stationLenMax
	for subRomas.Len() > 0 {
		var wordList = RootTree.SearchLeafWordList(romas, stationLenMax)
		if wordList.Len() != 0 {
			word := wordList.Eval()
			StoreData.push(word)
		} else if popWord, ok := StoreData.pop; ok {
			var tmpRomas = roma.Romas{popWord.Romas}
			*romas = append(tmp, romas)
		}
	}

	//var LeafWordList, matched = currentTree.GetLeafWordList(inputRomas)
	fmt.Printf("%v\n%v\n", (*wordList)[0], depth)

	/*
		var popdata StoredStation
		var StoredStationList StoredStationList
		tmp := "まんがよむならぶっくらいぶ"
		tmp = jaconv.ToHebon(tmp)
		print(tmp)
		itr_start, itr_end := 0, len(tmp)

		for i := 5; len(tmp) != itr_start && i > 0; i-- {
			matchLen, matchWord, err := WordSearch(tmp[itr_start:itr_end], jsonData)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			if matchLen != 0 {
				PushStationList(&StoredStationList, matchLen, matchWord)
				itr_start += matchLen
			} else if len(StoredStationList) > 0 {
				popdata, err = PopStationList(&StoredStationList)
				if err != nil {
					fmt.Println(err)
				}
				itr_start -= popdata.Len
			} else {
				break
			}
			fmt.Printf("%v\n", StoredStationList)
		}
	*/

}
