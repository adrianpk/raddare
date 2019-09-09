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

```shell
getRoutesHandler:

Sorted:

[{Legs:[{Summary: Weight:958.7 Duration:717.1 Steps:[] Distance:4128.3}] WeightName:routability Weight:958.7 Duration:717.1 Distance:4128.3}]
[{Legs:[{Summary: Weight:6293.7 Duration:5594.6 Steps:[] Distance:104758}] WeightName:routability Weight:6293.7 Duration:5594.6 Distance:104758}]
[{Legs:[{Summary: Weight:6902.4 Duration:6332.8 Steps:[] Distance:121581.3}] WeightName:routability Weight:6902.4 Duration:6332.8 Distance:121581.3}]
[{Legs:[{Summary: Weight:14411.3 Duration:9954.2 Steps:[] Distance:207313.8}] WeightName:routability Weight:14411.3 Duration:9954.2 Distance:207313.8}]
```

**Note**
This is not yet the output required by the specification. Work in progress.

Work in progress.

## Status
**This is a draft document**
App and docs are under development.

More info to come...