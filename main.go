package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
} 

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("time", 30, "The time limit for the quiz in seconds")
	flag.Parse()
	
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open csv file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Println("Failed to parse data from the provided csv file")
	}

	problems := parseLines(lines)
	
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	printProblem(problems, timer)
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

func printProblem( problems []problem, timer *time.Timer){
	correct := 0
	problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select{
			case <-timer.C:
				fmt.Println()
				break problemLoop
			case answer := <- answerCh:
				if answer == p.a {
					correct++
				}	
		}
	}
	fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}