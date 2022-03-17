package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	fileFlag := flag.String("file", "problem.csv", "load csv for question set")
	timeFlag := flag.Int("time", 30, "time limit for quiz in seconds")
	shuffleFlag := flag.Bool("shuffle", false, "enable random question order")
	flag.Parse()

	file, err := os.Open(*fileFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Handle CTRL+C for exit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Println()
		fmt.Println("Exiting...")
		file.Close()
		os.Exit(1)
	}()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if *shuffleFlag {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(rows), func(i, j int) { rows[i], rows[j] = rows[j], rows[i] })
	}
	total := len(rows)
	correct := 0

	fmt.Printf("You have %ds to answer %d questions.\nPress <ENTER> key to start quiz, <CTRL+C> to quit.", *timeFlag, total)
	fmt.Scanln()

	timer := time.NewTimer(time.Second * time.Duration(*timeFlag))
	defer timer.Stop()
	go func() {
		<-timer.C
		fmt.Printf("\n%ds up\n", *timeFlag)
		fmt.Printf("%d/%d correct\n", correct, total)
		file.Close()
		os.Exit(0)
	}()

	for _, r := range rows {
		fmt.Printf("%s=", r[0])

		var answer string
		fmt.Scanln(&answer)

		if strings.ToLower(strings.TrimSpace(answer)) == strings.ToLower(strings.TrimSpace(r[1])) {
			correct += 1
		}
	}
	fmt.Printf("%d/%d correct\n", correct, total)
}
