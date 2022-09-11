# Grafagram

Simple webhook to publish Grafana alerts in a Telegram Channel.

### Build

```bash
docker build -t grafagram .
```

### Run

```bash

$BOT_TOKEN: Create a bot using botfather and get the token.
$CHAT_ID: ID of the chat (group, channel) to publish the incoming events into it.

docker run --restart=unless-stopped -e BOT_TOKEN=$BOT_TOKEN -e CHAT_ID=$CHAT_ID -p 1323:1323 --name=grafagram -d grafagram
```