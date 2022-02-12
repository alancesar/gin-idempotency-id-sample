# Using a idempotent middleware in gin-gonic

## Starting the application

```shell
go mod download
go run main.go
```

## Request
### Create user
```shell
curl --location --request POST 'http://localhost:8080/user' \
--header 'Idempotency-Id: some-id' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "string",
    "email": "string"
}'
```
### Get user
```shell
curl --location --request GET 'http://localhost:8080/user/:id'
```

## Behaviour
The `middleware.Idempotency` will intercept all `POST` requests 
and with `Idempotency-Id` header present. If the return is not
an error (status code `4XX` and `5XX`), the  response body and
headers are going to be cached using the `idempotencyID` and 
request URL as key. If the `POST` request is received again with
both `Idempotency-Id` header and URL, the cached response will be
returned.