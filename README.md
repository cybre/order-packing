# Order Packing

This is a solution to a take-home assignment from RE Partners.
The app demo is publicly available at https://order-packing.fly.dev

## Running Locally

The solution is a multi-container application.  
docker-compose is used for running the containers.

To run the project locally, run:
```bash
make run
```

Now you should be able to access the web app locally at `http://localhost:3001`.

You can stop the containers by running:
```bash
make stop
```

## Tests
The test suite is containerized and can be executed by running:
```bash
make test
```

This will build a new Docker image using Dockerfile.test and run the test suite. 