This document tries to describe the work done, and the directory structure.

Done in Debian Buster system
Unit testing lib: Testify
For browsing the API from a web browser: grpcox
Editor : vim
Compilation directives in Makefiles


Directory structure:



How to use the docker compose:
------------------------------

In order to persist data in local hard drive please, create the "./pgsql-data" before starting the docker-compose

For starting the service on Linux machines:
sudo docker-compose up -d

The first time docker compose starts, it will create the tables and all necessary SQL objects. To warn the server for this initial creation execute the command in the root of the project "touch new". This command will create an empty file. If the server on start up detects this file will create the basic sql types of the APP and the special SQL function add_bid.

When it's not the first time the APP is started, remove the "new" file.

For shutting down the docker compose app, please execute:

sudo docker-compose down

For browing and interacting with the server API from a browser do the following:

* Open a web browser and points: http://localhost:6969/
* Once the web app appear type in the "gRPC Server Target" :  app:50051

And ready to play with the arex api demo.
