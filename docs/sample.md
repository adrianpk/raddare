# Sample execution

## CURL request
```shell
curl -X GET http://localhost:8080/routes?src=13.388860,52.517037&dst=12.997634,51.909407&dst=13.999528,51.113219&dst=13.558559,53.253225&dst=13.428555,52.523219&dst=13.758559,52.893225
```

## Response
```markdown
getRoutesHandler:

Sorted:

[{Legs:[{Summary: Weight:958.7 Duration:717.1 Steps:[] Distance:4128.3}] WeightName:routability Weight:958.7 Duration:717.1 Distance:4128.3}]
[{Legs:[{Summary: Weight:6293.7 Duration:5594.6 Steps:[] Distance:104758}] WeightName:routability Weight:6293.7 Duration:5594.6 Distance:104758}]
[{Legs:[{Summary: Weight:6902.4 Duration:6332.8 Steps:[] Distance:121581.3}] WeightName:routability Weight:6902.4 Duration:6332.8 Distance:121581.3}]
[{Legs:[{Summary: Weight:14411.3 Duration:9954.2 Steps:[] Distance:207313.8}] WeightName:routability Weight:14411.3 Duration:9954.2 Distance:207313.8}]
```
