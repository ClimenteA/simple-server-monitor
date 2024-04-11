# Is My Server Ok

Monitor your web apps with a simple chrome extension. Get notified on events you consider important (server resources CPU, RAM, Disk reached limit, new user signup, etc). 
Using [Go Fiber](https://gofiber.io/) and [Badger KV DB](https://dgraph.io/docs/badger/) to output maximum performance.  

## Browser extension UI

![](/pics/ismyserverok.png)


The UI is pretty simple with the following options:
- View Settings: view settings modal where you can add 1 or more servers to monitor;
- Pause notifications: you can turn off notifications (notifications are with sound);
- Clear events: delete all events saved in the browser's storage;
- Nuke all data: delete all events, settings, notifications as you would've just installed the chrome extension;
- Events table: all events displayed row by row (you can delete events one by one or use Clear events to delete all events);


## View settings modal

![](/pics/settings.png)

The View settings modal has 3 fields:
- Url: paste the base url of the Go api or your implementation;
- ApiKey: paste the ApiKey saved on the server's .env file;
- Request Interval (Minutes): at what interval in minutes should the chrome extension fetch events from the server;
- Settings table: you can delete settings one by one by clicking 'Delete' or you can edit one row by clicking 'Edit' and then clicking 'Save settings' button to save.


## Install browser extension

Clone the repo, go in your chromium based browser to `Settings > Manage Extensions > Load unpacked` and point to `chrome-extension folder` (you can stop here if you just want to know if your website is down or not).


## Running binary with a web app that does not use docker

Unless you run your server on a simple binary with an embeded database or using an external db via an uri. Not using docker and docker-compose makes your life harder than it needs to be. But, you can follow this [howtogeek tutorial](https://www.howtogeek.com/687970/how-to-run-a-linux-program-at-startup-with-systemd/) on how to run the server binary with systemd.


## Running binary with a web app that uses docker

You just need to add a new service in `docker-compose.yml` file next to your web application.

```yml

version: '3'

services:

  yourwebapp:
    etc
    
  ssm:
    build:
      context: .
      dockerfile: SSM.Dockerfile
    volumes:
      - ssmbadgerdata:/home/.badger
    env_file:
      - .env
    ports:
      - 4325:4325
    networks:
      - web

  otherservices:
    etc
    
networks:
  web:
    driver: bridge

volumes:
  ssmbadgerdata:

```

Here is the `SSM.Dockerfile` file that needs to be next to `docker-compose.yml`

```shell

FROM ubuntu:latest

WORKDIR /home

COPY .env .env

RUN apt-get update && apt-get install -y curl
RUN curl -L -o /home/server https://github.com/ClimenteA/simple-server-monitor/releases/download/v0.0.1/server
RUN chmod +x /home/server

CMD ["/home/server"]

```

Expected `.env` file next to binary:

```shell
SIMPLE_SERVER_MONITOR_PORT=4325
SIMPLE_SERVER_MONITOR_APIKEY=bdeef21a30cc0af802ac634ab2127817 # openssl rand -hex 16
SIMPLE_SERVER_MONITOR_CPU_MAX_USAGE=90
SIMPLE_SERVER_MONITOR_RAM_MAX_USAGE=90
SIMPLE_SERVER_MONITOR_DISK_MAX_USAGE=90
SIMPLE_SERVER_MONITOR_USAGE_INTERVAL_CHECK=5
SIMPLE_SERVER_MONITOR_HEALTH_URL=https://externally-acessible-url-no-auth
```

Maybe it looks like a lot of configuration, but don't worry it's just copy paste mostly. Checkout `docker` folder from this repo. With this setup you can send custom notifications to SIMPLE_SERVER_MONITOR server via RestApi provided.





# RestAPI

You can choose to use the Go binary setup or just recreate these routes in the webframework and programming language of your choice. 


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

