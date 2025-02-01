package main

import (
	"fmt"
	"github.com/peeta98/blog-aggregator/internal/config"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	if err = cfg.SetUser("Peeta"); err != nil {
		log.Fatal(err)
	}

	updatedUsrCfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(updatedUsrCfg.DbUrl)
}
