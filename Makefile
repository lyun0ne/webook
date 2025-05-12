.PHONY: docker
docker:
	@rm webook || true
	# @GOOS=linux GOARCH=arm64 go build -o webook .
	@GOOS=linux GOARCH=arm64 go build -tags=k8s -o webook .
	@docker rmi -f lyunone/webook:v0.0.1 
	@docker build -t lyunone/webook:v0.0.1 .