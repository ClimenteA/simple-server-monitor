# Simple server monitor

Monitor your web apps with a simple chrome extension. Get notified on events you consider important. 
Using [Go Fiber](https://gofiber.io/) and [Badger KV DB](https://dgraph.io/docs/badger/) to output maximum performance.


## Quickstart

- Clone the repo;
- Go in your chromium based browser to Settings > Manage Extensions > Load unpacked and point to `chrome-extension folder`;
- Run `make build` to create the binary which will run on the server (or use the Binary or Dockerfile provided);
- Create an APIKEY with: `openssl rand -hex 16`
- Place binary on your server;
- Expected `.env` file next to binary:

```shell
PORT=3000
APIKEY=bdeef21a30cc0af802ac634ab2127817
CPU_MAX_USAGE=90
RAM_MAX_USAGE=90
DISK_MAX_USAGE=90
USAGE_INTERVAL_CHECK=5
```


## Get events

Get all events saved in the Badger database.

Request (status code: `200` or `500`)::

```shell
curl  -X GET \
  'http://localhost:3000' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'ApiKey: bdeef21a30cc0af802ac634ab2127817'
```

Response (status code: `201` or `500`):

```json
{
  "data": [
    {
      "Title": "Test Title",
      "Message": "Test message",
      "Level": "info",
      "Timestamp": "20240408140055"
    }
  ]
}
```

## Clear database

Clear database if you handled some errors, upgraded your server to get some fresh events.

Request:

```shell
curl  -X POST \
  'http://localhost:3000/clear-database' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'ApiKey: bdeef21a30cc0af802ac634ab2127817'
```

Response (status code: `200` or `500`):

```json
{"message": "yay or nay"}
```

