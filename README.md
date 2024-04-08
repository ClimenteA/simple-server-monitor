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
SERVER_CPU_MAX_USAGE=90
SERVER_RAM_MAX_USAGE=90
SERVER_DISK_MAX_USAGE=90
SERVER_USAGE_INTERVAL_CHECK=1
```

### TODO

