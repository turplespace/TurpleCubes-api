# Turple Cube API Documentation

REST API for managing Docker workspaces and containers.

## Endpoints

### Get Images
```bash
curl -X GET http://localhost:8080/api/images
```

### Get Workspaces
```bash
curl -X GET http://localhost:8080/api/workspaces
```

### Create Workspace
```bash
curl -X POST http://localhost:8080/api/workspace/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "myworkspace",
    "desc": "This is a description of my workspace"
  }'
```

### Edit Workspace
```bash
curl -X PUT http://localhost:8080/api/workspace/edit \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "updated_workspace_name",
    "desc": "Updated description"
  }'
```

### Delete Workspace
```bash
curl -X DELETE "http://localhost:8080/api/workspace/delete?id=1"
```

### Get Cubes
```bash
curl -X GET "http://localhost:8080/api/cubes?workspace_id=1"
```

### Add Cubes
```bash
curl -X POST http://localhost:8080/api/cube/add \
  -H "Content-Type: application/json" \
  -d '{
    "workspace_id": 1,
    "cubes": [{
      "name": "golang-test-server1",
      "image": "golang-portos",
      "ports": ["80:80"],
      "environment_vars": [
        "NGINX_PORT=80"
      ],
      "resource_limits": {
        "cpus": "1.0",
        "memory": "512M"
      },
      "volumes": {
        "golang-test-server1-v1": "/data"
      },
      "labels": ["env=prod"]
    }]
  }'
```

### Edit Cube
```bash
curl -X PUT http://localhost:8080/api/cube/edit \
  -H "Content-Type: application/json" \
  -d '{
    "cube_id": 1,
    "updated_cube": {
      "name": "updated_cube_name",
      "image": "updated_image",
      "ports": ["8080:80"],
      "environment_vars": [
        "UPDATED_ENV_VAR=value"
      ],
      "resource_limits": {
        "cpus": "0.5",
        "memory": "256M"
      },
      "volumes": {
        "updated_volume": "/updated_data"
      },
      "labels": ["env=prod"]
    }
  }'
```

### Delete Cube
```bash
curl -X DELETE "http://localhost:8080/api/cube/delete?cube_id=1"
```

### Get Cube Data
```bash
curl -X GET "http://localhost:8080/api/cube?cube_id=1"
```

### Building Docker
```bash
docker build -t portos .
```

### Running Images
#### linux
```bash
sudo docker run -d --user root \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd)/docker_volume:/app/bin/ \
  -p 8080:8080 \
  --privileged \
  turplecubes


sudo docker rm -f turplecubes-proxy && \
sudo docker run -d --name turplecubes-proxy \
  -v $(pwd)/docker_volume/turplecubes_proxy:/etc/nginx/conf.d \
  -p 80:80 \
  nginx
```

