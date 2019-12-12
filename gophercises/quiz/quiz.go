package quiz

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// Quiz is a struct which defines question and answer params
type Quiz struct {
	question string
	answer   string
}

// Correct is used when user has given correct answer of the question
var Correct = 0

// Game process the file and records and gives the result of the quiz
func Game(csvFile string, timeLimit int) {
	records, _ := ReadFileData(csvFile)
	problems := ParseRecords(records)

	GetQuizResult(problems, timeLimit)
}

// ReadFileData reads quiz records from csv file
func ReadFileData(filename string) ([][]string, error) {
	file, error := os.Open(filename)
	if error != nil {
		fmt.Printf("Not able to open csv %s file", filename)
	}

	data := csv.NewReader(file)

	records, err := data.ReadAll()

	if err != nil {
		fmt.Printf("Not able to read records from file: %s", filename)
	}

	return records, nil

}

// ParseRecords creates formatted quiz records from the csv generated records
func ParseRecords(records [][]string) []Quiz {
	quizzes := make([]Quiz, len(records))
	for i, record := range records {
		quizzes[i] = Quiz{
			question: record[0],
			answer:   strings.TrimSpace(record[1]),
		}
	}
	return quizzes
}

// GetQuizResult gets answers of the problem from user and gives final result
func GetQuizResult(problems []Quiz, timeLimit int) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	for j, problem := range problems {
		fmt.Printf("Question %d: %s =  \n", j+1, problem.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You have scored %d out of %d\n", Correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == problem.answer {
				Correct++
			}
		}

	}
	fmt.Printf("You have scored %d out of %d\n", Correct, len(problems))
}
