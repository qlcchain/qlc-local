package main

import (
	"github.com/qlcchain/qlc-local/template"
	"log"
)

func main() {
	if _, _, err := template.Template("/Users/sidney/Desktop/testnet-compose.yml", 0, 0, 2, "", "0.10.5"); err != nil {
		log.Fatal(err)
	}
}
