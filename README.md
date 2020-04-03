# IP Cam Email Api

This was designed for amcrest ip cameras (but would work for other brands) to receive a motion alert email and trigger a http request for integration with a home automation system, I use this specifically for motion detection.


## Configuration

See sample config, action_url is the value that will be called when an email is received.

#### Replacements

%s in the action_url is replaced with the value before the @ symbol in the To: address in the email, for my use case I use this as device id for that camera

##### Example

```
url = "http://endpoint.example/%s/motion"
To: 123@system.example

Becomes:

url = "http://endpoint.example/123/motion"

```

## Install

## Linux

Copy the binary to /usr/local/bin, use the service file to start the program on linux

## Windows

I have created a windows binary but it is untested, the system should look for a config.json in the same directory as where the binary is

## TLS config

Default the program starts on port 25, if you need tls/port 587 you want to enable tls in the config (see example). You'll need to generate a server.crt and server.key and place in the same directory as the binary.


Below are openssl commands to generate the certs

```
openssl genrsa -out server.key 2048
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
