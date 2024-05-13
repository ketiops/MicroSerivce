Start
-------------


```
$ cd benchmark-broker && docker build -t kafka-client:0.1 .

$ docker run \
  -e DB_HOST='localhost' \
  -e DB_PORT='Port' \
  -e KAFKA_SERVERS='10.0.1.104:31001' \
  -e TABLE_NUMBER='1' \
  -e TOPIC_NUMBER='990' \
  --name go-kafka-app1 \
  kafka-client:0.1 && docker exec -it go-kafka-app /bin/bash
```
