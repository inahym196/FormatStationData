package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/inahym196/FormatStationData/src/roma"
	"github.com/inahym196/gojaconv/jaconv"
)

/*

	func PushStationList(store *StoredStationList, len int, Word Word.Word) {
	*store = append(*store, *NewStation(len, Word))
	}

	func PopStationList(store *StoredStationList) (StoredStation, error) {
	if len(*store) == 0 {
		return *NewStation(0, Word.Word{}), errors.New("StoredStationList is empty.\n")
	}
	popdata := (*store)[len(*store)-1]
	*store = (*store)[:len(*store)-1]
	return popdata, nil
	}

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

func ReadJson(filename string) (jsonData *stationTree.stationTree) {
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

	var RootTree *stationTree.stationTree = ReadJson("src/datalist.json")
	RootTree.Debug()
	var inputData = "ちょんだ"
	var inputRomas = roma.InitRomas(jaconv.ToHebon(inputData))
	var currentTree, reached = RootTree.SearchLeaf(inputRomas)
	//var LeafWordList, matched = currentTree.GetLeafWordList(inputRomas)
	fmt.Printf("%v\n%v\n", currentTree, reached)

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
