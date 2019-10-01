all: compile run

compile:
	go get
	GOOS=js GOARCH=wasm go build -o main.wasm

run:
	caddy

clean:
	rm main.wasm
