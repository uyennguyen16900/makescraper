package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
type Crypto struct {
	Name   string
	Price  string
	Symbol string
}

func main() {
	fName := "cryptocoinmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Price (USD)"})

	// Instantiate default collector
	c := colly.NewCollector()

	var coins []Crypto
	// var jsonText = []byte(`[
	//     {"Name": "", "Symbol": "", "Price": ""}
	// ]`)

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			crypto := Crypto{
				Name:   el.ChildText("td:nth-child(2)"),
				Symbol: el.ChildText("td:nth-child(3)"),
				Price:  el.ChildText(".price___3rj7O"),
			}

			writer.Write([]string{
				crypto.Name,
				crypto.Symbol,
				crypto.Price,
			})

			// if err := json.Unmarshal([]byte(jsonText), &coins); err != nil {
			// 	log.Println(err)
			// }
			coins = append(coins, crypto)

		})
		cryptoJSON, _ := json.MarshalIndent(coins, "", "")
		fmt.Println(cryptoJSON)
		_ = ioutil.WriteFile("output.json", cryptoJSON, 0644)
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
