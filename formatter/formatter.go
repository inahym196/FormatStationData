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
	"strings"
	"unicode"

	"../src/roma"
	"../src/stationTree"

	//"github.com/inahym196/FormatStationData/src/stationTree"
	"github.com/inahym196/gojaconv/jaconv"
)

func OpenReadFile(filename string) io.ReadCloser {
	fp, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return fp
}

func OpenWriteFile(filename string) io.WriteCloser {
	fp, err := os.Create(filename)
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
			romas := roma.InitRomas(jaconv.ToHebon(hira))
			//len := len(romas) //- strings.Count(vowel, "y") - strings.Count(vowel, "t")
			fmt.Fprintf(writeFp, "%v,%v,%v\n", kanji, hira, romas /*, len*/)
		}
	}
}
func CsvToJson(readfile, writefile string) {
	readFp := OpenReadFile(readfile)
	defer readFp.Close()
	reader := csv.NewReader(readFp)

	writeFp := OpenWriteFile(writefile)
	defer writeFp.Close()

	var RootTree = stationTree.NewStationTree("")
	//forStart := 0
	//forCount := 10000
	for /*debug := 0; debug < forStart+forCount; debug++*/ {
		record, err := reader.Read()
		if err == io.EOF {
			fmt.Printf("EOF")
			break
		}
		if err != nil {
			panic(err)
		}
		/*
			if debug < forStart {
				continue
			}
		*/
		RootTree.GrowTree(record)
	}
	jsonData, err := json.MarshalIndent(RootTree, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(writeFp, "%v", string(jsonData))
}

func main() {

	//TextToCsv("raw_datalist.txt", "datalist.csv")
	CsvToJson("datalist.csv", "datalist.json")

}
