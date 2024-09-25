build:
	GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	chmod +x bootstrap

deploy: build
	serverless deploy --stage prod

clean:
	rm -f bootstrap