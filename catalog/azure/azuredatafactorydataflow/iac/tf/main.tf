# One kind, two provider resources: Azure stores mapping data flows
# and flowlets in the SAME factory-scoped dataflow namespace, differing
# only in the ARM type token -- so the spec's `flowlet` flag selects
# which resource is created, and flipping it replaces the object (the
# resources share an ID shape but are distinct types to the provider).
#
# The transformation logic itself travels in the script (Azure owns
# that language); the source/sink/transformation blocks declare the
# named endpoints the script references.

resource "azurerm_data_factory_data_flow" "main" {
  count = var.spec.flowlet ? 0 : 1

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  script       = var.spec.script != "" ? var.spec.script : null
  script_lines = length(var.spec.script_lines) > 0 ? var.spec.script_lines : null

  dynamic "source" {
    for_each = var.spec.sources
    content {
      name        = source.value.name
      description = source.value.description != "" ? source.value.description : null

      dynamic "dataset" {
        for_each = source.value.dataset != null ? [source.value.dataset] : []
        content {
          name       = dataset.value.name
          parameters = length(dataset.value.parameters) > 0 ? dataset.value.parameters : null
        }
      }

      dynamic "flowlet" {
        for_each = source.value.flowlet != null ? [source.value.flowlet] : []
        content {
          name               = flowlet.value.name
          parameters         = length(flowlet.value.parameters) > 0 ? flowlet.value.parameters : null
          dataset_parameters = flowlet.value.dataset_parameters != "" ? flowlet.value.dataset_parameters : null
        }
      }

      dynamic "linked_service" {
        for_each = source.value.linked_service != null ? [source.value.linked_service] : []
        content {
          name       = linked_service.value.name
          parameters = length(linked_service.value.parameters) > 0 ? linked_service.value.parameters : null
        }
      }

      dynamic "schema_linked_service" {
        for_each = source.value.schema_linked_service != null ? [source.value.schema_linked_service] : []
        content {
          name       = schema_linked_service.value.name
          parameters = length(schema_linked_service.value.parameters) > 0 ? schema_linked_service.value.parameters : null
        }
      }
    }
  }

  dynamic "sink" {
    for_each = var.spec.sinks
    content {
      name        = sink.value.name
      description = sink.value.description != "" ? sink.value.description : null

      dynamic "dataset" {
        for_each = sink.value.dataset != null ? [sink.value.dataset] : []
        content {
          name       = dataset.value.name
          parameters = length(dataset.value.parameters) > 0 ? dataset.value.parameters : null
        }
      }

      dynamic "flowlet" {
        for_each = sink.value.flowlet != null ? [sink.value.flowlet] : []
        content {
          name               = flowlet.value.name
          parameters         = length(flowlet.value.parameters) > 0 ? flowlet.value.parameters : null
          dataset_parameters = flowlet.value.dataset_parameters != "" ? flowlet.value.dataset_parameters : null
        }
      }

      dynamic "linked_service" {
        for_each = sink.value.linked_service != null ? [sink.value.linked_service] : []
        content {
          name       = linked_service.value.name
          parameters = length(linked_service.value.parameters) > 0 ? linked_service.value.parameters : null
        }
      }

      dynamic "schema_linked_service" {
        for_each = sink.value.schema_linked_service != null ? [sink.value.schema_linked_service] : []
        content {
          name       = schema_linked_service.value.name
          parameters = length(schema_linked_service.value.parameters) > 0 ? schema_linked_service.value.parameters : null
        }
      }

      # Rejected-row routing exists on sinks only -- Azure's data flow
      # model (the provider accepts it on sources but silently drops
      # it, so the spec never offers it there).
      dynamic "rejected_linked_service" {
        for_each = sink.value.rejected_linked_service != null ? [sink.value.rejected_linked_service] : []
        content {
          name       = rejected_linked_service.value.name
          parameters = length(rejected_linked_service.value.parameters) > 0 ? rejected_linked_service.value.parameters : null
        }
      }
    }
  }

  dynamic "transformation" {
    for_each = var.spec.transformations
    content {
      name        = transformation.value.name
      description = transformation.value.description != "" ? transformation.value.description : null

      dynamic "dataset" {
        for_each = transformation.value.dataset != null ? [transformation.value.dataset] : []
        content {
          name       = dataset.value.name
          parameters = length(dataset.value.parameters) > 0 ? dataset.value.parameters : null
        }
      }

      dynamic "flowlet" {
        for_each = transformation.value.flowlet != null ? [transformation.value.flowlet] : []
        content {
          name               = flowlet.value.name
          parameters         = length(flowlet.value.parameters) > 0 ? flowlet.value.parameters : null
          dataset_parameters = flowlet.value.dataset_parameters != "" ? flowlet.value.dataset_parameters : null
        }
      }

      dynamic "linked_service" {
        for_each = transformation.value.linked_service != null ? [transformation.value.linked_service] : []
        content {
          name       = linked_service.value.name
          parameters = length(linked_service.value.parameters) > 0 ? linked_service.value.parameters : null
        }
      }
    }
  }

  annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  description = var.spec.description != "" ? var.spec.description : null
  folder      = var.spec.folder != "" ? var.spec.folder : null
}

