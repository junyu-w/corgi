
# Corgi  
  
Corgi is a command-line tool that helps with your repetitive command usages by organizing them into reusable snippet. See usage by simply running `corgi` or `corgi --help` 
  
## Installation  
Since this project is still in its very early phase, installation via package managers like `brew` or `apt-get` is not supported. Here are the steps to follow if you would like to try it out:
1. Download the latest package from release tab
2. Create a soft link of the `corgi` executable to your local `bin` folder  - (if you are on Mac, you can use `ln -s ./corgi /usr/local/bin/corgi`)
3. Start `Corgi`ing
  
## Features 
To view usage of a specific action, run `corgi <action> --help`  
  
### Create snippet (`corgi new --help`)  
corgi provides an interactive CLI interface to create snippet, and you can start by running  
```  
corgi new  
```  
If you would like to quickly combine the last couple commands you just executed into a snippet, you could also run  
```  
corgi new --last <number of commands to look back>  
```  
Furthermore, you could also save a command template (with or without default value) as part of the snippet, for example:  
```  
ssh -i <aws-key-file> <user=ec2-user>@<ec2-instance-address>  
```  
And you can enter the values for those parameters when the snippet executes.  
  
### List snippets  
To view all snippets saved on your system, run  
```  
corgi list  
```  
  
### Execute snippet  
To execute a snippet, simply run  
```  
corgi exec --title <title of the snippet>  
```  
  
### Edit snippet  
To edit a snippet, run  
```  
corgi edit --title <title of the snippet>  
```  
  
### Configure CLI  
Currently the only editable option is your text editor choice (default is `vim`), to configure the corgi CLI, run  
```  
corgi config --editor <editor of your choice>  
```  
  
## Roadmap  
Here are some features that are currently on the roadmap:  
1. Support concurrent execution of steps  
2. Support remote server configuration, so that snippet can run seamlessly on a remote computer

## Note
Corgi is inspired by [Pet](https://github.com/knqyf263/pet), and aims to advance Pet's command-level usage to a workflow level.