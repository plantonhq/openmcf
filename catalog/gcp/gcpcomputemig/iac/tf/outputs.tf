# Branch-independent outputs: one(concat()) collapses the zonal/regional
# resource pair so the same output keys carry the same shapes on every
# manifest — the contract the cross-engine conformance case relies on.

output "instance_group" {
  description = "Full URL of the group's instance group — the LB backend handle (a GcpBackendService backend's group takes exactly this)"
  value       = one(concat(google_compute_instance_group_manager.this[*].instance_group, google_compute_region_instance_group_manager.this[*].instance_group))
}

output "self_link" {
  description = "Self link of the instance group MANAGER resource"
  value       = one(concat(google_compute_instance_group_manager.this[*].self_link, google_compute_region_instance_group_manager.this[*].self_link))
}

output "current_template_self_link" {
  description = "The active template's link (unique form on zonal groups — changes on every template rotation)"
  value       = one(concat(google_compute_instance_template.this[*].self_link_unique, google_compute_region_instance_template.this[*].self_link))
}

output "mig_name" {
  description = "Name of the managed instance group as it exists in GCP"
  value       = one(concat(google_compute_instance_group_manager.this[*].name, google_compute_region_instance_group_manager.this[*].name))
}

output "location" {
  description = "The group's location: the zone of a zonal group or the region of a regional one"
  value       = local.location
}
