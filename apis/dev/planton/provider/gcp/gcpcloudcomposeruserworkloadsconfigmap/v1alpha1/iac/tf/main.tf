# The user workloads ConfigMap — a Kubernetes ConfigMap Composer manages
# in the environment's GKE cluster. Airflow DAGs (KubernetesPodOperator
# tasks) consume it by name for non-secret configuration: feature flags,
# endpoints, tuning parameters.
#
# The data updates in place; name, environment, region, and project are
# immutable. The environment must already exist — it is a first-class
# resource this one composes against by reference.
#
# No API enablement here: the Composer API is enabled by the environment
# this ConfigMap is delivered into (a ConfigMap cannot exist without one).
resource "google_composer_user_workloads_config_map" "config_map" {
  name        = local.config_map_name
  environment = local.environment
  region      = local.region
  project     = local.project_id

  data = var.spec.data
}
