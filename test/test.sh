docker image rm -f test
docker rm -f chatApp_test
docker build -t test:latest -f Dockerfile .
docker run --name chatApp_test --network chatnet test:latest
