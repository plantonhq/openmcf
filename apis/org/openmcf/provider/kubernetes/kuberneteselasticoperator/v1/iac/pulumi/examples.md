# Kubernetes Elastic Operator - Pulumi Examples

This document provides practical examples for deploying the ECK operator using Pulumi.

## Example 1: Basic Deployment

**stack-input.json**:
```json
{
  "target": {
    "spec": {
      "targetCluster": {
        "clusterName": "my-gke-cluster"
      },
      "namespace": {
        "value": "elastic-system"
      }
    }
  },
  "kubernetesElasticOperator": {
    "metadata": {
      "name": "eck-operator",
      "id": "eck-op-prod"
    },
    "spec": {
      "targetCluster": {
        "clusterName": "my-gke-cluster"
      },
      "namespace": {
        "value": "elastic-system"
      },
      "container": {
        "resources": {
          "requests": {
            "cpu": "50m",
            "memory": "100Mi"
          },
          "limits": {
            "cpu": "1000m",
            "memory": "1Gi"
          }
        }
      }
    }
  }
}
```

**Deploy**:
```bash
pulumi stack init prod
pulumi up
```

## Example 2: High-Availability Production

**stack-input.json**:
```json
{
  "target": {
    "spec": {
      "targetCluster": {
        "clusterName": "production-gke-cluster"
      },
      "namespace": {
        "value": "elastic-system"
      }
    }
  },
  "kubernetesElasticOperator": {
    "metadata": {
      "name": "eck-operator-ha",
      "id": "eck-op-prod-ha",
      "org": "platform-team",
      "env": "production"
    },
    "spec": {
      "targetCluster": {
        "clusterName": "production-gke-cluster"
      },
      "namespace": {
        "value": "elastic-system"
      },
      "container": {
        "resources": {
          "requests": {
            "cpu": "200m",
            "memory": "512Mi"
          },
          "limits": {
            "cpu": "2000m",
            "memory": "2Gi"
          }
        }
      }
    }
  }
}
```

## Example 3: Development Environment

**stack-input.json**:
```json
{
  "target": {
    "spec": {
      "targetCluster": {
        "clusterName": "dev-gke-cluster"
      },
      "namespace": {
        "value": "elastic-system"
      }
    }
  },
  "kubernetesElasticOperator": {
    "metadata": {
      "name": "eck-operator-dev",
      "id": "eck-op-dev",
      "env": "development"
    },
    "spec": {
      "targetCluster": {
        "clusterName": "dev-gke-cluster"
      },
      "namespace": {
        "value": "elastic-system"
      },
      "container": {
        "resources": {
          "requests": {
            "cpu": "25m",
            "memory": "64Mi"
          },
          "limits": {
            "cpu": "500m",
            "memory": "512Mi"
          }
        }
      }
    }
  }
}
```

## Verification Commands

```bash
# Check stack outputs
pulumi stack output namespace

# Verify deployment
kubectl get pods -n elastic-system
kubectl get crds | grep elastic

# View operator logs
kubectl logs -n elastic-system -l control-plane=elastic-operator
```

## Cleanup

```bash
pulumi destroy
```

