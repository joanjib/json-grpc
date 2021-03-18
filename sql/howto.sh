sudo docker start 49882a1f0fe2
# for dropping the database enter this command to log in the postgresql instance.
psql -U postgres -d ""
# creating the datbase num
psql -U postgres -d "num" -f create-db.sql 
# finally, creating the objects of the database:
psql -d num -f tables.sql 
# backup of the database num:
sudo docker exec -i  postgres-arex  pg_dump --username postgres num > bk
