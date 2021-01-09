package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	//"github.com/inahym196/GoToStation/src/stationTree"
	//"github.com/inahym196/GoToStation/src/word"
	"./src/stationTree"
)

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
	var inputData = "きょうのごはんはハンバーグ"
	//var inputData =
	var romas = word.StrToRomas(inputData)
	fmt.Printf("input: %s\n", romas)
	//var StoreData = newWordStoreList()
	//var StoreData = Search(RootTree,romas)
	var StoreData = RootTree.Search(romas)
	fmt.Printf("Search Complate.StoreData:%v\n\t", StoreData)

}
