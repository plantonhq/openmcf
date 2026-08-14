# Amazon Bedrock knowledge base: the RAG store (vector, managed, Kendra,
# or SQL retrieval) plus its folded data sources (S3, web crawl,
# Confluence, Salesforce, SharePoint, managed connectors), each with its
# chunk-parse-transform ingestion pipeline.
#
# Nearly every argument below is create-time only upstream -- changing the
# type, storage, or ingestion configuration replaces the knowledge base or
# data source (AWS re-ingests afterwards). The provider retries the
# IAM/data-access propagation classes at create; the module adds none.

resource "aws_bedrockagent_knowledge_base" "this" {
  # Create-time naming basis; doubles as the Name tag. metadata.name on
  # both engines.
  name = local.knowledge_base_name

  role_arn = var.spec.role_arn

  description = var.spec.description != "" ? var.spec.description : null

  # -------------------------------------------------------------------
  # Knowledge-base type (the discriminator is derived from which spec arm
  # is set -- exactly one, per the spec's CEL guards)
  # -------------------------------------------------------------------
  knowledge_base_configuration {
    type = local.kb_type

    dynamic "vector_knowledge_base_configuration" {
      for_each = var.spec.vector != null ? [var.spec.vector] : []
      content {
        embedding_model_arn = vector_knowledge_base_configuration.value.embedding_model_arn

        dynamic "embedding_model_configuration" {
          for_each = vector_knowledge_base_configuration.value.embedding_model != null ? [vector_knowledge_base_configuration.value.embedding_model] : []
          content {
            dynamic "bedrock_embedding_model_configuration" {
              for_each = [embedding_model_configuration.value]
              content {
                dimensions          = bedrock_embedding_model_configuration.value.dimensions != 0 ? bedrock_embedding_model_configuration.value.dimensions : null
                embedding_data_type = bedrock_embedding_model_configuration.value.embedding_data_type != "" ? bedrock_embedding_model_configuration.value.embedding_data_type : null
                dynamic "audio" {
                  for_each = bedrock_embedding_model_configuration.value.audio_segmentation_seconds != 0 ? [bedrock_embedding_model_configuration.value.audio_segmentation_seconds] : []
                  content {
                    segmentation_configuration {
                      fixed_length_duration = audio.value
                    }
                  }
                }
                dynamic "video" {
                  for_each = bedrock_embedding_model_configuration.value.video_segmentation_seconds != 0 ? [bedrock_embedding_model_configuration.value.video_segmentation_seconds] : []
                  content {
                    segmentation_configuration {
                      fixed_length_duration = video.value
                    }
                  }
                }
              }
            }
          }
        }

        # The supplemental S3 location rides a fixed two-level wrapper
        # upstream (storage_location type S3 + s3_location.uri); the spec
        # carries the URI leaf directly.
        dynamic "supplemental_data_storage_configuration" {
          for_each = vector_knowledge_base_configuration.value.supplemental_data_s3_uri != "" ? [vector_knowledge_base_configuration.value.supplemental_data_s3_uri] : []
          content {
            storage_location {
              type = "S3"
              s3_location {
                uri = supplemental_data_storage_configuration.value
              }
            }
          }
        }
      }
    }

    dynamic "managed_knowledge_base_configuration" {
      for_each = var.spec.managed != null ? [var.spec.managed] : []
      content {
        embedding_model_arn = managed_knowledge_base_configuration.value.embedding_model_arn != "" ? managed_knowledge_base_configuration.value.embedding_model_arn : null
        # Derived discriminator, ALWAYS sent: AWS's embeddingModelType is
        # CUSTOM exactly when an embedding-model ARN is brought, MANAGED
        # otherwise. Sending it keeps the value known at plan time -- the
        # provider marks it Optional+Computed, and leaving it unknown
        # strands the created KB outside Pulumi state on the twin module
        # (bridge apply error), so both engines own the derivation.
        embedding_model_type = managed_knowledge_base_configuration.value.embedding_model_arn != "" ? "CUSTOM" : "MANAGED"

        dynamic "embedding_model_configuration" {
          for_each = managed_knowledge_base_configuration.value.embedding_model != null ? [managed_knowledge_base_configuration.value.embedding_model] : []
          content {
            dynamic "bedrock_embedding_model_configuration" {
              for_each = [embedding_model_configuration.value]
              content {
                dimensions          = bedrock_embedding_model_configuration.value.dimensions != 0 ? bedrock_embedding_model_configuration.value.dimensions : null
                embedding_data_type = bedrock_embedding_model_configuration.value.embedding_data_type != "" ? bedrock_embedding_model_configuration.value.embedding_data_type : null
                dynamic "audio" {
                  for_each = bedrock_embedding_model_configuration.value.audio_segmentation_seconds != 0 ? [bedrock_embedding_model_configuration.value.audio_segmentation_seconds] : []
                  content {
                    segmentation_configuration {
                      fixed_length_duration = audio.value
                    }
                  }
                }
                dynamic "video" {
                  for_each = bedrock_embedding_model_configuration.value.video_segmentation_seconds != 0 ? [bedrock_embedding_model_configuration.value.video_segmentation_seconds] : []
                  content {
                    segmentation_configuration {
                      fixed_length_duration = video.value
                    }
                  }
                }
              }
            }
          }
        }

        dynamic "server_side_encryption_configuration" {
          for_each = managed_knowledge_base_configuration.value.kms_key_arn != "" ? [managed_knowledge_base_configuration.value.kms_key_arn] : []
          content {
            kms_key_arn = server_side_encryption_configuration.value
          }
        }
      }
    }

    dynamic "kendra_knowledge_base_configuration" {
      for_each = var.spec.kendra != null ? [var.spec.kendra] : []
      content {
        kendra_index_arn = kendra_knowledge_base_configuration.value.kendra_index_arn
      }
    }

    dynamic "sql_knowledge_base_configuration" {
      for_each = var.spec.sql != null ? [var.spec.sql] : []
      content {
        # REDSHIFT is the only SQL engine AWS defines -- the module owns
        # the constant.
        type = "REDSHIFT"

        redshift_configuration {
          query_engine_configuration {
            type = sql_knowledge_base_configuration.value.provisioned != null ? "PROVISIONED" : "SERVERLESS"

            dynamic "provisioned_configuration" {
              for_each = sql_knowledge_base_configuration.value.provisioned != null ? [sql_knowledge_base_configuration.value.provisioned] : []
              content {
                cluster_identifier = provisioned_configuration.value.cluster_identifier
                auth_configuration {
                  type                         = provisioned_configuration.value.auth.type
                  database_user                = provisioned_configuration.value.auth.database_user != "" ? provisioned_configuration.value.auth.database_user : null
                  username_password_secret_arn = provisioned_configuration.value.auth.username_password_secret_arn != "" ? provisioned_configuration.value.auth.username_password_secret_arn : null
                }
              }
            }

            dynamic "serverless_configuration" {
              for_each = sql_knowledge_base_configuration.value.serverless != null ? [sql_knowledge_base_configuration.value.serverless] : []
              content {
                workgroup_arn = serverless_configuration.value.workgroup_arn
                auth_configuration {
                  type                         = serverless_configuration.value.auth.type
                  username_password_secret_arn = serverless_configuration.value.auth.username_password_secret_arn != "" ? serverless_configuration.value.auth.username_password_secret_arn : null
                }
              }
            }
          }

          storage_configuration {
            type = sql_knowledge_base_configuration.value.warehouse.data_catalog != null ? "AWS_DATA_CATALOG" : "REDSHIFT"

            dynamic "aws_data_catalog_configuration" {
              for_each = sql_knowledge_base_configuration.value.warehouse.data_catalog != null ? [sql_knowledge_base_configuration.value.warehouse.data_catalog] : []
              content {
                table_names = aws_data_catalog_configuration.value.table_names
              }
            }

            dynamic "redshift_configuration" {
              for_each = sql_knowledge_base_configuration.value.warehouse.redshift != null ? [sql_knowledge_base_configuration.value.warehouse.redshift] : []
              content {
                database_name = redshift_configuration.value.database_name
              }
            }
          }

          dynamic "query_generation_configuration" {
            for_each = sql_knowledge_base_configuration.value.query_generation != null ? [sql_knowledge_base_configuration.value.query_generation] : []
            content {
              execution_timeout_seconds = query_generation_configuration.value.execution_timeout_seconds != 0 ? query_generation_configuration.value.execution_timeout_seconds : null

              dynamic "generation_context" {
                for_each = (length(query_generation_configuration.value.curated_queries) > 0 || length(query_generation_configuration.value.tables) > 0) ? [query_generation_configuration.value] : []
                content {
                  dynamic "curated_query" {
                    for_each = generation_context.value.curated_queries
                    content {
                      natural_language = curated_query.value.natural_language
                      sql              = curated_query.value.sql
                    }
                  }
                  dynamic "table" {
                    for_each = generation_context.value.tables
                    content {
                      name        = table.value.name
                      description = table.value.description != "" ? table.value.description : null
                      inclusion   = table.value.inclusion != "" ? table.value.inclusion : null
                      dynamic "column" {
                        for_each = table.value.columns
                        content {
                          name        = column.value.name
                          description = column.value.description != "" ? column.value.description : null
                          inclusion   = column.value.inclusion != "" ? column.value.inclusion : null
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }

  # -------------------------------------------------------------------
  # Vector store (required with the vector type, absent otherwise -- the
  # spec's CEL guards enforce the pairing)
  # -------------------------------------------------------------------
  dynamic "storage_configuration" {
    for_each = var.spec.storage != null ? [var.spec.storage] : []
    content {
      type = local.storage_type

      dynamic "opensearch_serverless_configuration" {
        for_each = storage_configuration.value.opensearch_serverless != null ? [storage_configuration.value.opensearch_serverless] : []
        content {
          collection_arn    = opensearch_serverless_configuration.value.collection_arn
          vector_index_name = opensearch_serverless_configuration.value.vector_index_name
          field_mapping {
            vector_field   = opensearch_serverless_configuration.value.field_mapping.vector_field
            text_field     = opensearch_serverless_configuration.value.field_mapping.text_field
            metadata_field = opensearch_serverless_configuration.value.field_mapping.metadata_field
          }
        }
      }

      dynamic "opensearch_managed_cluster_configuration" {
        for_each = storage_configuration.value.opensearch_managed != null ? [storage_configuration.value.opensearch_managed] : []
        content {
          domain_arn        = opensearch_managed_cluster_configuration.value.domain_arn
          domain_endpoint   = opensearch_managed_cluster_configuration.value.domain_endpoint
          vector_index_name = opensearch_managed_cluster_configuration.value.vector_index_name
          field_mapping {
            vector_field   = opensearch_managed_cluster_configuration.value.field_mapping.vector_field
            text_field     = opensearch_managed_cluster_configuration.value.field_mapping.text_field
            metadata_field = opensearch_managed_cluster_configuration.value.field_mapping.metadata_field
          }
        }
      }

      dynamic "s3_vectors_configuration" {
        for_each = storage_configuration.value.s3_vectors != null ? [storage_configuration.value.s3_vectors] : []
        content {
          index_arn         = s3_vectors_configuration.value.index_arn != "" ? s3_vectors_configuration.value.index_arn : null
          index_name        = s3_vectors_configuration.value.index_name != "" ? s3_vectors_configuration.value.index_name : null
          vector_bucket_arn = s3_vectors_configuration.value.vector_bucket_arn != "" ? s3_vectors_configuration.value.vector_bucket_arn : null
        }
      }

      dynamic "rds_configuration" {
        for_each = storage_configuration.value.rds != null ? [storage_configuration.value.rds] : []
        content {
          resource_arn           = rds_configuration.value.resource_arn
          credentials_secret_arn = rds_configuration.value.credentials_secret_arn
          database_name          = rds_configuration.value.database_name
          table_name             = rds_configuration.value.table_name
          field_mapping {
            vector_field          = rds_configuration.value.field_mapping.vector_field
            text_field            = rds_configuration.value.field_mapping.text_field
            metadata_field        = rds_configuration.value.field_mapping.metadata_field
            primary_key_field     = rds_configuration.value.field_mapping.primary_key_field
            custom_metadata_field = rds_configuration.value.field_mapping.custom_metadata_field != "" ? rds_configuration.value.field_mapping.custom_metadata_field : null
          }
        }
      }

      dynamic "pinecone_configuration" {
        for_each = storage_configuration.value.pinecone != null ? [storage_configuration.value.pinecone] : []
        content {
          connection_string      = pinecone_configuration.value.connection_string
          credentials_secret_arn = pinecone_configuration.value.credentials_secret_arn
          namespace              = pinecone_configuration.value.namespace != "" ? pinecone_configuration.value.namespace : null
          field_mapping {
            text_field     = pinecone_configuration.value.field_mapping.text_field
            metadata_field = pinecone_configuration.value.field_mapping.metadata_field
          }
        }
      }

      dynamic "mongo_db_atlas_configuration" {
        for_each = storage_configuration.value.mongodb_atlas != null ? [storage_configuration.value.mongodb_atlas] : []
        content {
          endpoint               = mongo_db_atlas_configuration.value.endpoint
          database_name          = mongo_db_atlas_configuration.value.database_name
          collection_name        = mongo_db_atlas_configuration.value.collection_name
          vector_index_name      = mongo_db_atlas_configuration.value.vector_index_name
          text_index_name        = mongo_db_atlas_configuration.value.text_index_name != "" ? mongo_db_atlas_configuration.value.text_index_name : null
          endpoint_service_name  = mongo_db_atlas_configuration.value.endpoint_service_name != "" ? mongo_db_atlas_configuration.value.endpoint_service_name : null
          credentials_secret_arn = mongo_db_atlas_configuration.value.credentials_secret_arn
          field_mapping {
            vector_field   = mongo_db_atlas_configuration.value.field_mapping.vector_field
            text_field     = mongo_db_atlas_configuration.value.field_mapping.text_field
            metadata_field = mongo_db_atlas_configuration.value.field_mapping.metadata_field
          }
        }
      }

      dynamic "neptune_analytics_configuration" {
        for_each = storage_configuration.value.neptune_analytics != null ? [storage_configuration.value.neptune_analytics] : []
        content {
          graph_arn = neptune_analytics_configuration.value.graph_arn
          field_mapping {
            text_field     = neptune_analytics_configuration.value.field_mapping.text_field
            metadata_field = neptune_analytics_configuration.value.field_mapping.metadata_field
          }
        }
      }

      dynamic "redis_enterprise_cloud_configuration" {
        for_each = storage_configuration.value.redis_enterprise_cloud != null ? [storage_configuration.value.redis_enterprise_cloud] : []
        content {
          endpoint               = redis_enterprise_cloud_configuration.value.endpoint
          vector_index_name      = redis_enterprise_cloud_configuration.value.vector_index_name
          credentials_secret_arn = redis_enterprise_cloud_configuration.value.credentials_secret_arn
          dynamic "field_mapping" {
            for_each = redis_enterprise_cloud_configuration.value.field_mapping != null ? [redis_enterprise_cloud_configuration.value.field_mapping] : []
            content {
              vector_field   = field_mapping.value.vector_field != "" ? field_mapping.value.vector_field : null
              text_field     = field_mapping.value.text_field != "" ? field_mapping.value.text_field : null
              metadata_field = field_mapping.value.metadata_field != "" ? field_mapping.value.metadata_field : null
            }
          }
        }
      }
    }
  }

  tags = local.aws_tags
}

# Document connectors, keyed by their stable entry names. The connector
# type discriminator is derived from which spec arm is set. Entries are
# `any`-typed (heterogeneous arms -- see variables.tf), so every optional
# field is read with try().
resource "aws_bedrockagent_data_source" "this" {
  for_each = local.data_sources

  name              = each.value.name
  knowledge_base_id = aws_bedrockagent_knowledge_base.this.id

  description          = try(each.value.description, "") != "" ? each.value.description : null
  data_deletion_policy = try(each.value.data_deletion_policy, "") != "" ? each.value.data_deletion_policy : null

  dynamic "server_side_encryption_configuration" {
    for_each = try(each.value.kms_key_arn, "") != "" ? [each.value.kms_key_arn] : []
    content {
      kms_key_arn = server_side_encryption_configuration.value
    }
  }

  data_source_configuration {
    type = try(each.value.s3, null) != null ? "S3" : (try(each.value.web, null) != null ? "WEB" : (try(each.value.confluence, null) != null ? "CONFLUENCE" : (try(each.value.salesforce, null) != null ? "SALESFORCE" : (try(each.value.sharepoint, null) != null ? "SHAREPOINT" : "MANAGED_KNOWLEDGE_BASE_CONNECTOR"))))

    dynamic "s3_configuration" {
      for_each = try(each.value.s3, null) != null ? [each.value.s3] : []
      content {
        bucket_arn              = s3_configuration.value.bucket_arn
        inclusion_prefixes      = try(s3_configuration.value.inclusion_prefix, "") != "" ? [s3_configuration.value.inclusion_prefix] : null
        bucket_owner_account_id = try(s3_configuration.value.bucket_owner_account_id, "") != "" ? s3_configuration.value.bucket_owner_account_id : null
      }
    }

    dynamic "web_configuration" {
      for_each = try(each.value.web, null) != null ? [each.value.web] : []
      content {
        source_configuration {
          url_configuration {
            dynamic "seed_urls" {
              for_each = web_configuration.value.seed_urls
              content {
                url = seed_urls.value
              }
            }
          }
        }
        crawler_configuration {
          scope             = try(web_configuration.value.scope, "") != "" ? web_configuration.value.scope : null
          inclusion_filters = length(try(web_configuration.value.inclusion_filters, [])) > 0 ? web_configuration.value.inclusion_filters : null
          exclusion_filters = length(try(web_configuration.value.exclusion_filters, [])) > 0 ? web_configuration.value.exclusion_filters : null
          user_agent        = try(web_configuration.value.user_agent, "") != "" ? web_configuration.value.user_agent : null
          dynamic "crawler_limits" {
            for_each = (try(web_configuration.value.max_pages, 0) != 0 || try(web_configuration.value.rate_limit, 0) != 0) ? [web_configuration.value] : []
            content {
              max_pages  = try(crawler_limits.value.max_pages, 0) != 0 ? crawler_limits.value.max_pages : null
              rate_limit = try(crawler_limits.value.rate_limit, 0) != 0 ? crawler_limits.value.rate_limit : null
            }
          }
        }
      }
    }

    dynamic "confluence_configuration" {
      for_each = try(each.value.confluence, null) != null ? [each.value.confluence] : []
      content {
        source_configuration {
          # SAAS is the only Confluence host type AWS defines -- the
          # module owns the constant.
          host_type              = "SAAS"
          host_url               = confluence_configuration.value.host_url
          auth_type              = confluence_configuration.value.auth_type
          credentials_secret_arn = confluence_configuration.value.credentials_secret_arn
        }
        dynamic "crawler_configuration" {
          for_each = length(try(confluence_configuration.value.filters, [])) > 0 ? [confluence_configuration.value.filters] : []
          content {
            filter_configuration {
              # PATTERN is the only filter type AWS defines -- the module
              # owns the constant.
              type = "PATTERN"
              pattern_object_filter {
                dynamic "filters" {
                  for_each = crawler_configuration.value
                  content {
                    object_type       = filters.value.object_type
                    inclusion_filters = length(try(filters.value.inclusion_filters, [])) > 0 ? filters.value.inclusion_filters : null
                    exclusion_filters = length(try(filters.value.exclusion_filters, [])) > 0 ? filters.value.exclusion_filters : null
                  }
                }
              }
            }
          }
        }
      }
    }

    dynamic "salesforce_configuration" {
      for_each = try(each.value.salesforce, null) != null ? [each.value.salesforce] : []
      content {
        source_configuration {
          # OAUTH2_CLIENT_CREDENTIALS is the only Salesforce auth type AWS
          # defines -- the module owns the constant.
          auth_type              = "OAUTH2_CLIENT_CREDENTIALS"
          host_url               = salesforce_configuration.value.host_url
          credentials_secret_arn = salesforce_configuration.value.credentials_secret_arn
        }
        dynamic "crawler_configuration" {
          for_each = length(try(salesforce_configuration.value.filters, [])) > 0 ? [salesforce_configuration.value.filters] : []
          content {
            filter_configuration {
              type = "PATTERN"
              pattern_object_filter {
                dynamic "filters" {
                  for_each = crawler_configuration.value
                  content {
                    object_type       = filters.value.object_type
                    inclusion_filters = length(try(filters.value.inclusion_filters, [])) > 0 ? filters.value.inclusion_filters : null
                    exclusion_filters = length(try(filters.value.exclusion_filters, [])) > 0 ? filters.value.exclusion_filters : null
                  }
                }
              }
            }
          }
        }
      }
    }

    dynamic "share_point_configuration" {
      for_each = try(each.value.sharepoint, null) != null ? [each.value.sharepoint] : []
      content {
        source_configuration {
          # ONLINE is the only SharePoint host type AWS defines -- the
          # module owns the constant.
          host_type              = "ONLINE"
          site_urls              = share_point_configuration.value.site_urls
          domain                 = share_point_configuration.value.domain
          tenant_id              = try(share_point_configuration.value.tenant_id, "") != "" ? share_point_configuration.value.tenant_id : null
          auth_type              = share_point_configuration.value.auth_type
          credentials_secret_arn = share_point_configuration.value.credentials_secret_arn
        }
        dynamic "crawler_configuration" {
          for_each = length(try(share_point_configuration.value.filters, [])) > 0 ? [share_point_configuration.value.filters] : []
          content {
            filter_configuration {
              type = "PATTERN"
              pattern_object_filter {
                dynamic "filters" {
                  for_each = crawler_configuration.value
                  content {
                    object_type       = filters.value.object_type
                    inclusion_filters = length(try(filters.value.inclusion_filters, [])) > 0 ? filters.value.inclusion_filters : null
                    exclusion_filters = length(try(filters.value.exclusion_filters, [])) > 0 ? filters.value.exclusion_filters : null
                  }
                }
              }
            }
          }
        }
      }
    }

    dynamic "managed_knowledge_base_connector_configuration" {
      for_each = try(each.value.managed_connector, null) != null ? [each.value.managed_connector] : []
      content {
        connector_parameters = try(managed_knowledge_base_connector_configuration.value.connector_parameters, null) != null ? jsonencode(managed_knowledge_base_connector_configuration.value.connector_parameters) : null

        dynamic "deletion_protection_configuration" {
          for_each = try(managed_knowledge_base_connector_configuration.value.deletion_protection, null) != null ? [managed_knowledge_base_connector_configuration.value.deletion_protection] : []
          content {
            deletion_protection_status    = try(deletion_protection_configuration.value.enabled, false) ? "ENABLED" : "DISABLED"
            deletion_protection_threshold = try(deletion_protection_configuration.value.threshold_percent, 0) != 0 ? deletion_protection_configuration.value.threshold_percent : null
          }
        }

        dynamic "media_extraction_configuration" {
          for_each = try(managed_knowledge_base_connector_configuration.value.media_extraction, null) != null ? [managed_knowledge_base_connector_configuration.value.media_extraction] : []
          content {
            audio_extraction_configuration {
              audio_extraction_status = try(media_extraction_configuration.value.audio, false) ? "ENABLED" : "DISABLED"
            }
            image_extraction_configuration {
              image_extraction_status = try(media_extraction_configuration.value.image, false) ? "ENABLED" : "DISABLED"
            }
            video_extraction_configuration {
              video_extraction_status = try(media_extraction_configuration.value.video, false) ? "ENABLED" : "DISABLED"
            }
          }
        }
      }
    }
  }

  dynamic "vector_ingestion_configuration" {
    for_each = try(each.value.vector_ingestion, null) != null ? [each.value.vector_ingestion] : []
    content {
      dynamic "chunking_configuration" {
        for_each = try(vector_ingestion_configuration.value.chunking, null) != null ? [vector_ingestion_configuration.value.chunking] : []
        content {
          chunking_strategy = chunking_configuration.value.strategy

          dynamic "fixed_size_chunking_configuration" {
            for_each = try(chunking_configuration.value.fixed_size, null) != null ? [chunking_configuration.value.fixed_size] : []
            content {
              max_tokens         = fixed_size_chunking_configuration.value.max_tokens
              overlap_percentage = fixed_size_chunking_configuration.value.overlap_percentage
            }
          }

          dynamic "hierarchical_chunking_configuration" {
            for_each = try(chunking_configuration.value.hierarchical, null) != null ? [chunking_configuration.value.hierarchical] : []
            content {
              overlap_tokens = hierarchical_chunking_configuration.value.overlap_tokens
              dynamic "level_configuration" {
                for_each = hierarchical_chunking_configuration.value.levels
                content {
                  max_tokens = level_configuration.value.max_tokens
                }
              }
            }
          }

          dynamic "semantic_chunking_configuration" {
            for_each = try(chunking_configuration.value.semantic, null) != null ? [chunking_configuration.value.semantic] : []
            content {
              breakpoint_percentile_threshold = semantic_chunking_configuration.value.breakpoint_percentile_threshold
              buffer_size                     = try(semantic_chunking_configuration.value.buffer_size, 0)
              # The provider spells this one singular.
              max_token = semantic_chunking_configuration.value.max_tokens
            }
          }
        }
      }

      dynamic "parsing_configuration" {
        for_each = try(vector_ingestion_configuration.value.parsing, null) != null ? [vector_ingestion_configuration.value.parsing] : []
        content {
          parsing_strategy = parsing_configuration.value.strategy

          dynamic "bedrock_data_automation_configuration" {
            for_each = parsing_configuration.value.strategy == "BEDROCK_DATA_AUTOMATION" && try(parsing_configuration.value.multimodal, false) ? [1] : []
            content {
              # MULTIMODAL is the only parsing modality AWS defines -- the
              # spec models it as a bool and the module owns the constant.
              parsing_modality = "MULTIMODAL"
            }
          }

          dynamic "bedrock_foundation_model_configuration" {
            for_each = try(parsing_configuration.value.foundation_model, null) != null ? [parsing_configuration.value.foundation_model] : []
            content {
              model_arn        = bedrock_foundation_model_configuration.value.model_arn
              parsing_modality = try(bedrock_foundation_model_configuration.value.multimodal, false) ? "MULTIMODAL" : null
              dynamic "parsing_prompt" {
                for_each = try(bedrock_foundation_model_configuration.value.parsing_prompt, "") != "" ? [bedrock_foundation_model_configuration.value.parsing_prompt] : []
                content {
                  parsing_prompt_string = parsing_prompt.value
                }
              }
            }
          }
        }
      }

      dynamic "custom_transformation_configuration" {
        for_each = try(vector_ingestion_configuration.value.custom_transformation, null) != null ? [vector_ingestion_configuration.value.custom_transformation] : []
        content {
          intermediate_storage {
            s3_location {
              uri = custom_transformation_configuration.value.intermediate_s3_uri
            }
          }
          transformation {
            # POST_CHUNKING is the only transformation step AWS defines --
            # the module owns the constant.
            step_to_apply = "POST_CHUNKING"
            transformation_function {
              transformation_lambda_configuration {
                lambda_arn = custom_transformation_configuration.value.lambda_arn
              }
            }
          }
        }
      }
    }
  }
}
