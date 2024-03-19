# Build container image.
build:
	docker build -t go-image-opencv .

# Step into the container shell.
cli:
	docker exec -it goimageopencv sh

# Run the container.
run:
	docker run -p 3000:3000 -d --name goimageopencv go-image-opencv

# Stop and remove the container.
stop:
	docker stop goimageopencv && docker rm goimageopencv

# Run test and coverage report.
test:
	go test -cover

# Run test and generate detailed coverage report.
testreport:
	go test -coverprofile cover.out
	go tool cover -html=cover.out -o coverage.html