This document tries to describe the work done, and the directory structure.

Done in Debian Buster system
Unit testing lib: Testify
For browsing the API from a web browser: grpcox

Directory structure:

db    							: go package to create a database object, and the its unit test. 
docker-entrypoint-initdb.d    	: directory mounted from docker compose to load a script that initialices the database
scripts    						: docker and postgresql scripts for: apropiate configuration of the postgresql in the transaction
								  issolation level, bash postgresql variables (development). Some docker usual commands in the dev day
								  by day.
constants    					: package containing constants of the server.
utils    						: package containing 3 useful functions used in varios places: 
									- InitClient: setup a client in go for connecting to the server. 
									- GenNumeric: Generates a String from the proto buffer Money type
									- GetMoney: generates a Money type (proto buffer) from string type.
								(no time for unit testing this package)
investor    					: package containing the investor CRUD implementation
financing    					: package containing the financing implementation (not finished)
issuer    						: package containing the issuer implementation
basic-client    				: in this directory we can find the rcp unit test like, for issuers, and investors.
server    						: the main.go is the server of the app. Also added a unit test too.
arexservices    				: package containing the gRPC implementation and the .proto file of the service
rdbms    						: various sql scripts for creating the database, the tables, and a PSQL FUNCTION where is implemented
								  the add_bid functionality.
financing-client   				: Initial test of the financing package, just a call to place a financing process (sell order + invoice)
								  Also some scripts in sql to add some data (sell orders, invoices, etc) to test at the SQL level 
								  the add_bid functionality. Add_bid only implemented in a SQL function and tested interatively 
								  from the psql utility ( no more time for more ).


How to use the docker compose:
------------------------------

In order to persist data in local hard drive please, create the "./pgsql-data" before starting the docker-compose

For starting the service on Linux machines:
sudo docker-compose up -d

The first time docker compose starts, it will create the tables and all necessary SQL objects. Before the first time the container has started, it will not re-create any more the sql objects.  Because add_bid.sql has to be used from psql utility is not loaded during this process.

For shutting down the docker compose app, please execute:

sudo docker-compose down

For browing and interacting with the server API from a browser do the following:

* Open a web browser and points: http://localhost:6969/
* Once the web app appear type in the "gRPC Server Target" :  app:50051

And ready to play with the arex api demo.
