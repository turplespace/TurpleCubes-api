# Turple Cube API Documentation

## Running Turple Cube

### Prerequisites
- Install [Docker](https://docs.docker.com/get-docker/)
- Ensure Docker daemon is running

### Build the Docker Image (Optional, if not using pre-built image)
```bash
docker build -t turplecubes .
```

### Run the API Server
#### Linux
```bash
sudo docker run -d --user root \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd)/docker_volume:/app/bin/ \
  -p 8080:8080 \
  --privileged \
  sanjaysagar12/turplecubes
```

### Run Nginx Proxy
```bash
sudo docker rm -f turplecubes-proxy && \
sudo docker run -d --name turplecubes-proxy \
  -v $(pwd)/docker_volume/turplecubes_proxy:/etc/nginx/conf.d \
  -p 80:80 \
  nginx
```

