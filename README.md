# reapy

Achieve complex tasks by breaking them down in small steps and guiding the user through it step by step.

## Concept

The idea is to break down complex tasks into their individual steps.

A typical installation of software might consist out of the following steps:

1. Download the software
2. Unzip the archive the software is distributed int
3. Copy the unzipped files to the target folder
4. Generate/Adjust the configuration
5. Install the software (i.e. as a windows service)

Such a process is defined inside of a YAML file. It consits out of the following 

To break the process down we distinguish the following building blocks:

* Plan  
The overaching goal we want to achieve. (i.e. installing a piece of software). Each plan consists out of one ore more tasks.
* Task  
Each plan consists out of multiple tasks. A task could be the update of a database and checking that the update was successful. Each task therefore consists out of one or more steps and checks.
* Step
A small piece of work contributing to the fulfillment of a task or a check to ensure that a previous step was executed successfull. (i.e. unzipping a file or copying a folder from one place to another) (i.e. having the user manually confirm something or checking that a log file contains no errors)

### MVP

The following gives an overview of the bare minimum that needs to work in order for this to be a usable tool.

* [ ] Read plan definition from YAML file
* [ ] Protocol the execution of a plan in a file
* [ ] Validate the plan definition file
* [ ] Allow for the definition of variables that can be used in the YAML file

#### Steps

* [ ] Unzip file
* [ ] Copy file/folder
* [ ] Delete file/folder
* [ ] Replace/Add strings to file  
This should work interactively (ask the user what to fill in)
* [ ] Stop service
* [ ] Start service
* [ ] Ask user a question (Yes/No)
* [ ] Check that a file/folder is present
* [ ] Define which steps to repeat if a check was unsuccessfull
  * [ ] Validate that repeating the defined steps results in the same check that was unsuccessfull
* [ ] Read a file and check it for the occurens of string (i.e. errors in a log file)
* [ ] Wait for something (i.e. wait for the occurens of a specific string in a log file / wait for a service to start / when deleting the a `.war` file in the *webapps* folder of Tomcat, wait that the deployed distribution was delted)

### Additional Features

Thinks outside of the scope of the MVP that would add value.

* [ ] Allow for hot-reloading of the plan definition
* [ ] Execute SQL against database to check the result (Check)

## Libraries

* [go-prompt](https://github.com/c-bata/go-prompt)  
Make the shell interactive (Works with Powershell/Command Prompt)
* [survey](https://github.com/AlecAivazis/survey)  
Allow for user input (Works with Powershell/Command Prompt)
* [tview](https://github.com/rivo/tview)  
Alternative to *survey* (Works with Powershell/Command Prompt)
* [termui](https://github.com/gizak/termui)  
Charts, lists and other widgets (Works with Powershell/Command Prompt but doesn't look to good)
* [color](github.com/gookit/color)  
Allow printing to the console in color

## Build

The tool [goversioninfo](https://github.com/josephspurrier/goversioninfo) has to be installed in order to be able to build this project.

It should be enough to execute `go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo` to install it.

To build the project execute:

```shell
go generate
go build
```
