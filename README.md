# journal-2-logentries 

Ship systemd journal entires to logentries.com over SSL.

## Usage

```
sudo LOGENTRIES_TOKEN=<token> journal-2-logentries
```

```
sudo docker run -d -e 'LOGENTRIES_TOKEN=<token>' -v /run/journald.sock:/run/journald.sock \
quay.io/octoblu/journal-2-logentries
```

## Configuration

All configuration is done through env vars.

* `LOGENTRIES_JOURNAL_SOCKET` - The systemd-journal-gatewayd socket. `/run/journald.sock`
* `LOGENTRIES_URL` - The log entry url. `data.logentries.com:433`
* `LOGENTRIES_TOKEN` - The logentries.com TCP token -- See https://logentries.com/doc/input-token

Note: Make sure that systemd-journal-gatewayd is actually listening on
`/run/journald.sock`. This is not done by default on CoreOS -- See
[example cloud-config](cloud-config.yaml)

## Building

```
GO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .
```

### Docker

```
docker build -t quay.io/octoblu/journal-2-logentries .
docker push quay.io/octoblu/journal-2-logentries
```

## Fleet integration

```
etcdctl set /logentries.com/token <token>
```

Edit `journal-2-logentries.service`

```
[Unit]
Description=Forward Systemd Journal to logentries.com

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker kill journal-2-logentries
ExecStartPre=-/usr/bin/docker rm journal-2-logentries
ExecStartPre=/usr/bin/docker pull quay.io/kelseyhightower/journal-2-logentries
ExecStart=/usr/bin/bash -c \
"/usr/bin/docker run --name journal-2-logentries \
-v /run/journald.sock:/run/journald.sock \
-e LOGENTRIES_TOKEN=`etcdctl get /logentries.com/token` \
quay.io/kelseyhightower/journal-2-logentries"

[X-Fleet]
Global=true
```

```
fleetctl start journal-2-logentries.service
```
