# Should build for raspbery pi 2 and 3 (model b)
# Not tested on anything other than RPi3 running raspian-stretch
build:
	GOARM=7 GOARCH=arm GOOS=linux go build
