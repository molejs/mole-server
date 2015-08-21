# mole-server [![Build Status](https://travis-ci.org/foo/mole-server.svg?branch=master)](https://travis-ci.org/foo/mole-server)

Mole server is one of the three independent packages that form the whole Mole service. It is a server with one endpoint to report logs and to retrieve them from the database.

All the reported logs must follow the **mole log specification**.

## Install

```
go get https://github.com/foo/mole-server
```

## Configure and run

Configuration of mole server is done via environment variables.

Name              | Description                    | Default value
-----------------------------------------------------------------------
`MOLE_ADDR`       | Server address, e.g `:8080`    | `:8080`
`MOLE_MONGO_ADDR` | MongoDB server address         | `127.0.0.1:27017`
`MOLE_DB_NAME`    | MongoDB database name for mole | `127.0.0.1:27017`
`MOLE_KEY`        | SSL certificate key            | none [1]
`MOLE_CERT`       | SSL certificate                | none [1]
-----------------------------------------------------------------------

[1] If cert and key are not empty the server will be started as an HTTPS server.

To run the server just:
```
go run server.go
```

## Retrieving logs
```
curl http://example.com/logs?limit=20&skip=0
```
**Query string params**

Param name | Description                       | Required | Default value
--------------------------------------------------------------------------
`limit`    | Max number of logs to retrieve    | no       | `25`
`skip`     | Skip first n logs when retrieving | no       | `0`
--------------------------------------------------------------------------

All results are ordered by descending creation date.
A successful result looks like:
```
{
  "error": false,
  "logs": [array of logs]
  "count": 10, // The number of logs retrieved
  "total": 50 // Total number of logs
}
```

## Reporting logs
```
curl -X POST --data "{the json log object}" http://example.com/logs
```

If the request is successful the request will look like:
```
{
  "error": false
}
```

## Errors

All error responses look like:
```
{
  "error": true,
  "msg": "error message"
}
```
