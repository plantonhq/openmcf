# DigitalOcean Database Kafka Schema
#
# Registers one schema subject in a DigitalOcean managed Kafka cluster's
# schema registry -- the complete
# digitalocean_database_kafka_schema_registry resource surface.
#
# EVERY argument is create-only: the provider has no update path, so any
# change -- including evolving the schema definition, even a
# whitespace-only reformat (the definition is compared verbatim) --
# destroys the subject and re-registers it, which DROPS all previously
# registered versions. Treat schema evolution as a deliberate replacement.

resource "digitalocean_database_kafka_schema_registry" "schema" {
  cluster_id   = var.spec.cluster
  subject_name = var.spec.subject_name
  schema_type  = var.spec.schema_type
  schema       = var.spec.schema
}
