# Using Zigbee Home Container

## Building the Docker Image

To build the Docker image, use the following command:
```sh
docker buildx build -t ncs-builder .
```

### For Mac M1 Users

If you are using a Mac with an M1 chip, use this command to build the image for the appropriate platform:
```sh
docker buildx build --platform=linux/amd64 -t ncs-builder .
```

### Custom Toolchain or Zigbee Home Versions

If you need to specify custom toolchain or Zigbee Home versions, use the following command:
```sh
docker buildx build --platform=linux/amd64 --build-arg toolchain_version=v2.6.1 --build-arg sdk_nrf_branch=v2.6.1 -t ncs-builder:v2.6.1 --build-arg zigbee_home_version=develop .
```

## Building an Example Using the Docker Image

1. Navigate to the examples folder:
    ```sh
    cd examples
    ```

2. Run the Docker container:
    ```sh
    docker run --rm -v .:/workdir/examples -ti --platform linux/amd64 ncs-builder
    ```

3. The container will start an interactive shell. Use it to navigate to the desired example and then run:
    ```sh
    zigbee_home firmware build
    ```
