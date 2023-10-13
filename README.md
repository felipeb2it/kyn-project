# Small web service using gorilla 

This Servce Accept a json document with a single value to securely encrypt and decrypt.

# Installation

1) Build the Docker image:

```bash
docker build -t kyn-project .
```
2) Run your Go application inside a Docker container:

```bash
docker run -p 8080:8080 kyn-project
```

# REST API
The REST API to the example app is described below.

## Encrypt
The encrypt method accept a json object with a single value or a single line value.

### Request

`POST /api/encrypt`

    curl -X POST http://localhost:8080/api/encrypt -H "Content-Type: application/json" -d @data.json

### Response

    {"encrypted_value":"rdVDv9k68m2fFNfPorQOJsNfMchrP9RvkKreM10sVEWEU+TQleUU/45KoogJdVOAMXI="}

## Decrypt
The decrypt method accept a single line with the encrypted string returned from encrypt method.

`POST /api/decrypt`

    curl -X POST http://localhost:8080/api/decrypt -H "Content-Type: application/json" -d @datad.json

### Response

    {"decrypted_value":"Test"}


## License

[MIT](https://choosealicense.com/licenses/mit/)