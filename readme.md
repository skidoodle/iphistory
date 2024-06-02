# IPHistory

The IPHistory project is a simple yet effective solution for tracking and logging the public IP address of your network. It periodically fetches the public IP address and logs it to a file, while also providing a web interface to view the IP history in a clean UI and an endpoint (`/history`) for JSON format.

![iphistory](https://github.com/skidoodle/iphistory/assets/53189968/ca1b41fd-8295-4455-9b89-3f1b3c2b030e)

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
    container_name: iphistory
    image: 'ghcr.io/skidoodle/iphistory:main'
    restart: unless-stopped
    ports:
      - '8080:8080'
    volumes:
      - data:/app

volumes:
  data:
    driver: local
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
