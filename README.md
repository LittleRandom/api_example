# API Example
Example API skeleton built with Go for learning purposes.

This repository was built by me to explore how to create an API using go. It uses GORM as the database interface, sqlite + driver for lightweight sql database, and serves a static version of Swagger API for documentation.

# Demo

Try out the project [demo](https://api.plainrandom.com/).

# Features

Some features I put in the repository are:

    - CI/CD with github actions to build, publish, and deploy.
    - Containerized multi-stage build to reduce overhead.
    - Router that serves SPA from root and API from /items endpoint.
    - Integration tests to ensure handlers work properly.
    - Modularized packages to be a starting point for any API.

# Quick Setup

1. Install Docker and Docker-Compose

- [Docker Install documentation](https://docs.docker.com/install/)
- [Docker-Compose Install documentation](https://docs.docker.com/compose/install/)

2. Create a docker-compose.yml file similar to this:

```yml
services:
  api:
    image: ghcr.io/littlerandom/api_example:main
    restart: unless-stopped
    ports:
      - 80:5050
networks: {}
```

This is the bare minimum configuration required to host the API.

3. Bring up stack by running

```bash
docker compose up -d
```

4. Open up a browser to the local address.

[http://127.0.0.1:81](http://127.0.0.1:81)

# How I hosted the application

I'm currently self-hosting the application because I'm too cheap to rent a VPS. My networking stack involves:

1. Reverse-proxy by Cloudflare CDN.
2. [Nginx Proxy Manager](https://nginxproxymanager.com/) to proxy it throughout my LAN.
3. Debian LXC with Docker compose to host the application.
4. Github actions-runner for CI/CD automation for build, publish, and deploy.