main : main.go
	go run main.go
main.go : main_.go
	error_track $<
clean : 
	rm main.go
