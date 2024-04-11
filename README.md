# Is My Server Ok

# Work in progress..

Monitor your web apps with a simple chrome extension. Get notified on events you consider important (server resources CPU, RAM, Disk reached limit, new user signup, etc). 
Using [Go Fiber](https://gofiber.io/) and [Badger KV DB](https://dgraph.io/docs/badger/) to output maximum performance.  


## Quickstart

- Clone the repo;
- Go in your chromium based browser to Settings > Manage Extensions > Load unpacked and point to `chrome-extension folder` (you can stop here if you just want to know if your website is down or not);
- Run `make build` to create the binary which will run on the server (or use the Binary or Dockerfile provided);
- Create an APIKEY with: `openssl rand -hex 16`;
- Place binary on your server;
- Expected `.env` file next to binary:

```shell
SIMPLE_SERVER_MONITOR_PORT=4325
SIMPLE_SERVER_MONITOR_APIKEY=bdeef21a30cc0af802ac634ab2127817
SIMPLE_SERVER_MONITOR_CPU_MAX_USAGE=90
SIMPLE_SERVER_MONITOR_RAM_MAX_USAGE=90
SIMPLE_SERVER_MONITOR_DISK_MAX_USAGE=90
SIMPLE_SERVER_MONITOR_USAGE_INTERVAL_CHECK=5
SIMPLE_SERVER_MONITOR_USAGE_HEALTH_URL=http://localhost:3000/health
```

**Browser extension main page:**
![](/pics/ismyserverok.png)


The UI is pretty simple with the following options:
- View Settings: view settings modal where you can add 1 or more servers to monitor;
- Pause notifications: you can turn off notifications (notifications are with sound);
- Clear events: delete all events saved in the browser's storage;
- Nuke all data: delete all events, settings, notifications as you would've just installed the chrome extension;
- Events table: all events displayed row by row (you can delete events one by one or use Clear events to delete all events);

**View settings modal:**
![](/pics/settings.png)

The View settings modal has 3 fields:
- Url: paste the base url of the Go api or your implementation;
- ApiKey: paste the ApiKey saved on the server's .env file;
- Request Interval (Minutes): at what interval in minutes should the chrome extension fetch events from the server;
- Settings table: you can delete settings one by one by clicking 'Delete' or you can edit one row by clicking 'Edit' and then clicking 'Save settings' button to save.


# RestAPI

You can choose to use the Go binary setup or just recreate these routes in the webframework and programming language of your choice. 

## Get events

Get all events saved in the Badger database. The events are deleted after they are sent because they will be saved in the chrome extension. 

Request (status code: `200` or `500`)::

```shell
curl  -X GET \
  'http://localhost:3000/simple-server-monitor/notifications' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'ApiKey: bdeef21a30cc0af802ac634ab2127817'
```

Response (status code: `200` or `500`):

```json
{
  "data": [
    {
      "EventId": "0febebf4-37e4-4a0b-8fec-a5e32dd645b5",
      "Title": "Test Title",
      "Message": "Test message",
      "Level": "info",
      "Timestamp": "20240408140055"
    }
  ]
}
```

## Save event

From your web app send this POST request to `Is My Server Ok`.

```shell
curl  -X POST \
  'http://localhost:3000/simple-server-monitor/save' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'ApiKey: bdeef21a30cc0af802ac634ab2127817' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "Title":  "Test Title",
	"Message":   "Test message",
	"Level":     "info"
}'
```

Response (status code: `201` or `500`):

```json
{"message": "yay or nay"}
```

## Delete event

You can also delete and event by it's EventId value.

```shell

curl  -X DELETE \
  'http://localhost:3000/simple-server-monitor/delete/{path parameter EventId}' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'ApiKey: bdeef21a30cc0af802ac634ab2127817'

```

Response (status code: `200` or `500`):

```json
{"message": "yay or nay"}
```


## Clear database

Clear database if you handled some errors and upgraded your server to get some fresh events.

Request:

```shell
curl  -X POST \
  'http://localhost:3000/simple-server-monitor/clear-database' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'ApiKey: bdeef21a30cc0af802ac634ab2127817'
```

Response (status code: `200` or `500`):

```json
{"message": "yay or nay"}
```

