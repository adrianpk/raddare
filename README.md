# Raddare

A simple route planner service.

## Docker
```shell
$ docker run --rm adrianpksw/raddare:stage
```

## Clone
```shell
$ git clone https://github.com/adrianpk/raddare.git
```

## Run
```shell
$ make run
```

## API Call
```shell
$ make curl-routes
```

You should get something like this

```json
{
  "source": "13.388860,52.517037",
  "routes": [
    {
      "destination": "13.397631,52.529432",
      "duration": 433,
      "distance": 1999.6
    },
    {
      "destination": "13.428554,52.523239",
      "duration": 717.1,
      "distance": 4128.3
    }
  ]
}
```

## Status
**This is a draft document**
App and docs are under development.

More info to come...
