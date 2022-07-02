package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	gpt3 "github.com/louisbarrett/gpt3client"
)

var (
	apiKey             = os.Getenv("OPEN_AI_APIKEY")
	flagUserInput      = flag.String("p", "generate golang code to run paralell commands", "prompt to send to gpt3 such as generate <lang> code to <something>")
	flagPromptSuffix   = flag.String("e", "and end the content with //QED.", "suffix to append to the prompt")
	flagIterations     = flag.Int("i", 1, "number of iterations to allow for content gen")
	flagOutputfileName = flag.String("o", "", "Output file name")

	prompt string
)

func init() {
	flag.Parse()
	prompt = fmt.Sprintf("%s %s", *flagUserInput, *flagPromptSuffix)
	if apiKey == "" {
		fmt.Println("Please set API key via OPEN_AI_APIKEY variable")
		os.Exit(1)
	}

}

func writeOutput(data string) {
	// Clean up the prompt that is returned by the API
	data = strings.Replace(data, *flagUserInput, "", -1)
	data = strings.Replace(data, *flagPromptSuffix, "", -1)
	// Write the generated code to a file
	ioutil.WriteFile(*flagOutputfileName, []byte(data), 0755)
}

func main() {
	iterations := *flagIterations
	inputPrompt := prompt
	var newContent string
	for i := 0; i < iterations; i++ {
		matchesEnd, err := regexp.Compile("QED")
		if err != nil {
			log.Fatal(err)
		}
		endofGen := matchesEnd.MatchString(newContent)
		if endofGen && *flagIterations != 0 {
			if *flagOutputfileName != "" {
				writeOutput(inputPrompt)
			}
			os.Exit(0)
		}

		// Send the prompt to the API
		inputPrompt, newContent = gpt3.SendOpenAIPrompt(inputPrompt)
		// Print the generated content
		fmt.Println(newContent)
	}

	if *flagOutputfileName != "" {
		writeOutput(inputPrompt)
	}
}
