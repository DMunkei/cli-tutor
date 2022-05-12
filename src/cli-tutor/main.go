package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"

	"cli-tutor/src/input"
	"cli-tutor/src/lesson"
	"cli-tutor/src/printer"
	"cli-tutor/src/prompt"

	"github.com/charmbracelet/glamour"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
)

func main() {
	// NOTE: This seems unix only needs to be tested
	logFile, err := os.OpenFile("tutor-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	temp := template.Must(template.New("lesson1.md").Funcs(lesson.FuncMap).ParseFiles("./lessons/lesson1.md"))
	f, _ := os.Create("expanded.md")
	defer os.Remove(f.Name())

	err = temp.Execute(f, "")
	if err != nil {
		log.Panic(err)
	}

	content, _ := os.ReadFile("expanded.md")
	currentLesson := lesson.ParseLesson(content)

	termenv.ClearScreen()
	printer.Print("Welcome to Chistole", "welcome")
	time.Sleep(1 * time.Second)
	termenv.ClearScreen()
	printer.Print("This lesson is titled:", "tip")
	printer.Print(currentLesson.Name, "")
	printer.Print("\n"+currentLesson.Description, "guide")
	printer.Print("\n\nWhen you are ready press the enter key to begin", "note")
	fmt.Scanln() // Any key
	termenv.ClearScreen()
	printer.Print("Welcome to the shell", "tip")
	printer.Print("Try out some commands or type 'exit'/'quit' to quit the shell", "note")
	time.Sleep(1 * time.Second)

	rl, err := readline.New(prompt.BuildPrompt() + " > ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	// Markdown renderer
	r, _ := glamour.NewTermRenderer(
		glamour.WithPreservedNewLines(),
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(130),
	)

	// Readline loop
	currentTask := 0
	for {
		log.Print(rl.Config.Prompt)
		if currentTask < 0 {
			currentTask = 0
		}

		if currentTask >= len(currentLesson.Tasks) {
			printer.Print("All Task Completed!, Exiting tutor!", "note")
			break
		}

		tracker := fmt.Sprintf("\n\n%s : Current Task [%d/%d]:", currentLesson.Name, currentTask, len(currentLesson.Tasks)-1)
		printer.Print(tracker, "note")

		out, _ := r.Render(currentLesson.Tasks[currentTask].Description)
		printer.Print(out, "")

		line, err := rl.Readline()
		printer.Print(line, "")
		if err != nil { // io.EOF
			break
		}

		if line == "" {
			continue
		}

		if line == "next" || line == "n" {
			currentTask++
			continue
		}

		if line == "prev" || line == "p" {
			currentTask--
			continue
		}

		if line == "quit" || line == "exit" {
			printer.Print("Exiting tutor, Good bye!", "note")
			break
		}

		command := input.InputFilter(line, currentLesson.Vocabulary)

		var cmd *exec.Cmd

		if len(command) > 0 {
			if len(command) > 1 {
				args := command[1:]
				cmd = exec.Command(command[0], args...)
			} else {
				cmd = exec.Command(command[0])
			}
			output, _ := cmd.CombinedOutput()
			printer.Print(string(output), "")
		} else {
			printer.Print("Let's stick to the basics", "error")
		}
	}
}
