# Delivery Tracker PoC

This application provides a simple microservice that estimates the ability of a given driver to make a delivery on time.

## Running the service

We use Docker Compose to facilitate running the application. Before running the service the first time, you need to build
the image first. There is a script provided in `./bin/build.sh` if you'd like to use it.

 - `docker build . -t order-tracker:latest` to build the image
 - `docker-compose up` to bring the service up
 
## Accessing the containers

From the root of the project simply run: `docker-compose exec tracker bash`