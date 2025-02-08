# MQTT Broker
First of all clone the [mqtt-server](https://github.com/SoroushSaDev/mqtt-server) project from Github and then proceed with the installation 

**Installation tutorial is only for Microsoft Windows & Ubuntu Linux**

***In all senarios below, use "api.go" file (instead of "main.go") to build the API version, "main.go" file is for the Static (default) version***

## Windows Installation
To install and run the server on a Windows machine, you need to install **[Go Lang](https://go.dev/doc/install)** first.

Then open **CMD** in the directory below in the root directory of the project :
> /examples/auth/basic

and run the script below :
```bash
go build main.go
```
After the execution is done a new file named "main.exe" will appear next to the "main.go" in the current directory. run the exe file to start the server.

*You can set the exe file in windows startup to run the server on reboots automatically*
## Ubuntu Installation
To install and run the server on a Ubuntu machine, first install **Go Lang** with root privileges using the command below :
```bash
sudo apt install golang
```
Then cd to the directory below in the root directory of the project :
> /examples/auth/basic

and run this command to build the server from the source :
```bash
go build -o server main.go
```
a bash file named "server" will be created next to the "main.go" file in the current directory which running it starts the server

*You can set a supervisor config file to run the server on reboots automatically*

### Notes
- Building the server for the first time might take longer due to the installation of required packages & libraries (Network connection required)
- If you wanna use the API version, set the value of the "apiURL" variable in the ".env" file to retreive users & ACL data from your desired endpoint, by default it is set to the local address (127.0.0.1)
- If you wanna use the Static version, you can set your users info in the "authRules" variable
