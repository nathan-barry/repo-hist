run: build
	./cmd/bin/main

build:
	clear && ~/tailwindcss -i views/input.css -o static/tailwind.css --minify && go build -o ./cmd/bin ./cmd/main.go 
