build:
	docker build -t http_helloworld .

run:
	docker run --rm -p 8080:8080 -e PORT=8080 -v ${PWD}/custom_responses:/app/custom_responses http_helloworld