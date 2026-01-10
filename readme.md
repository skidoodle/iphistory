# IPHistory

A simple tool for tracking your network's public IP address. It periodically checks for changes and logs them to a local database, providing a clean, searchable web interface to browse your history. It's designed to be lightweight, self-reliant, and fast without requiring any maintenance.

## Deploy

```yaml
services:
  iphistory:
    image: ghcr.io/skidoodle/iphistory
    container_name: iphistory
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - data:/app

volumes:
  data:
```

## License

[GPL-3.0](https://github.com/skidoodle/iphistory/blob/main/license)
