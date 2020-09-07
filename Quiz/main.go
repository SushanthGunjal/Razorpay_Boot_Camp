package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// user defined type to store question and answer
type problem struct {
	q string
	a string
}

// function to shuffle the array every time the program is executed
func Shuffle(problems []problem) []problem {
	// Seed function to generate different shuffle pattern every time.
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})

	return problems
}

// function to convert the contents of csv to the slice of struct type 'problem'
func parseFile(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret

}

func main() {

	// Flag to take the input from command line arguments.
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "set true to shuffle the elements")
	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	problems := parseFile(lines)

	if *shuffle == true {
		problems = Shuffle(problems)
	}

	// Create a timer to add the timelimit.
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct_answer := 0

	for i, p := range problems {
		fmt.Printf("Problem number %d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s \n", &answer)
			answerCh <- answer
		}()

		select {
		/* This case is executed if timer runs out. */
		case <-timer.C:
			fmt.Printf("\n Oops Time up!!! \n You scored %d out of %d \n ", correct_answer, len(problems))
			return

		/* This case is executed if the input is recieved from the answerCh channel. */
		case answer := <-answerCh:
			if answer == p.a {
				correct_answer++
			}
		}
	}
	fmt.Printf("\n Good Game \n Yeah! You scored %d out of %d \n ", correct_answer, len(problems))
}
