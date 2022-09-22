# reapy

Make complex tasks repeatable by breaking them down into small automated steps.

## Concept

The idea is to break down complex tasks into their individual steps.

A typical installation of software might consist out of the following steps:

1. Download the software
2. Unzip the archive the software is distributed int
3. Copy the unzipped files to the target folder
4. Generate/Adjust the configuration
5. Install the software (i.e. as a windows service)

Such a process is defined inside of a YAML file. It consits out of the following pieces.

To break the process down we distinguish the following building blocks:

* Plan  
The overaching goal we want to achieve. (i.e. installing a piece of software). Each plan consists out of one ore more tasks.
* Task  
Each plan consists out of multiple tasks. A task could be the update of a database and checking that the update was successful. Each task therefore consists out of one or more steps.
* Step  
A small piece of work contributing to the fulfillment of a task or a check to ensure that a previous step was executed successfull. (i.e. unzipping a file or copying a folder from one place to another) (i.e. having the user manually confirm something or checking that a log file contains no errors)
