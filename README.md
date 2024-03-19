# go-image-opencv

Image processing with opencv

## Features
- Convert image  from PNG to JPEG.
- Resize image to specified dimensions.
- Compress image

## Requirements

- Docker
- Make

## Development

- Development using DevContainer.
- Navigate to `.devcontainer` directory and run from there.

## Run Locally

### Run service

```shell
make image
make run
```

### Run test

```shell
# Make sure you have run `Run Service` above
make cli
make test
```

> Type `exit` and enter to quit from shell.

### Stop service

```shell
make stop
```

### API Documentation (OpenAPI 3.1)

Please refer to `/api` directory.  
Preview online: [api doc](https://elements-demo.stoplight.io/?spec=https://raw.githubusercontent.com/alfi-repo/go-image-opencv/main/api/api.yaml)

### Technologies

- GoCV 4.8.1 (via docker image)
- Go 1.21