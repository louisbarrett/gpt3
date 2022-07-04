# gpt3

A tool for generating code from a prompt.

## Installation

```
export OPEN_AI_APIKEY=<api key>

go get -u github.com/louisbarrett/gpt3

OR 

go install github.com/louisbarrett/gpt3@latest
```

## Usage

`gpt3 -i 10 -p "generate golang code to run parallel commands" -o output.txt`

```go
func main() {
    cmd1 := exec.Command("echo", "Hello World")
    cmd2 := exec.Command("echo", "Goodbye World")

    err := cmd1.Start()
    if err != nil {
        log.Fatal(err)
    }

    err := cmd2.Start()
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Waiting for commands to finish...")
    err = cmd1.Wait()
    log.Printf("Command 1 finished with error: %v", err)

    err = cmd2.Wait()
    log.Printf("Command 2 finished with error: %v", err)

    log.Printf("All commands finished.")
}

//QED%
```

This will generate code that runs parallel commands. The cli will submit appended prompts 10 times. The output will be saved to `output.txt`.

## Options

```
  -e string
        suffix to append to the prompt (default "and end the content with //QED.")
  -i int
        number of iterations to allow for content gen (default 1)
  -o string
        Output file name
  -p string
        prompt to send to gpt3 such as generate <lang> code to <something> (default "generate golang code to run paralell commands")
```
