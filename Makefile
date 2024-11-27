gomod:
	@go mod tidy

gobin: gomod
	GOOS=linux GOARCH=amd64 go build -o perftest

image: gobin
	@docker build -t sqlperftest:20241118 -f ./Dockerfile .
	@docker tag sqlperftest:20241118 quay.io/ybrillou/sqlperftest:20241118
	@docker push quay.io/ybrillou/sqlperftest:20241118
