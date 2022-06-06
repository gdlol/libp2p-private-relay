# libp2p-private-relay
A libp2p private relay node. It will provide unlimited relay service for whitelisted peerIds and IP addresses.

## Configuration (`config.json`)
- **`listenAddrStrings`**: List of libp2p multiaddress strings to listen on. By default it listens on port 55555.
- **`whitelistPeers`**: List of whitelisted peerIds.
- **`whitelistAddrs`**: List of whitelisted IP addresses.

## Build (export DOCKER_BUILDKIT=1)
```
docker compose build
```

## Run
```
docker compose up --detach
```
