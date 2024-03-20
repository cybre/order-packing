run:
	docker-compose up --build -d
stop:
	docker-compose down
test:
	docker build -t order-packing-tests --file Dockerfile.test . && docker run -it --rm order-packing-tests