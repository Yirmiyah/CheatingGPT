package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

const EXERCISES_NAMES = "all.txt"
const PATH = "../../../piscine-go/piscine/"
const NAME_PATH = "all.txt"
const EXTENSION = ".go"
const FILE_SOLUTION = "/home/student07/projet_perso/Big_Projects/Cheating/excercise_solution.txt"
const TEST_FILE = "/test.txt"

var countResponse int

func main() {

	gpt3Response := "nothing..."
	cmd0 := exec.Command("xsel", "-b", "-d")
	err := cmd0.Run()
	if err != nil {
		fmt.Println("Erreur exec.Command delete primary buffer!")
	}
	GetNameFile()
	for {
		time.Sleep(time.Second)
		if gpt3Response != WriteCtrlV() {
			gpt3Response = WriteCtrlV()

			fmt.Printf("gpt3Response: %v\n", gpt3Response)

		} else {
			cmd0 := exec.Command("xsel", "-b", "-d")
			err := cmd0.Run()
			if err != nil {
				fmt.Println("Erreur exec.Command delete primary buffer!")
			}
			time.Sleep(time.Second)
			gpt3Response = "nothing..."
		}
	}
}

func WriteCtrlV() string {

	cmd2 := exec.Command("xsel", "-b", "-o")
	var out1 bytes.Buffer
	cmd2.Stdout = &out1
	err := cmd2.Run()
	if err != nil {
		fmt.Println("erreur copy content into the primary buffer")

	}
	clipboard := out1.String()

	// fmt.Printf("clipboard content: %v\n", clipboard)

	if strings.Contains(clipboard, "-gpt3") {

		client := openai.NewClient("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")
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: clipboard,
					},
				},
			},
		)
		if err != nil {
			fmt.Println("Error GPT3 Response...")
			log.Fatal(err)
		}

		gpt3Response := resp.Choices[0].Message.Content
		countResponse++
		fmt.Printf("countResponse: %v\n", countResponse)

		if gpt3Response != "" {
			fmt.Printf("gpt3Response: %v\n", gpt3Response)
			gpt3 := strings.NewReader(gpt3Response)

			cmd2 = exec.Command("xsel", "-i", "-b")
			cmd2.Stdin = gpt3
			if err := cmd2.Run(); err != nil {
				fmt.Println("erreur cmd2.Run()")
			}
			return gpt3Response
		} else {
			return "nothing..."
		}

	} else {
		f, err := os.ReadFile(EXERCISES_NAMES)
		if err != nil {
			fmt.Println("Name not found")
		}
		element := strings.Split(string(f), "\n")

		cmd2 := exec.Command("xsel", "-b", "-o")
		var out1 bytes.Buffer
		cmd2.Stdout = &out1
		err = cmd2.Run()
		if err != nil {
			fmt.Println("erreur copy content into the primary buffer")

		}
		clipboard := out1.String()
		textClipboard := strings.ToLower(clipboard)

		if textClipboard != "" {

			for _, v := range element {

				if textClipboard == v {

					fmt.Printf("name_exercise: %v\n", v)

					r, err := os.ReadFile(PATH + v + EXTENSION)
					if err != nil {
						fmt.Println("erreur ReadFile")

					}

					cmd := exec.Command("echo", "-n", string(r))
					stdout, err := cmd.StdoutPipe()
					if err != nil {
						fmt.Println("erreur echo")

					}
					if err := cmd.Start(); err != nil {
						fmt.Println("erreur Start cmd")
					}

					cmd2 := exec.Command("xsel", "-i", "-b")
					cmd2.Stdin = stdout
					if err := cmd2.Run(); err != nil {
						fmt.Println("erreur cmd2.Run()")
					}

					if err := cmd.Wait(); err != nil {
						fmt.Println("erreur wait cmd")

					}
					return "nothing..."
				}
			}
		}

		if textClipboard == "all" {

			fmt.Printf("name_exercise: %v\n", textClipboard)

			r, err := os.ReadFile(textClipboard + ".txt")
			if err != nil {
				fmt.Println("erreur ReadFile")

			}

			cmd := exec.Command("echo", "-n", string(r))
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Println("erreur echo")

			}
			if err := cmd.Start(); err != nil {
				fmt.Println("erreur Start cmd")
			}

			cmd2 := exec.Command("xsel", "-i", "-b")
			cmd2.Stdin = stdout
			if err := cmd2.Run(); err != nil {
				fmt.Println("erreur cmd2.Run()")

			}

			if err := cmd.Wait(); err != nil {
				fmt.Println("erreur wait cmd")

			}
			return "nothing..."
		}

	}

	return "nothing..."
}

func GetNameFile() {
	dir := PATH // replace with the path to your folder
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	f, err := os.Create("all.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()
	for _, file := range files {
		if !file.IsDir() {
			filename := filepath.Base(file.Name())
			name := strings.Split(filename, ".")
			fmt.Fprintln(f, name[0])
		}
	}
	fmt.Println("File names written to all.txt")

}
