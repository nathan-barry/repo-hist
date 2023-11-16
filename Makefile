run: build
	./bin/main

build:
	clear && ~/tailwindcss -i views/input.css -o static/tailwind.css --minify && go build -o ./bin main.go 
