package main

import (
	"flag"
	"gophercises/quiz"
)

func main() {
	csvFile := flag.String("csv", "problems.csv", "A csv file containing problems and answers.")
	timeLimit := flag.Int("limit", 5, "the time limit for the quiz in seconds")
	flag.Parse()
	quiz.Game(*csvFile, *timeLimit)
}
