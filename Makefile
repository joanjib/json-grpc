main : main.go
	go run main.go
main.go : main_.go
	error_track $<
~/bin/error_track : utilities/error_track.go
	go build  -o ~/bin/error_track utilities/error_track.go
clean : 
	rm main.go
