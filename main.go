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
	apiKey                       = os.Getenv("OPEN_AI_APIKEY")
	flagStdIn                    = flag.Bool("s", false, "Use stdin as input")
	flagCombine                  = flag.Bool("c", false, "Combine stdin and input prompt into one, stdin will be below the input prompt")
	flagStripInputPromptFromFile = flag.Bool("strip", true, "Strip the input prompt from the input file")
	flagIterations               = flag.Int("i", 1, "number of iterations to allow for content gen")
	flagUserInput                = flag.String("p", "", "prompt to send to gpt3 such as generate <lang> code to <something>")
	flagPromptSuffix             = flag.String("e", "", "suffix to append to the prompt")
	flagInputfileName            = flag.String("f", "", "input file to read from")
	flagOutputfileName           = flag.String("o", "", "Output file name")

	prompt         string
	promptFromFile string
	err            error
	stdinData      = ""
)

func init() {
	// Parse the flags
	flag.Parse()
	// If the user input is not empty, use that as the prompt
	if *flagUserInput != "" {
		prompt = *flagUserInput
	}
	// if we have data on stdin use the as the prompt
	if *flagStdIn {
		stdinData, err = getPipedData()
		if err != nil {
			log.Println(err)
		}
		prompt = stdinData
	}
	// if we have a file to read from, read it and use it as the prompt
	if *flagInputfileName != "" {
		promptFromFile, err = readInputFile()
		if err != nil {
			log.Println(err)
		}
		prompt = promptFromFile
	}

	// if multiple input are given but combine is not set, combine them into one

	if prompt == "" && stdinData == "" && promptFromFile == "" {
		log.Fatal("Please provide a prompt")
		flag.Usage()
	}

	inputFromParam := len(*flagUserInput) > 0
	inputFromFile := len(*flagInputfileName) > 0
	inputFromStdin := len(stdinData) > 0

	// convert boolean to int
	var inputFromParamInt int
	if inputFromParam {
		inputFromParamInt = 1
	}
	var inputFromFileInt int
	if inputFromFile {
		inputFromFileInt = 1
	}
	var inputFromStdinInt int
	if inputFromStdin {
		inputFromStdinInt = 1
	}
	shouldCombine := inputFromParamInt+inputFromFileInt+inputFromStdinInt > 1

	// Write warning that combine flag was not set
	if shouldCombine && !*flagCombine {
		// print to error stream
		// color the output orange
		fmt.Fprintln(os.Stderr, "\033[33mWarning: combining input prompts is not set, only the last input prompt will be used")
	}

	if *flagCombine {
		prompt = *flagUserInput + "\n" + stdinData + "\n" + promptFromFile
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
	if *flagStripInputPromptFromFile {
		data = strings.Replace(data, promptFromFile, "", -1)
	}
	// Write the generated code to a file
	ioutil.WriteFile(*flagOutputfileName, []byte(data), 0755)
}

func getPipedData() (pipedData string, err error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		pipedData += line
	}

	return pipedData, err
}

func readInputFile() (inputPrompt string, err error) {
	inputFile, err := os.Open(*flagInputfileName)
	if err != nil {
		return "", err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		inputPrompt += scanner.Text()
	}

	return inputPrompt, nil
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
		fmt.Println("\033[96m", newContent)
	}

	if *flagOutputfileName != "" {
		writeOutput(inputPrompt)
	}
}
