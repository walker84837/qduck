run:
	go run src/main.go

build:
	GOOS=windows go build -o dist/qduck-windows.exe src/main.go
	GOOS=linux go build -o dist/qduck-linux src/main.go
	GOOS=darwin go build -o dist/qduck-mac.dmg src/main.go
