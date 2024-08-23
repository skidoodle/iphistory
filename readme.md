# IPHistory

The IPHistory project is a simple yet effective solution for tracking and logging the public IP address of your network. It periodically fetches the public IP address and logs it to a file, while also providing a web interface to view the IP history in a clean UI and an endpoint (`/history`) for JSON format.

![iphistory](https://github.com/user-attachments/assets/4c4ea170-94ef-42fc-b3b5-5db735789855)

## Running Locally

### With Docker

```sh
git clone https://github.com/skidoodle/iphistory
cd iphistory
docker build -t iphistory:main .
docker run -p 8080:8080 iphistory:main
```

### Without Docker

```sh
git clone https://github.com/skidoodle/iphistory
cd iphistory
go run main.go
```

## Deploying

### Docker Compose

```yaml
version: '3.9'

services:
  iphistory:
    image: ghcr.io/skidoodle/iphistory:main
    container_name: iphistory
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - iphistory_data:/app

volumes:
  iphistory_data:
    external: false
```

### Docker Run

```sh
docker run \
  -d \
  --name=iphistory \
  --restart=unless-stopped \
  -p 8080:8080 \
  ghcr.io/skidoodle/iphistory:main
```

## License

[GPL-3.0](https://github.com/skidoodle/iphistory/blob/main/license)
