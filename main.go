package main

import (
	"bufio"
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
	flagUserInput      = flag.String("p", "", "prompt to send to gpt3 such as generate <lang> code to <something>")
	flagPromptSuffix   = flag.String("e", "and end the content with //QED.", "suffix to append to the prompt")
	flagIterations     = flag.Int("i", 1, "number of iterations to allow for content gen")
	flagOutputfileName = flag.String("o", "", "Output file name")
	flagStdIn          = flag.Bool("s", false, "Use stdin as input")
	prompt             string
	err                error
	stdinData          = ""
)

func init() {
	flag.Parse()
	if *flagStdIn {
		stdinData, err = getPipedData()
		if err != nil {
			log.Println(err)
		}
		prompt = stdinData

	} else {
		prompt = *flagUserInput
	}

	if prompt == "" && stdinData == "" {
		log.Fatal("Please provide a prompt")
		flag.Usage()
	}

	// prompt = fmt.Sprintf("%s", *flagUserInput)

	if *flagPromptSuffix != "" {
		prompt = fmt.Sprintf("%s %s", prompt, *flagPromptSuffix)
	}

	if apiKey == "" {
		log.Fatal("Please set API key via OPEN_AI_APIKEY variable")
	}

}

func writeOutput(data string) {
	// Clean up the prompt that is returned by the API
	data = strings.Replace(data, *flagUserInput, "", -1)
	data = strings.Replace(data, *flagPromptSuffix, "", -1)
	// Write the generated code to a file
	ioutil.WriteFile(*flagOutputfileName, []byte(data), 0755)
}

func getPipedData() (pipedData string, err error) {
	reader := bufio.NewReader(os.Stdin)

	pipedData, err = reader.ReadString('\n')
	if err != nil {
		return "", err

	}
	return pipedData, err
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
