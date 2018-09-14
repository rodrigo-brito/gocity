package main

import (
	"fmt"
	"log"
)

func main() {
	packageName := "github.com/rodrigo-brito/go-async-benchmark"

	analyzer := NewAnalyzer(packageName)
	err := analyzer.FetchPackage()
	if err != nil {
		log.Fatal(err)
	}

	summary, err := analyzer.Analyze()
	if err != nil {
		log.Fatalf("error on analyzetion %s", err)
	}

	for key, value := range summary {
		fmt.Printf("%s: LOC=%d NOM=%d NOA=%d\n", key, value.NumberLines, value.NumberFunctions, value.NumberAttributes)
	}
}
