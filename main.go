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

		/*
			if matched, _ := regexp.MatchString(`岩城`, rawText); matched {
							println(rawText, hiraSubText, matchKanji[1], KanaToHira(matchHira[1]))
						}
		*/
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
	var RootTree *StationTree = new(StationTree)
	RootTree.ChildTree = map[string]*StationTree{}

	for i := 0; i < 10000; i++ {
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
		var currentVowel = word.Vowel[0:1]
		var totalVowel = word.Vowel[0:1]
		var vowelSize = 1
		var smallVowelCount int
		var regSmallVowel = regexp.MustCompile(`[yt]`)
		for l := 0; l < wordLen; l++ {
			matchSmallVowel := regSmallVowel.MatchString(word.Vowel[l : l+1])
			if matchSmallVowel == true {
				vowelSize = 2
			} else {
				vowelSize = 1
			}
			currentVowel = word.Vowel[l : l+vowelSize]
			totalVowel = word.Vowel[:l+vowelSize]
			if _, ok := (*currentTree).ChildTree[currentVowel]; ok {
				//fmt.Printf("%v-tree is exist. move it.\n", totalVowel)
				currentTree = (*currentTree).ChildTree[currentVowel]
			} else {
				//fmt.Printf("%v-tree is not exist. create it.\n", totalVowel)
				var nextTree = new(StationTree)
				nextTree.Len = l + 1 - smallVowelCount
				nextTree.CurrentVowel = totalVowel
				nextTree.ChildTree = map[string]*StationTree{}
				(*currentTree).ChildTree[currentVowel] = nextTree
				//fmt.Printf("%#v\n\n", (*currentTree))
				currentTree = nextTree
			}

			if vowelSize == 2 {
				l++
				smallVowelCount++
			}

		}
		//matchSmallVowel := regSmallVowel.MatchString(currentTree.CurrentVowel[wordLen : wordLen])
		//if matchSmallVowel == true {
		currentTree.WordList = append(currentTree.WordList, word)
		//}
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
