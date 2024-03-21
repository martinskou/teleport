# Teleport

With Teleport you can easily send and retrieve folders.

When you upload a folder, its automatically zipped and uploaded to a teleport server. After upload, you receive a download code. With this download code, you can download the folder from another machine. After download the folder is automatically unzipped and placed where you want it.

Teleport needs a server running somewhere. The server supports upload and download and a timeout after which any uploaded files are deleted.

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

It sound like FTP just with extra steps.


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

{
  "Server": "0.0.0.0",      
  "Port": 31345,
  "AuthToken": "1234",
  "TmpFolder": "tmp",
  "TimeOut": 3600
}

Server: IP or hostname of the Teleport server.
Port: Port the Teleport server runs on.
AuthToken: This token must be identical on the clients and the server.
TmpFolder: On the server its where the files are store.
TimeOut: Filen on the server are deleted after this interval (seconds)




## Usage

1. Start a server on a machine with a reachable IP 

A server needs to be running somewhere.

> teleport server

The IP of the server must be placed in config.json



2. Upload a folder to the server:

> teleport upload /development/misc

Folder /development/misc uploaded with retrieval code liyilak-0-qave



3. To retrieve the folder on same or different machine:

> teleport download liyilak-0-qave /somewhere/dev

Downloaded liyilak-0-qave (8831 bytes) to /somewhere/dev


## Server deployment

You can look inside the deploy folder of this project to find a script and systemd config file to run the server.

This script is hardcoded for transport to be placed in the folder /transport/

Example setup on Ubuntu server:

cd /
git clone https://github.com/martinskou/teleport.git
cd teleport




