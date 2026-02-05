#!/bin/bash
# NUDEX RabbitMQ Initialization Script

# Wait for RabbitMQ to be ready
sleep 10

# Create exchanges
rabbitmqadmin declare exchange name=nudex.events type=topic durable=true
rabbitmqadmin declare exchange name=nudex.jobs type=topic durable=true

# Create queues for events
rabbitmqadmin declare queue name=nudex.events.catalog durable=true
rabbitmqadmin declare queue name=nudex.events.users durable=true
rabbitmqadmin declare queue name=nudex.events.library durable=true
rabbitmqadmin declare queue name=nudex.events.playback durable=true

# Create queues for jobs
rabbitmqadmin declare queue name=nudex.jobs.ingestion durable=true

# Bind event queues
rabbitmqadmin declare binding source=nudex.events destination=nudex.events.catalog routing_key="catalog.*"
rabbitmqadmin declare binding source=nudex.events destination=nudex.events.users routing_key="user.*"
rabbitmqadmin declare binding source=nudex.events destination=nudex.events.library routing_key="library.*"
rabbitmqadmin declare binding source=nudex.events destination=nudex.events.playback routing_key="playback.*"

# Bind job queues
rabbitmqadmin declare binding source=nudex.jobs destination=nudex.jobs.ingestion routing_key="ingestion.*"

echo "RabbitMQ NUDEX exchanges and queues initialized successfully"