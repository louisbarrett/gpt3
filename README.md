# gpt3

A tool for generating code from a prompt.

## Installation

```
go get -u github.com/louisbarrett/gpt3
```

## Usage

```
gpt3 -i 10 -p "generate golang code to run parallel commands" -o output.txt
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
