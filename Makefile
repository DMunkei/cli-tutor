CC = go
MAIN = main.go

all:
	${CC} run ${MAIN} --no-welcome

build:
	mkdir -p bin
	${CC} build -o ./bin

podman:
	podman build -t qasimwarraich/cli-tutor .

podmanrun:
	podman run -it qasimwarraich/cli-tutor

podmanpull:
	podman pull ghcr.io/qasimwarraich/cli-tutor:latest

install:
	mkdir -p bin
	${CC} build -o ./bin
	cp ./bin/cli-tutor /usr/local/bin

uninstall:
	rm -f /usr/local/bin/cli-tutor

clean:
	rm -rf bin

goinstall:
	${CC} install

gouninstall:
	rm ~/go/bin/cli-tutor

test:
	gotest -v ./pkg/...
