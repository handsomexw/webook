.PHONY: docker
docker:
	@rm webook || true
	@docker rmi -f webook:v0.0.1
	@docker build -t webook:v0.0.1 .