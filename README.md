[![CI](https://github.com/flagship-io/flagship-sync-agent/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/flagship-io/flagship-sync-agent/actions/workflows/ci.yml)

# flagship-sync-agent

flagship-sync-agent is a binary that performs the bucketing polling process. learn more [here](https://developers.flagship.io/docs/sdk/php/v2.0#bucketing-polling)

## Project setup

```bash

go mod download

```

## build

```bash

go build -o app .

```

## Test

```bash

go test  -coverprofile=coverage.txt -covermode=atomic ./... && go tool cover -html=coverage.txt -o cover.html

```

## Run

```bash

  $ ./app --envId=envId --pollingInterval=2000 --port=3000 --address=0.0.0.0

```

arguments:

| argument        | type   | description                                                            |
| --------------- | ------ | ---------------------------------------------------------------------- |
| envId           | string | Environment id provided by Flagship.                                   |
| pollingInterval | int    | Define time interval between two bucketing updates. Default is 2000ms. |
| Port            | int    | Endpoint listen port. Default is 8080                                  |
| address         | string | Address where the endpoint is served. Default is 0.0.0.0               |

## Docker

```bash

docker pull flagshipio/sync-agent

env FS_ENV_ID="YOUR_ENV_ID" docker run -p 3000:8080 -e FS_ENV_ID flagshipio/sync-agent

```

Environment variables:

| argument            | type   | description                                                            |
| ------------------- | ------ | ---------------------------------------------------------------------- |
| FS_ENV_ID           | string | Environment id provided by Flagship.                                   |
| FS_POLLING_INTERVAL | int    | Define time interval between two bucketing updates. Default is 2000ms. |
| FS_PORT             | int    | Endpoint listen port. Default is 8080                                  |
| FS_ADDRESS          | string | Address where the endpoint is served. Default is 0.0.0.0               |

## API docs

| route         | Description                                           |
| ------------- | ----------------------------------------------------- |
| /bucketing    | Get the Json bucketing file                           |
| /health_check | Return http code 200 to check if the sync agent is up |
