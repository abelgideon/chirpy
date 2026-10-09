# Chirpy

## Prerequisites
- go
- postgres
- .env file with the following key value pairs:

  ```
    DB_URL="your_database_url"
    PLATFROM="dev_or_prod"
    JWT_SECRET="your_jwt_secret"
    POLKA_KEY="your_webhook_key"
  ```

## Installation
`go install github.com/abelgideon/chirpy`

## Endpoints
- `GET /app/`: access file server

- `GET /api/healthz`: check if server is up

- `GET /admin/metrics`: hit count for file server

- `POST /admin/reset`: delete all users

- `POST /api/users`: create user

- `POST /api/chirps`: create chirp

- `GET /api/chirps`: get all chirps with sort and author_id filtering

- `GET /api/chirps/{chirpID}`: get single chirp

- `POST /api/login`: user login

- `POST /api/refresh`: refresh access token

- `POST /api/revoke`: revoke refresh token

- `PUT /api/users`: update user

- `DELETE /api/chirps/{chirpID}`: delete chirp

- `POST /api/polka/webhooks`: webhook handler
