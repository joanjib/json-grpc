# removing all the contents of all the tables:
psql -U joan -d num -f empty.sql
go test -v client_test.go
