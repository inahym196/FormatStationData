package main

import (
	"bufio"
	"encoding/csv"
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

func OpenReadFile(filename string) io.ReadCloser {
	fp, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return fp
}

func OpenWriteFile(filename string) io.WriteCloser {
	fp, err := os.OpenFile(filename, os.O_WRONLY, 0666)
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
	regexpHiraStation := regexp.MustCompile(`[ぁ-んァヶ]+駅`)
	if matchKanji != nil {
		kanjiSubText := rawText[matchKanjiIndex[0]:matchKanjiIndex[1]]
		hiraSubText := rawText[matchKanjiIndex[1]:]
		matchHiraStation := regexpHiraStation.MatchString(kanjiSubText)
		if len(hiraSubText) == 0 || matchHiraStation == true {
			return matchKanji[1], KanaToHira(matchKanji[1])
		}
		regexpHira := regexp.MustCompile(`（(.+?)えき`)
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
			len := len(vowel)
			fmt.Fprintf(writeFp, "%v,%v,%v,%v\n", kanji, hira, vowel, len)
		}
	}
}

func GrowTree(reader *csv.Reader) {
	var RootTree *StationTree = new(StationTree)
	RootTree.ChildTree = map[string]*StationTree{}
	//fmt.Printf("%#v", *currentTree.ChildTree["a"])

	var regVowel = regexp.MustCompile(`yt`)

	record, err := reader.Read()
	if err != nil {
		panic(err)
	}
	word := Word{
		Kanji: string(record[0]),
		Hira:  string(record[1]),
		Vowel: string(record[2]),
	}
	len, _ := strconv.Atoi(record[3])
	fmt.Printf("%#v\n", word)

	var currentTree = RootTree
	for l := 0; l < len; l++ {
		if regVowel.MatchString(word.Vowel[l:l+1]) == true {
			if _, ok := (*currentTree).ChildTree[word.Vowel[l:l+2]]; ok {
				print("2chr is exist")
			}
		} else {
			if _, ok := (*currentTree).ChildTree[word.Vowel[l:l+1]]; ok {
				print("1chr is exist")
			} else {
				fmt.Printf("%v-tree is not exist. create it.\n", word.Vowel[0:l+1])
				var nextTree = new(StationTree)
				nextTree.Len = l + 1
				nextTree.CurrentVowel = word.Vowel[0 : l+1]
				nextTree.ChildTree = map[string]*StationTree{}
				(*currentTree).ChildTree[word.Vowel[l:l+1]] = nextTree
				fmt.Printf("%#v\n\n", (*currentTree))
				currentTree = nextTree
			}
		}
	}
	currentTree.WordList = append(currentTree.WordList, word)
	fmt.Printf("%#v\n", (*currentTree))

	/*
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			word := Word{
				Kanji: record[0],
				Hira: record[1],
				Vowel: record[2]
			}
			len := record[3]
			for l:=0;l<len;l++ {
				if word.Vowel[l]{

				}
			}
			//fmt.Fprintf(writeFp, "%v,%v,%v,%v\n", kanji, hira, vowel, len)

		}
	*/
}
func CsvToJson(readfile, writefile string) {
	readFp := OpenReadFile(readfile)
	defer readFp.Close()
	reader := csv.NewReader(readFp)

	writeFp := OpenWriteFile(writefile)
	defer writeFp.Close()

	GrowTree(reader)
}

func main() {

	//TextToCsv("raw_datalist.txt", "datalist.csv")
	CsvToJson("datalist.csv", "datalist.json")

}