# The flowlet form: identical surface, but sources and sinks are
# optional (the embedding data flow supplies them).
resource "azurerm_data_factory_flowlet_data_flow" "main" {
  count = var.spec.flowlet ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  script       = var.spec.script != "" ? var.spec.script : null
  script_lines = length(var.spec.script_lines) > 0 ? var.spec.script_lines : null

  dynamic "source" {
    for_each = var.spec.sources
    content {
      name        = source.value.name
      description = source.value.description != "" ? source.value.description : null

      dynamic "dataset" {
        for_each = source.value.dataset != null ? [source.value.dataset] : []
        content {
          name       = dataset.value.name
          parameters = length(dataset.value.parameters) > 0 ? dataset.value.parameters : null
        }
      }

      dynamic "flowlet" {
        for_each = source.value.flowlet != null ? [source.value.flowlet] : []
        content {
          name               = flowlet.value.name
          parameters         = length(flowlet.value.parameters) > 0 ? flowlet.value.parameters : null
          dataset_parameters = flowlet.value.dataset_parameters != "" ? flowlet.value.dataset_parameters : null
        }
      }

      dynamic "linked_service" {
        for_each = source.value.linked_service != null ? [source.value.linked_service] : []
        content {
          name       = linked_service.value.name
          parameters = length(linked_service.value.parameters) > 0 ? linked_service.value.parameters : null
        }
      }

      dynamic "schema_linked_service" {
        for_each = source.value.schema_linked_service != null ? [source.value.schema_linked_service] : []
        content {
          name       = schema_linked_service.value.name
          parameters = length(schema_linked_service.value.parameters) > 0 ? schema_linked_service.value.parameters : null
        }
      }
    }
  }

  dynamic "sink" {
    for_each = var.spec.sinks
    content {
      name        = sink.value.name
      description = sink.value.description != "" ? sink.value.description : null

      dynamic "dataset" {
        for_each = sink.value.dataset != null ? [sink.value.dataset] : []
        content {
          name       = dataset.value.name
          parameters = length(dataset.value.parameters) > 0 ? dataset.value.parameters : null
        }
      }

      dynamic "flowlet" {
        for_each = sink.value.flowlet != null ? [sink.value.flowlet] : []
        content {
          name               = flowlet.value.name
          parameters         = length(flowlet.value.parameters) > 0 ? flowlet.value.parameters : null
          dataset_parameters = flowlet.value.dataset_parameters != "" ? flowlet.value.dataset_parameters : null
        }
      }

      dynamic "linked_service" {
        for_each = sink.value.linked_service != null ? [sink.value.linked_service] : []
        content {
          name       = linked_service.value.name
          parameters = length(linked_service.value.parameters) > 0 ? linked_service.value.parameters : null
        }
      }

      dynamic "schema_linked_service" {
        for_each = sink.value.schema_linked_service != null ? [sink.value.schema_linked_service] : []
        content {
          name       = schema_linked_service.value.name
          parameters = length(schema_linked_service.value.parameters) > 0 ? schema_linked_service.value.parameters : null
        }
      }

      dynamic "rejected_linked_service" {
        for_each = sink.value.rejected_linked_service != null ? [sink.value.rejected_linked_service] : []
        content {
          name       = rejected_linked_service.value.name
          parameters = length(rejected_linked_service.value.parameters) > 0 ? rejected_linked_service.value.parameters : null
        }
      }
    }
  }

  dynamic "transformation" {
    for_each = var.spec.transformations
    content {
      name        = transformation.value.name
      description = transformation.value.description != "" ? transformation.value.description : null

      dynamic "dataset" {
        for_each = transformation.value.dataset != null ? [transformation.value.dataset] : []
        content {
          name       = dataset.value.name
          parameters = length(dataset.value.parameters) > 0 ? dataset.value.parameters : null
        }
      }

      dynamic "flowlet" {
        for_each = transformation.value.flowlet != null ? [transformation.value.flowlet] : []
        content {
          name               = flowlet.value.name
          parameters         = length(flowlet.value.parameters) > 0 ? flowlet.value.parameters : null
          dataset_parameters = flowlet.value.dataset_parameters != "" ? flowlet.value.dataset_parameters : null
        }
      }

      dynamic "linked_service" {
        for_each = transformation.value.linked_service != null ? [transformation.value.linked_service] : []
        content {
          name       = linked_service.value.name
          parameters = length(linked_service.value.parameters) > 0 ? linked_service.value.parameters : null
        }
      }
    }
  }

  annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  description = var.spec.description != "" ? var.spec.description : null
  folder      = var.spec.folder != "" ? var.spec.folder : null
}
