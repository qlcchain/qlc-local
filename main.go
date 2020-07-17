package main

import (
	"github.com/qlcchain/qlc-local/template"
	"log"
)

func main() {
	if _, _, err := template.Template("/Users/xxx/Desktop/testnet-compose.yml", 3, 0, 3, "", "0.10.5"); err != nil {
		log.Fatal(err)
	}
}
