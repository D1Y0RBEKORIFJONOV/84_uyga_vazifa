#!/bin/bash

echo "Waiting for Kafka to start..."
cub kafka-ready -b broker:9092 1 20

echo "Creating topics..."
kafka-topics --create --topic USER_CREATE --partitions 1 --replication-factor 1 --if-not-exists --bootstrap-server broker:9092
kafka-topics --create --topic VFY_VFY --partitions 1 --replication-factor 1 --if-not-exists --bootstrap-server broker:9092
