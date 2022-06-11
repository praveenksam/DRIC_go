# A guide to using the service order downloader for DRIC

This is meant to be a guide to use the DRIC service order downloader to generate the service order every Sunday

## Identifying the script folder you need

The program to run the script depends on your operating system (windows, mac or linux), generally you should need the script in the x64 folder, however if you specifically know you have a ARM or M1 machine then use that script but the x64 script will also work for you.

For example:

1. If you have a Windows computer then you can use the files in the 'windows\x64' folder
2. If you have a Mac computer then you can use the files in the 'mac/x64' folder

## Configuration

In the script folder you use there is a config file which you can use to add and modify configuration settings.

### Necessary Configuration

Here is an overview of the configuration you might need to change:

- BoardId: The Board ID can be found with a script but if we leave the Trello Board "This Sunday" intact, this needs to be changed only if the board is deleted for any reason or if you need to connect to another board
- OutputLocation: This is the location where the PDF file should be saved, by default this is the current user's Desktop location, but if there is something mentioned here that location will be used instead

### Other Configuration

- Url: The Url endpoint where Trello will connect for the API, only needs to be changed if Trello ever changes this Url
- ApiKey: The ApiKey is provided by Trello to login, you can follow this tutorial to regenerate the API key if required: https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/. It is important to create this while logged in as the "Damascus Road IC" user account.
- Token: The token is the Trello Api login provided again by Trello based on your user account. You can follow this tutorial to regenerate the token if required: https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/

### Example layout of config.json

{
"Url": "https://api.trello.com",
"ApiKey": "123",
"Token": "abc123",
"BoardId": "456def",
"OutputLocation": "C:\Service order"
}

## How to run the script

It is best to run the script from Command Prompt/Powershell (Windows) or Terminal (Mac) as this will also allow you to add special options. In Windows you can also run the script by double clicking the script in the folder and this will use the default option.

### Initial setup

In Windows the script should generally be able to be run immediately either by double clicking (for default options) or from Powershell/Command Prompt

On a Mac or Linux you will have to first tell the system to allow the script to be run, to do this you will have to do some initial steps:

1. Open Terminal or any Terminal application on your computer
2. Go to the folder with the script for example if the script is in a folder ~/Documents/DRIC Service Download you will have to run the command `cd "~/Documents/DRIC\ Service\ Download"`
3. To give the script executing privileges run the command `chmod +x ./serviceListDownload`

### Running the script from Command Prompt/Powershell/Terminal

#### Windows

1. Open Powershell/Command Prompt
2. Go to the folder with the script for example if the script is in a folder C:\DRIC Service Download you will have to run the command `cd "C:\DRIC Service Download"`
3. Now you can run the script by calling the .exe file as follows `serviceListDownload.exe`
4. There are also some options you can select at this stage. To see an overview of the available options you can run `serviceListDownload.exe --help`. A further explanation of these options is given below
5. Every Sunday if you just want to use the default options and automatic sizing to make sure the Service List is in one page then you will only need to run the script with this command `serviceListDownload.exe`

#### Mac/Linux

1. Open Terminal or any Terminal application on your computer
2. Go to the folder with the script for example if the script is in a folder ~/Documents/DRIC Service Download you will have to run the command `cd "~/Documents/DRIC\ Service\ Download"`
3. Now you can run the script by calling the script as follows `./serviceListDownload`
4. There are also some options you can select at this stage. To see an overview of the available options you can run `.\serviceListDownload --help`. A further explanation of these options is given below
5. Every Sunday if you just want to use the default options and automatic sizing to make sure the Service List is in one page then you will only need to run the script with this command `.\serviceListDownload`

### Addtional options to generate Service List

The selected options are shown at the start of every attempt to download the service order

- startTime
  This allows to set the starting time of the Service as a 4 digit number and it is used to calculate the times for each part of the service, eg: 10:30 starting time should be entered as 1030. The default set is 1100 (default 1100)
  Usage: `.\serviceListDownload -startTime 0900`
- magFactor float
  This sets the magnification factor of text besides description. This will default to 1.6 (default 1.6)
  Usage: `.\serviceListDownload -magFactor 1.2`
- descMagFactor
  This allows to set the magnification factor for the description, this will only affect the description added under a card. This will default to 1.6 (default 1.6)
  Usage: `.\serviceListDownload -descMagFactor 1.2`
  -forceSize
  This allows the user to force the specified value for magFactor and descMagFactor to be used when making the PDF. In general the script will try to keep the service order to one page, however if this is used multiple pages are possible. This option is a boolean i.e. `true` or `false`, the default is `false`.
  Usage: `.\serviceListDownload -forceSize true`
  -listName string
  List name can be selected if there are multiple services, if not specified then the default value is Sunday Service (default "Sunday Service"). The list name should be placed in quotes and should be same as the name of the list on Trello. Assuming there is a list called "Sunday Service Two" the following usage would help.
  Usage: `.\serviceListDownload -listName "Sunday Service Two"`

These options can be combined if required. For example:

- If I have a Trello list called Monday Service
- it started at 18:00
- I wanted to make sure the list is developed with `magFactor` and `descMagFactor` set to 1.6 (the default value)
  then I can run this command: `.\serviceListDownload -forceSize true -listName "Monday Service" -startTime 1800`
