![Authfox project logo](.github/media/authfox.svg)

---
Authfox is a simple session and authentication server that poweres and secures the Puroto social media. The basic design concepts are simplicity and security. It's one part of the Puroto stack and of the Puroto backend API.

# Contribute
We highly appreciate all kinds of contributions to Berry. For learning more about contributing to berry in an open source matter, please review our [contribution guidelines]() (coming soon!) to get started quickly.

# Compiling locally
Authfox builds up on some third party services and will refuse to run when they are not present. To make the local deployment easier, we have provided you with a simple docker-compose configuration. Since the docker-comtainers are only running on linux, you might need to set up other tools like WSL (Windows subsystem for Linux) for preceeding. Additionally a manual setup can be done as well, the services you'll need are:
- PostgreSQL
- Redis
- webserver for setting the CORS headers (optional)

## Prepare using docker-compose
Docker compose will configure and start all services needed, along side a few additional services that might come in handy in some cases. The following services will be started:
- Postgres (port 5432)
    - user: `user`, password: `pass`
- PGAdmin (port 80)
    - email: `user@test.lcl`, password: `pass`
- Redis (port 6379)
    - the password can be left blank
- Authfox (port 3622)
    - the local authfox service running for reference
    - swagger page: `http://localhost:3622` (you'll be redirected automatically)
- Caddy (port 3623)
    - A proxy that sets the headers used by CORS to allow local testing with web front-ends

### Running the docker-compose
```bash
# enter the deploy directory
cd deploy/

# start the docker-compose setup
docker-compose up --build
```

### Stop the docker-compose
```bash
docker-compose down
```
### Update containers
```bash
# download the new containers
docker-compose pull
# start the compose file again,
# docker-compose will re-start containers if needed
docker-compose up --build
```

## compiling
Authfox can be compiled using the official go compiler. But first, we need to give authfox the database configuration. This is done via environment variables like the following:
```bash
export POSTGRES_HOST=localhost
export POSTGRES_USER=user
export POSTGRES_PASS=pass
export POSTGRES_DB=authfox
export POSTGRES_PORT=5432
export POSTGRES_SSLMODE=disable
export POSTGRES_TIMEZONE=Europe/Berlin
export REDIS_HOST=localhost:6379
# optional, can be kept blank
export REDIS_PASS=""
```
> Note that if you don't use the docker-compose setup, you have to change the values to the ones of your own databases.

When that's done, all needed to do now is just:
```bash
go build .
```
> Please make sure that your terminal is in the root directory of this repository.

This will download the required go dependencies and should start authfox on the port `3621` for testing.

# License
Authfox is a project created and maintained by [Puroto](https://puroto.net) under the GPLv3.
