# Simple server monitor

Monitor your web apps with a simple chrome extension. Get notified on events you consider important. 


# TODO

## Quickstart

- Clone the repo;
- Go in your chromium based browser to Settings > Manage Extensions > Load unpacked and point to `chrome-extension folder`;
- Run `make build` to create the binary which will run on the server (or use the Binary or Dockerfile provided);
- manually generate keys with: `openssl rand -base64 32`
- https://dgraph.io/docs/badger/get-started/


// CPU Usage
// cat /proc/stat |grep cpu |tail -1|awk '{print ($5*100)/($2+$3+$4+$5+$6+$7+$8+$9+$10)}'|awk '{print 100-$1"%"}'

// RAM Usage
// free -h | awk '/^Mem:/ {print ($3/$2)*100"%"}'

// Disc Usage
// df -h | awk '$6 == "/" {print $5}'
