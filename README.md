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

  $ ./app --envId=envId --pollingInterval=2000 --bucketingPath=customDirectory

```

arguments:

| argument        | type   | description                                                                                                                                                   |
| --------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| envId           | string | Environment id provided by Flagship.                                                                                                                          |
| pollingInterval | int    | Define time interval between two bucketing updates. Default is 2000ms.                                                                                        |
| bucketingPath   | string | Directory path where bucketing file will be saved. <br/> Default path is : - `./flagship` for flagship-sync-agent                                             |
| config          | string | flagship configuration file path. **See [flagship configuration file](https://developers.flagship.io/docs/sdk/php/v2.0#flagship-configuration-file) section** |
