# Teleport

With Teleport you can easily send and retrieve folders.

When you upload a folder, its automatically zipped and uploaded to a teleport server. The upload is either named after the folder or supplying a name to the upload command. With this name, you can download the folder from another machine. After download the folder is automatically unzipped and placed where you want it (download argument).

Teleport needs a server running somewhere. The server supports upload and download and a timeout after which any uploaded files are deleted.

Teleport uses regular HTTP to upload and download files.

 ```
    @@@@@@@  @@@@@@@@  @@@       @@@@@@@@  @@@@@@@    @@@@@@   @@@@@@@   @@@@@@@ 
    @@@@@@@  @@@@@@@@  @@@       @@@@@@@@  @@@@@@@@  @@@@@@@@  @@@@@@@@  @@@@@@@ 
      @@!    @@!       @@!       @@!       @@!  @@@  @@!  @@@  @@!  @@@    @@!   
      !@!    !@!       !@!       !@!       !@!  @!@  !@!  @!@  !@!  @!@    !@!   
      @!!    @!!!:!    @!!       @!!!:!    @!@@!@!   @!@  !@!  @!@!!@!     @!!   
      !!!    !!!!!:    !!!       !!!!!:    !!@!!!    !@!  !!!  !!@!@!      !!!   
      !!:    !!:       !!:       !!:       !!:       !!:  !!!  !!: :!!     !!:   
      :!:    :!:        :!:      :!:       :!:       :!:  !:!  :!:  !:!    :!:   
       ::     :: ::::   :: ::::   :: ::::   ::       ::::: ::  ::   :::     ::   
       :     : :: ::   : :: : :  : :: ::    :         : :  :    :   : :     :    
 ```

*It sounds like FTP just with extra steps.*


## Example

On home computer, upload a folder:

> teleport upload /docs/report23

On work computer, download the folder report23 and place it in documents:

> teleport download report23 /documents


## Why?

I work on both mac and linux and often need to transfer a folder with a project or documents from the one to the other.


## Install

Clone repro:

> git clone https://github.com/martinskou/teleport.git

Build teleport

> go build .

Build a config file

> ./teleport

Configure the server and authtoken in the config.json file.

Repeat install on a server somewhere online.


## Config

There are a config.json file:

```
 {
  "Server": "0.0.0.0",      
  "Port": 31345,
  "AuthToken": "1234",
  "TmpFolder": "tmp",
  "TimeOut": 3600
  "UseTLS": true,
  "TLSCert": "",
  "TLSKey": "",
  "AllowSelfSigned": true
}
```
 
Server: IP or hostname of the Teleport server.
Port: Port the Teleport server runs on.
AuthToken: This token must be identical on the clients and the server.
TmpFolder: On the server its where the files are store.
TimeOut: Filen on the server are deleted after this interval (seconds)
UseTLS: If true only transmit with TLS/SSL cert
TLSCert: Certificate file
TLSKey: Key file
AllowSelfSigned: If true and TLSCert and TLSKey are empty, a selfsigned cert is used


## Usage

1. Start a server on a machine with a reachable IP 

A server needs to be running somewhere.

> teleport server

The IP of the server must be placed in config.json



2. Upload a folder to the server:

> teleport upload /development/misc

Folder /development/misc uploaded with retrieval code misc


3. List files on server:

> teleport list

Will print a list of all files on the server.


4. To retrieve the folder on same or different machine:

> teleport download misc /somewhere/dev

Downloaded misc (8831 bytes) to /somewhere/dev


## Server deployment

You can look inside the deploy folder of this project to find a script and systemd config file to run the server.

This script is hardcoded for transport to be placed in the folder /transport/

Example setup on Ubuntu server:

> cd /
> git clone https://github.com/martinskou/teleport.git
> cd teleport
> cd deploy
> source deploy.sh
> cd ..

Edit config.json

> systemctl start teleport.service


## Security

TLS/SSL can be configured in the config.json file.


## Tests

No tests was written during the production.


## Requirements

Go 1.22 is required.

https://go.dev/

Working on Linux and macOs, not tested on Windows.

Uses https://github.com/urfave/cli for handleling the command line arguments.


## Alternatives

Dropbox, etc. if you are into that.

2 bash scripts. One that zips and uploads to an existing FTP server. And another the downloads and unzip.

Rsync to/from a server is also an option.

But where is the fun in that?