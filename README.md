
# Corgi - CLI workflow manager

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/DrakeW/corgi/blob/master/LICENSE)

Corgi is a command-line tool that helps with your repetitive command usages by organizing them into reusable snippet. See usage by simply running `corgi` or `corgi --help`

## Examples

Create a new snippet to automate the commands you run repetitively
<img src="images/corgi-new.gif" width="700">

Execute an existing snippet knowing what command is being run and its output
<img src="images/corgi-exec.gif" width="700">
  
## Installation  
Since this project is still in its very early phase, installation via package managers like `brew` or `apt-get` is not supported. Here are the steps to follow if you would like to try it out:
1. Download the latest `corgi` executable from releases tab
2. `chmod a+x ./corgi` to give execution permission to all users & groups
3. Create a soft link of the `corgi` executable to your local `bin` folder  - (if you are on Mac, you can use `ln -s $(pwd)/corgi /usr/local/bin/corgi`)
4. Start `Corgi`ing
  
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
Furthermore, you could also save a command template (with or without default value) as part of the snippet and reuse the same parameter, for example:  
```  
tar -xzf <project>.tgz && scp <project> <user=ec2-user>@<ec2-instance-address>:~/
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