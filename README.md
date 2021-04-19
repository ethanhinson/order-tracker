# Delivery Tracker PoC

This application provides a simple microservice that estimates the ability of a given driver to make a delivery on time.

## Running the service

We use Docker Compose to facilitate running the application. Before running the service the first time, you need to build
the image first. There is a script provided in `./bin/build.sh` if you'd like to use it.

### Using the binary locally

 - `docker build . -t order-tracker:latest` to build the image
 - `docker-compose up` to bring the service up
 
### Development running the source code
 - `docker-compose up -f docker-compose-dev.yml` to bring the service up
 
## Accessing the app container

From the root of the project simply run: `docker-compose exec tracker bash`

## Deploying

The application is deploy to GKE. The `tracker` app/service is a managed workload in GKE. There is a private github
repository that GKE accesses for the `Dockerfile`. I have not created a CI process for this demo. 