Made with:

- Debian Buster
- Unit testing lib: Testify
- For browsing the API from a web browser: grpcox
- Editor: vim
- ORM: gorm
- Database: Postgresql
- Compilation directives in Makefiles
- Docker compose file in the root directory

The following project is an improvement of the arex repository APP.

Two improvements are made into the "utilities" directory. 
- error_track: this utility tries to bring a kind of Exceptions to the Golang (check server/main_.go)
- macro_expansion: sometimes macros are very useful specially when your are dealing with generated code (gRPC). Thanks to this utility it's reduced the amount of code necessary to implement repetitive APIs. This utility brings a kind of macro system to golang.

It's also used a ORM : gorm, to reduce the amount of code. All entities are centralized in the models.go file. In the persiscence are used personalized domains and types (SQL), Gorm (Go) and SQL functions to do the most costly operarions into the database itself.

Directory structure:

- server: contains the main.go file needed to start the server.
- utils: contains an utility to connect a client to the server using gRPC.
- models: single go file with all models and cast functions to cast instances of gorm to instances of gRPC and viceversa.
- sql: all sql scrips. The main ones are:
  - add_bid.sql: processes a bid
  - types-domains.sql: basic sql types and domains used in the app
  - The rest are just examples and utilities
- arexservices: the definition of the gRPC  API in .proto file and the generated .go files. Contains an script to generate the go files from the .proto file.
- client: single file with all test to the server made from the client side (gRPC)
- db: configuration and creation of the gorm db object.
 
 
File nomenclature:
- Files ended with "_.go" are files that need to be processed with the utility error_track because code contains exceptions.
- Files ended with "__go" are files that contains macro expansions and need to be processed with the macro_expansion utilitity. 
- Files ended with ".go" are the actual code.
The final code is always .go file, and Makefiles into each directoty make easier this process of code generation and transformation.



How to use the docker compose:
------------------------------

In order to persist data in local hard drive please, create the "./pgsql-data" before starting the docker-compose

For starting the service on Linux machines:
sudo docker-compose up -d

The first time docker compose starts, it will create the tables and all necessary SQL objects. To warn the server for this initial creation execute the command in the root of the project "touch new". This command will create an empty file named "new". If the server on start up detects this file will create the basic sql types of the APP and the special SQL function add_bid.

When it's not the first time the APP is started, remove the "new" file.

For shutting down the docker compose app, please execute:

sudo docker-compose down

For browing and interacting with the server API from a browser do the following:

* Open a web browser and points: http://localhost:6969/
* Once the web app appear type in the "gRPC Server Target" :  app:50051

And ready to play with the arex api demo.
