# Stack outputs — exactly the DigitalOceanDatabaseKafkaTopicStackOutputs
# contract, identical across both provisioners. The (cluster, topic name)
# pair is the topic's API identity; DigitalOcean mints no standalone id.

output "cluster_id" {
  description = "UUID of the Kafka database cluster the topic lives in"
  value       = digitalocean_database_kafka_topic.topic.cluster_id
}

output "topic_name" {
  description = "Name of the Kafka topic (its API identity within the cluster)"
  value       = digitalocean_database_kafka_topic.topic.name
}

output "state" {
  description = "Provisioning state of the topic as reported by DigitalOcean at apply time"
  value       = digitalocean_database_kafka_topic.topic.state
}
