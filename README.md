# OPEN AI terminal assistant


## What's this ?

This is an assistant for your terminal, using [OpenAI ChatGPT](https://chat.openai.com/) to build and run commands for you. You just need to describe them in your everyday language, it will take care or the rest. 

## Installation

To install the project, you need to have Go installed on your machine. Once you have Go installed, you can clone the repository and install the dependencies.

```sh
git clone https://github.com/AkhilSharma90/go-openai-terminal-assistant.git
cd repository
```

To run the project, you can use the go run command:
```
go run main.go
```

## How to get started?

Create a .config folder - 
```
mkdir ~/.config
cd ~/.config
```

create a config file here - ```terminal-assistant.json```
and mention the following details -

```
{
    "openai_key": "REPLACE WITH YOUR OPEN AI KEY",      
    "openai_model": "gpt-3.5-turbo",   
    "openai_proxy": "",               
    "openai_temperature": 0.0,        
    "openai_max_tokens": 2000,         
    "user_default_prompt_mode": "exec",
    "user_preferences": ""             
  }
```

## Testing
This project includes unit tests for the various modules. You can run these tests using the go test command. For example, to run the tests for the history module, you can use the following command:

```go test ./history```

You can run all tests in the project using the following command:

```go test ./...```