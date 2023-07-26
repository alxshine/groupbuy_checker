package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := flag.String("url", "https://candykeys.com/group-buys/updates", "URL to parse")
	targets := flag.String("targets", "awaken,terror below", "Comma-separated list of targets to look for")
	target_file := flag.String("target-file", "", "File containing targets to look for (overwrites -targets)")
	flag.Parse()

	var interested = strings.Split(*targets, ",")

	if *target_file != "" {
		targetFileContent, err := os.ReadFile(*target_file)
		if err != nil {
			panic(err)
		}
		interested = strings.Split(string(targetFileContent), "\n")
	}

	res, err := http.Get(*url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.keyboards-single-info").Each(func(_ int, keyboardInfo *goquery.Selection) {
		name := keyboardInfo.Find("h3").Text()

		contained := false
		for _, target := range interested {
			if strings.Contains(strings.ToLower(name), target) {
				contained = true
				break
			}
		}
		if !contained {
			return
		}

		shippingInfo := keyboardInfo.Find("div.ship-date").Children().Map(func(i int, s *goquery.Selection) string { return s.Text() })
		shippingExpected := shippingInfo[0:2]
		shippingActual := shippingInfo[2:4]

		fmt.Printf("%s:\n\t%s\n\t%s\n", name, shippingExpected, shippingActual)
	})
}
