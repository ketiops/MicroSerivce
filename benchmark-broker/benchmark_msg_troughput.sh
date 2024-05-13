#/bin/bash
  
TOPIC_NUMBER=$((RANDOM % 100+ 1050))
docker stop go-kafka-app
docker rm go-kafka-app
docker run \
  -e DB_HOST='10.0.1.104' \
  -e DB_PORT='32000' \
  -e KAFKA_SERVERS='10.0.1.110:31001' \
  -e TABLE_NUMBER='1' \
  -e TOPIC_NUMBER="$TOPIC_NUMBER" \
  --name go-kafka-app \
  eca560bb5be1
