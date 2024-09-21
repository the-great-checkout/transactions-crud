# transactions-crud
Transactions SSOT

## Swagger commands
To update swagger:
```shell
docker run --rm -v $(pwd):/code ghcr.io/swaggo/swag:latest fmt
docker run --rm -v $(pwd):/code ghcr.io/swaggo/swag:latest init
```

## Kafka commands
To develop with Kafka, create topic:
```shell
docker exec -it the-great-checkout-kafka-1 kafka-topics.sh --create --topic transactions --bootstrap-server localhost:9092
```

To list topics:
```shell
docker exec -it the-great-checkout-kafka-1 kafka-topics.sh --list --bootstrap-server localhost:9092
```

> See more in the-great-checkout on github!