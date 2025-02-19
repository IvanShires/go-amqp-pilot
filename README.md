# Go AMQP Pilot

This lab test validates a Golang script that constructs and publishes a JSON-formatted dictionary to an AMQP message queue (e.g., RabbitMQ). The script connects to the AMQP broker, serializes a structured dictionary into JSON, and publishes it to a specified queue. The test includes:

    Setting up an AMQP connection and declaring a queue.
    Serializing a sample dictionary into JSON format.
    Publishing the JSON message to the queue.
    Verifying message delivery via a consumer or log output.

The test ensures proper connectivity, message serialization, and successful publishing.