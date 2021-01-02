package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/inahym196/gojaconv/jaconv"
)

type Word struct {
	Kanji string
	Hira  string
	Vowel string
}

type StationTree struct {
	Len          int
	CurrentVowel string
	WordList     []Word
	ChildTree    map[string]*StationTree
}

func NewStationTree(length int, vowel string) *StationTree {
	var tree = new(StationTree)
	if matched, _ := regexp.MatchString(`[yt]`, vowel); matched {
		tree.Len = length
	} else {
		tree.Len = length + 1
	}
	tree.CurrentVowel = vowel
	tree.ChildTree = map[string]*StationTree{}
	return tree
}

func (tree StationTree) addChildTree(vowel string, nextTree *StationTree) {
	tree.ChildTree[vowel] = nextTree
}

func (tree *StationTree) addWordList(word Word) {
	(*tree).WordList = append((*tree).WordList, word)
}

func (tree StationTree) getChildTree(vowel string) (*StationTree, bool) {
	if childTree, ok := tree.ChildTree[vowel]; ok {
		return childTree, true
	}
	return NewStationTree(0, ""), false
}

func OpenReadFile(filename string) io.ReadCloser {
	fp, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return fp
}

func OpenWriteFile(filename string) io.WriteCloser {
	fp, err := os.OpenFile(filename, os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return fp
}

func ExtractText(rawText string /*, gomifp io.WriteCloser*/) (string, string) {

	rawText = strings.Replace(rawText, "・", "", -1)
	regexpKanji := regexp.MustCompile(`(.+?)駅`)
	matchKanji := regexpKanji.FindStringSubmatch(rawText)
	matchKanjiIndex := regexpKanji.FindStringIndex(rawText)
	if matchKanji != nil {
		kanjiSubText := rawText[matchKanjiIndex[0]:matchKanjiIndex[1]]
		hiraSubText := rawText[matchKanjiIndex[1]:]

		if matched, _ := regexp.MatchString(`^[ぁ-んァヶ]+駅`, kanjiSubText); matched || len(hiraSubText) == 0 {
			return matchKanji[1], KanaToHira(matchKanji[1])
		}
		regexpHira := regexp.MustCompile(`（([ぁ-んァヶ]+)えき`)
		matchHira := regexpHira.FindStringSubmatch(hiraSubText)
		if matchHira != nil {
			return matchKanji[1], KanaToHira(matchHira[1])
		}
		//fmt.Fprintf(gomifp, "raw: %v\nsub: %v\tkanji: %v\n\n", rawText, hiraSubText, matchKanji[1])
	}
	//fmt.Fprintf(gomifp, "raw: %v\n", rawText)
	return "", ""
}

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

func TextToCsv(readfile, writefile string) {

	readFp := OpenReadFile(readfile)
	defer readFp.Close()
	scanner := bufio.NewScanner(readFp)

	writeFp := OpenWriteFile(writefile)
	defer writeFp.Close()

	/*
		gomiFp := OpenWriteFile("gomi.txt")
		defer gomiFp.Close()
	*/

	for scanner.Scan() {
		kanji, hira := ExtractText(scanner.Text() /*, gomiFp*/)
		if kanji != "" {
			vowel := jaconv.ToHebon(hira)
			len := len(vowel) - strings.Count(vowel, "y") - strings.Count(vowel, "t")
			fmt.Fprintf(writeFp, "%v,%v,%v,%v\n", kanji, hira, vowel, len)
		}
	}
}

func GrowTree(reader *csv.Reader, writer io.Writer) {
	var RootTree = NewStationTree(0, "")

	for i := 0; ; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if i < 0 {
			continue
		}

		word := Word{
			Kanji: string(record[0]),
			Hira:  string(record[1]),
			Vowel: string(record[2]),
		}
		wordLen, _ := strconv.Atoi(record[3])
		//fmt.Printf("\n%#v\n", word)

		var currentTree = RootTree
		for l := 0; l < wordLen; l++ {
			var currentVowel = word.Vowel[l : l+1]
			var totalVowel = word.Vowel[:l+1]
			var childTree, ok = currentTree.getChildTree(currentVowel)
			if ok {
				//fmt.Printf("%v-tree is exist. move it.\n", totalVowel)
				currentTree = childTree
			} else {
				//fmt.Printf("%v-tree is not exist. create it.\n", totalVowel)
				var nextTree = NewStationTree(l, totalVowel)
				currentTree.addChildTree(currentVowel, nextTree)
				currentTree = nextTree
			}
		}
		currentTree.addWordList(word)
	}
	jsonData, err := json.MarshalIndent(RootTree, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(writer, "%v", string(jsonData))
}

func CsvToJson(readfile, writefile string) {
	readFp := OpenReadFile(readfile)
	defer readFp.Close()
	reader := csv.NewReader(readFp)

	writeFp := OpenWriteFile(writefile)
	defer writeFp.Close()

	GrowTree(reader, writeFp)
}

func main() {

	//TextToCsv("raw_datalist.txt", "datalist.csv")
	CsvToJson("datalist.csv", "datalist.json")

}
