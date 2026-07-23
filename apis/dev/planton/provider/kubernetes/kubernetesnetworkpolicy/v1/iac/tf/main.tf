# Kubernetes NetworkPolicy Terraform module.
#
# The three peer forms (pod selector, namespace selector, IP block) pass
# through with their exact combination semantics — pod+namespace selectors in
# ONE peer are an AND (pods matching the selector in the selected namespaces),
# which is why each spec peer maps to exactly one API peer and is never split.

resource "kubernetes_network_policy_v1" "network_policy" {
  metadata {
    name        = var.spec.name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  spec {
    # An absent pod_selector means "all pods in the namespace" — the empty
    # selector block is the correct wire form for that (the default-deny
    # building block).
    pod_selector {
      match_labels = try(var.spec.pod_selector.match_labels, null) != null && length(try(var.spec.pod_selector.match_labels, {})) > 0 ? var.spec.pod_selector.match_labels : null
      dynamic "match_expressions" {
        for_each = try(var.spec.pod_selector.match_expressions, null) != null ? var.spec.pod_selector.match_expressions : []
        content {
          key      = match_expressions.value.key
          operator = match_expressions.value.operator
          values   = length(match_expressions.value.values) > 0 ? match_expressions.value.values : null
        }
      }
    }

    policy_types = local.policy_types

    dynamic "ingress" {
      for_each = try(var.spec.ingress_rules, [])
      content {
        dynamic "from" {
          for_each = try(ingress.value.from, [])
          content {
            dynamic "pod_selector" {
              for_each = from.value.pod_selector != null ? [from.value.pod_selector] : []
              content {
                match_labels = length(pod_selector.value.match_labels) > 0 ? pod_selector.value.match_labels : null
                dynamic "match_expressions" {
                  for_each = pod_selector.value.match_expressions
                  content {
                    key      = match_expressions.value.key
                    operator = match_expressions.value.operator
                    values   = length(match_expressions.value.values) > 0 ? match_expressions.value.values : null
                  }
                }
              }
            }
            dynamic "namespace_selector" {
              for_each = from.value.namespace_selector != null ? [from.value.namespace_selector] : []
              content {
                match_labels = length(namespace_selector.value.match_labels) > 0 ? namespace_selector.value.match_labels : null
                dynamic "match_expressions" {
                  for_each = namespace_selector.value.match_expressions
                  content {
                    key      = match_expressions.value.key
                    operator = match_expressions.value.operator
                    values   = length(match_expressions.value.values) > 0 ? match_expressions.value.values : null
                  }
                }
              }
            }
            dynamic "ip_block" {
              for_each = from.value.ip_block != null ? [from.value.ip_block] : []
              content {
                cidr   = ip_block.value.cidr
                except = length(ip_block.value.except) > 0 ? ip_block.value.except : null
              }
            }
          }
        }
        dynamic "ports" {
          for_each = try(ingress.value.ports, [])
          content {
            protocol = ports.value.protocol
            # "" matches all ports for the protocol; a numeric string matches
            # a port number, anything else a named container port.
            port     = ports.value.port != "" ? ports.value.port : null
            end_port = ports.value.end_port > 0 ? ports.value.end_port : null
          }
        }
      }
    }

    dynamic "egress" {
      for_each = try(var.spec.egress_rules, [])
      content {
        dynamic "to" {
          for_each = try(egress.value.to, [])
          content {
            dynamic "pod_selector" {
              for_each = to.value.pod_selector != null ? [to.value.pod_selector] : []
              content {
                match_labels = length(pod_selector.value.match_labels) > 0 ? pod_selector.value.match_labels : null
                dynamic "match_expressions" {
                  for_each = pod_selector.value.match_expressions
                  content {
                    key      = match_expressions.value.key
                    operator = match_expressions.value.operator
                    values   = length(match_expressions.value.values) > 0 ? match_expressions.value.values : null
                  }
                }
              }
            }
            dynamic "namespace_selector" {
              for_each = to.value.namespace_selector != null ? [to.value.namespace_selector] : []
              content {
                match_labels = length(namespace_selector.value.match_labels) > 0 ? namespace_selector.value.match_labels : null
                dynamic "match_expressions" {
                  for_each = namespace_selector.value.match_expressions
                  content {
                    key      = match_expressions.value.key
                    operator = match_expressions.value.operator
                    values   = length(match_expressions.value.values) > 0 ? match_expressions.value.values : null
                  }
                }
              }
            }
            dynamic "ip_block" {
              for_each = to.value.ip_block != null ? [to.value.ip_block] : []
              content {
                cidr   = ip_block.value.cidr
                except = length(ip_block.value.except) > 0 ? ip_block.value.except : null
              }
            }
          }
        }
        dynamic "ports" {
          for_each = try(egress.value.ports, [])
          content {
            protocol = ports.value.protocol
            port     = ports.value.port != "" ? ports.value.port : null
            end_port = ports.value.end_port > 0 ? ports.value.end_port : null
          }
        }
      }
    }
  }
}
