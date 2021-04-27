package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func xor(input string, key []byte) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i%len(key)])
	}
	return output
}

type Language struct {
	HelperFunction string
	Comment        string
	FormatFunction func(string, string, []byte) string
}

var lang = map[string]Language{
	"go": Language{
		HelperFunction: `
func xor(input string, key []byte) (output string) {
    for i := 0; i < len(input); i++ {
        output += string(input[i] ^ key[i%len(key)])
    }
    return output
}`,
		Comment: "//",
		FormatFunction: func(plain, cover string, key []byte) string {
			return fmt.Sprintf("%s := xor(\"%s\", %#v) // %s", strings.Replace(plain, " ", "", -1), cover, key, plain)
		},
	},
	"powershell": Language{
		HelperFunction: `
function xor {
    Param(
        [Parameter(Position = 0, Mandatory = $True)] $plain,
        [Parameter(Position = 1, Mandatory = $True)] $key
    )
    $r = ""
    $plaintext.ToCharArray() | foreach-object -process {
        $r += [char]([byte][char]$_ -bxor $key[$r.Length % $key.Length])
    }
    return $r
}`,
		Comment: "#",
		FormatFunction: func(plain, cover string, key []byte) string {
			keyStr := ""
			for _, k := range key {
				keyStr += "0x" + fmt.Sprintf("%x", k) + ","
			}
			keyStr = keyStr[0 : len(keyStr)-1]
			return fmt.Sprintf("$%s = xor -plain \"%s\" -key @(%s)",
				strings.Replace(plain, " ", "", -1),
				cover,
				keyStr,
			)
		},
	},
}

func main() {
	var (
		plain    = "AmsiScanBuffer"
		wordlist = "words.txt"
		matches  = 10
		language = "go"
	)
	flag.StringVar(&plain, "plain", plain, "Plain text to hide")
	flag.StringVar(&wordlist, "wordlist", wordlist, "Wordlist for cover text")
	flag.IntVar(&matches, "matches", matches, "Number of matches")
	flag.StringVar(&language, "language", language, "Language for output")
	flag.Parse()

	file, err := os.Open(wordlist)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Printf("%s Helper xor function:%s\n", lang[language].Comment, lang[language].HelperFunction)

	fmt.Printf("%s Covers:\n", lang[language].Comment)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cover := scanner.Text()
		//if len(cover) >= len(plain) && len(cover) < len(plain)*2 {
		if len(cover) == len(plain) {
			var key = []byte(xor(cover, []byte(plain)))
			fmt.Println(lang[language].FormatFunction(plain, cover, key))
			matches--
			if matches == 0 {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
