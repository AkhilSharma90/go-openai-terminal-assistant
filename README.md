# OPEN AI terminal assistant


## What's this ?

This is an assistant for your terminal, using [OpenAI ChatGPT](https://chat.openai.com/) to build and run commands for you. You just need to describe them in your everyday language, it will take care or the rest. 

## How to get started?

Create a .config folder - 
```mkdir ~/.config
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