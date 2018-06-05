
# Corgi - CLI workflow manager

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/DrakeW/corgi/blob/master/LICENSE)

Corgi is a command-line tool that helps with your repetitive command usages by organizing them into reusable snippet. See usage by simply running `corgi` or `corgi --help`

Current version: **v0.1.2-alpha**
## Examples

Create a new snippet to automate the commands you run repetitively
<img src="images/corgi-new.gif" width="700">

Execute an existing snippet knowing what command is being run and its output
<img src="images/corgi-exec.gif" width="700">

# Table of Contents

- [Installation](#installation)
- [Features](#features)
    - [Create a snippet](#create-a-snippet)
    - [List snippets](#list-snippets)
    - [Describe a snippet](#describe-a-snippet)
    - [Execute a snippet](#execute-a-snippet)
    - [Edit a snippet](#edit-a-snippet)
    - [Share snippets](#share-snippets)
    - [Configure corgi](#configure-corgi)
- [Roadmap](#roadmap)
- [Note](#note)
  
## Installation  
Since this project is still in its very early phase, installation via package managers like `brew` or `apt-get` is not supported. Here are the steps to follow if you would like to try it out:
1. Download the latest `corgi` executable from releases tab
2. `chmod a+x ./corgi` to give execution permission to all users & groups
3. (Optional) If you already have a previous release of corgi installed, remove the soft link in your `bin` folder first
4. Create a soft link of the `corgi` executable to your local `bin` folder  - (if you are on Mac, you can use `ln -s $(pwd)/corgi /usr/local/bin/corgi`)
5. Start `Corgi`ing
  
## Features 
To view usage of a specific action, run `corgi <action> --help`  
  
### Create a snippet
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

### Describe a snippet
To see the details of a snippet, you can run
```
corgi describe <title of the snippet>
```
And it will print out each step of the snippet so that you don't have to memorize them.
  
### Execute a snippet  
To execute a snippet, simply run  
```  
corgi exec <title of the snippet>  
```  
Your commands will run smoothly just as they were being run on a terminal directly, and any prompt that asks for input (for example, password prompt when you `ssh` into a remote server) will work seamlessly as well.

Also note that if you run `cd` command in one of the steps, the current working directory will not change for subsequent steps. But you can always put `cd <target dir> && <do something else>` to get commands executed inside of your target directory. 
  
### Edit a snippet
To edit a snippet, run  
```  
corgi edit <title of the snippet>  
```  
You'll be able to edit the snippet json file directly with your preferred editor (configurable via `corgi config` command, details see below)

### Share snippets
If someone shares his/her snippet json file(s) with you, you can import it by running
```
corgi import <snippet json file 1> [<snippet json file 2>...]
```
And similarly, if you already have a workflow defined in a snippet, you can easily share it by exporting via
```
corgi export <title of the snippet> [-o <output file path>]
```
and send the json file to another person
  
### Configure `corgi` 
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