
# Corgi - CLI workflow manager

[![CodeFactor](https://www.codefactor.io/repository/github/drakew/corgi/badge)](https://www.codefactor.io/repository/github/drakew/corgi)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/DrakeW/corgi/blob/master/LICENSE)

<img src="images/corgi-logo.png" width="200">

Corgi is a command-line tool that helps with your repetitive command usages by organizing them into reusable snippet. See usage by simply running `corgi` or `corgi --help`

Current version: **v0.2.2**
## Examples

Create a new snippet to automate the commands you run repetitively
<img src="images/corgi-new.gif" width="900">

Execute an existing snippet knowing what command is being run and its output
<img src="images/corgi-exec.gif" width="900">

# Table of Contents

- [Installation](#installation)
    - [Install corgi](#install-corgi)
        - [Homebrew](#homebrew)
        - [Binary](#binary)
        - [Build from scratch](#build-from-scratch)
    - [Install a fuzzy-finder](#install-a-fuzzy-finder)
- [Features](#features)
    - [Create a snippet](#create-a-snippet)
        - [Create from command history](#create-from-command-history)
        - [Define template fields in snippet](#define-template-fields-in-snippet)
    - [List snippets](#list-snippets)
    - [Execute a snippet](#execute-a-snippet)
        - [Use default value without prompt](#use-default-value-without-prompt)
        - [Select steps to execute](#select-steps-to-execute)
        - [Interactive snippet selection](#interactive-snippet-selection)
    - [Edit a snippet](#edit-a-snippet)
    - [Describe a snippet](#describe-a-snippet)
    - [Share snippets](#share-snippets)
    - [Remove a snippet](#remove-a-snippet)
    - [Configure corgi](#configure-corgi)
- [Roadmap](#roadmap)
- [Note](#note)
  
## Installation

### Install `corgi`

#### Homebrew
Run `brew install drakew/corgi/corgi` and you are good to roll!

#### Binary
To install the binary directly:
1. Download the latest package tarball from [releases](https://github.com/DrakeW/corgi/releases) based on your system (currently support linux and macOS).
2. Unzip the package with `tar -xzvf <downloaded tar.gz file>`
3. `cd <directory extracted>`, and `chmod a+x ./corgi` to make sure execution bit is turned.
4. `ln -s $(pwd)/corgi <your bin folder>/corgi` to create a softlink
5. Start rolling with `corgi`

#### Build from scratch
To build the executable from scratch:
1. Install the most recent version of `go` language on your system
2. `git clone git@github.com:DrakeW/corgi.git`
3. Run `go build -o corgi` inside of the repo
4. `ln -s $(pwd)/corgi <your bin folder>/corgi`

### Install a fuzzy-finder
`corgi` will enable interactive selection if you install a fuzzy finder, the currently supported two are [fzf](https://github.com/junegunn/fzf) and [peco](https://github.com/peco/peco).
  
## Features 
To view usage of a specific action, run `corgi <action> --help`  
  
### Create a snippet
corgi provides an interactive CLI interface to create snippet, and you can start by running  
```  
corgi new  
```  

#### Create from command history
If you would like to quickly combine the last couple commands you just executed into a snippet, you could also run  
```  
corgi new --last <number of commands to look back>  
```  

#### Define template fields in snippet
Furthermore, you could also add template fields (with or without default value) to command of a step and reuse the same field, for example:  
```  
tar -xzf <project>.tgz && scp <project> <user=ec2-user>@<ec2-instance-address>:~/
```
And you will be prompted to enter values for those fields when the snippet executes. The value set for the same field name will be applied to **all steps** in a snippet, so you don't have to enter multiple times.

Also if you have field with **multiple default values**, the latest appearance will take precedence over the previous values.


### List snippets  
To view all snippets saved on your system, run  
```  
corgi list  
```  
  
### Execute a snippet
To execute a snippet, simply run  
```  
corgi exec [<title of the snippet>] [--use-default] [--step <step range>]
```  
Your commands will run smoothly just as they were being run on a terminal directly, and any prompt that asks for input (for example, password prompt when you `ssh` into a remote server) will work seamlessly as well.

Also note that if you run `cd` command in one of the steps, the current working directory will not change for subsequent steps. But you can always put `cd <target dir> && <do something else>` to get commands executed inside of your target directory.

#### Use default value without prompt
if `--use-default` is set, then the snippet will execute without asking explicitly for user to set a value for template fields defined, but if there are missing default values, the snippet will fail fast and tell you what fields are missed.

#### Select steps to execute
You can use the `--step` (or `-s`) flag to specify the steps (starting from index 1) you want to execute, for example
```
--step 3    # will only execute step 3
--step 3-5  # will execute step 3 to 5
--step 3-   # will execute step 3 to the last step
```
This can be particularly useful when your workflow fail midway, but you don't want to start the whole workflow from step 1 again.

#### Interactive snippet selection
This feature will guide you through the darkness if you don't have the title of your snippet memorized. Simply type
```
corgi exec [with or without options]
```
and you will be presented an interactive interface to fuzzy-find your snippet (To enable this feature, see [installation instruction](#install-a-fuzzy-finder))

<img src="images/corgi-fuzzy-select.gif" width="700">


### Edit a snippet
To edit a snippet, run  
```  
corgi edit [<title of the snippet>]  
```  
You'll be able to edit the snippet json file directly with your preferred editor (configurable via `corgi config` command, details see below)

Furthermore, `edit` also provides [fuzzy finding capabilities](#interactive-snippet-selection) when you omit the snippet title.

### Describe a snippet
To see the details of a snippet, you can run (without title will yield interactive title selection)
```
corgi describe [<title of the snippet>]
```
And it will print out each step of the snippet so that you don't have to memorize them.

### Share snippets
If someone shares his/her snippet json file(s) with you, you can import it by running
```
corgi import <snippet json file 1> [<snippet json file 2>...]
```
And similarly, if you already have a workflow defined in a snippet, you can easily share it by exporting via (without title will yield interactive title selection)
```
corgi export [<title of the snippet>] [-o <output file path>]
```
and send the json file to another person

### Remove a snippet
You can remove a snippet by running (without title will yield interactive title selection)
```
corgi remove [<title of the snippet>]
```
  
### Configure `corgi` 
You can configure your text editor choice (default is `vim`) for snippet editing, and filter tool for interactive snippet selection. To configure the corgi CLI, run  
```  
corgi config [--editor <editor of your choice>] [--filter-cmd <fzf, peco or something else>]
```  
  
## Roadmap  
Here are some features that are currently on the roadmap:  
1. Support concurrent execution of steps  
2. Support remote server configuration, so that snippet can run seamlessly on a remote computer

## Note
Corgi is inspired by [Pet](https://github.com/knqyf263/pet), and aims to advance Pet's command-level usage to a workflow level.