build:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' goroutinecreator.go

kill-process:
	sudo kill -9 $(shell pgrep goroutine)

clean:
	rm -rf goroutinecreator