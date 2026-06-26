CC := "go"
MAIN := "main.go"

all:
  {{CC}} run {{MAIN}} --no-welcome

[group('build')]
build:
  mkdir -p bin
  {{CC}} build -o ./bin
  
[group('container')]
podman:
  podman build -t qasimwarraich/cli-tutor .

[group('container')]
podmanrun:
  podman run -it qasimwarraich/cli-tutor

[group('container')]
podmanpull:
  podman pull ghcr.io/qasimwarraich/cli-tutor:latest

[group('install')]
install: build
  cp ./bin/cli-tutor /usr/local/bin

[group('install')]
uninstall:
  rm -f /usr/local/bin/cli-tutor

[group('build')]
clean:
  rm -rf bin

[group('install')]
goinstall:
  {{CC}} install

[group('install')]
gouninstall:
  rm ~/go/bin/cli-tutor

[group('test')]
test:
  gotest -v ./pkg/...
