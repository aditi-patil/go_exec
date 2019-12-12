package quiz

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestGame(t *testing.T) {

	content := []byte("10")
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
	os.Stdin = tmpfile
	Game("sample_quiz.csv", 10)
	if Correct != 1 {
		t.Errorf("You have scored %d out of 10\n", Correct)
	}
}

func TestReadFileData(t *testing.T) {
	// if provided filename is present
	records, _ := ReadFileData("sample_quiz.csv")
	fmt.Println(records)

	// if provided filename is not present
	_, error := ReadFileData("test.csv")
	if error != nil {
		t.Fatal(error)
	}
}

func TestParseRecords(t *testing.T) {
	records := [][]string{
		{"5 + 5", "10"},
	}
	quizzes := ParseRecords(records)
	if quizzes[0].answer != "10" {
		t.Errorf("Quiz answer of 5+5 is 10 but got %v", quizzes[0].answer)
	}
	fmt.Println(quizzes)
}

func TestGetQuizResult(t *testing.T) {
	problems := []Quiz{
		{"5 + 5", "10"},
	}

	content := []byte("10")
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
	os.Stdin = tmpfile
	Correct = 0
	GetQuizResult(problems, 12)
	if Correct != 1 {
		t.Errorf("You have scored %d out of %d\n", Correct, len(problems))
	}

	// if time out
	GetQuizResult(problems, 0)
	if Correct == 0 {
		fmt.Printf("You have scored %d out of %d\n", Correct, len(problems))
	}

	// Close the file
	if err := tmpfile.Close(); err != nil {
		fmt.Println(err)
	}
}
