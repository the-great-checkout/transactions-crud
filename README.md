# transactions-crud
Transactions SSOT

docker exec -it the-great-checkout-kafka-1 kafka-topics.sh --create --topic transactions --bootstrap-server localhost:9092
docker exec -it the-great-checkout-kafka-1 kafka-topics.sh --list --bootstrap-server localhost:9092