# DigitalOcean Load Balancer Examples

This document provides practical examples for deploying DigitalOcean Load Balancers using OpenMCF's manifest-based approach.

## Table of Contents

1. [Basic HTTP Load Balancer](#1-basic-http-load-balancer)
2. [Production HTTPS Load Balancer with SSL Termination](#2-production-https-load-balancer-with-ssl-termination)
3. [Multi-Port Load Balancer (HTTP + HTTPS)](#3-multi-port-load-balancer-http--https)
4. [TCP Load Balancer for Database](#4-tcp-load-balancer-for-database)
5. [Load Balancer with Droplet IDs](#5-load-balancer-with-droplet-ids)

---

## 1. Basic HTTP Load Balancer

**Use Case:** Simple HTTP load balancer for development or internal services.

**Features:**
- Tag-based Droplet targeting (dynamic backend discovery)
- HTTP health checks
- Basic configuration

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanLoadBalancer
metadata:
  name: dev-web-lb
spec:
  load_balancer_name: dev-web-lb
  region: nyc3
  
  # VPC for private network communication
  vpc:
    value: "vpc-123456"  # Replace with your VPC ID
  
  # Tag-based targeting - all Droplets with this tag will be added to the pool
  droplet_tag: "web-dev"
  
  # Forward HTTP traffic from port 80 to port 80 on backends
  forwarding_rules:
    - entry_port: 80
      entry_protocol: http
      target_port: 80
      target_protocol: http
  
  # TCP health check on port 80
  health_check:
    port: 80
    protocol: tcp
```

**Deploy:**
```bash
planton pulumi up --manifest dev-web-lb.yaml
```

---

## 2. Production HTTPS Load Balancer with SSL Termination

**Use Case:** Production web application with HTTPS SSL termination at the load balancer.

**Features:**
- HTTPS → HTTP (SSL termination)
- Let's Encrypt certificate integration
- HTTP health checks with custom path
- Sticky sessions enabled
- Custom health check interval

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanLoadBalancer
metadata:
  name: prod-web-lb
  labels:
    environment: production
    app: web-frontend
spec:
  load_balancer_name: prod-web-lb
  region: sfo3
  
  vpc:
    value: "vpc-prod-123"
  
  droplet_tag: "web-prod"
  
  # HTTPS entry with SSL termination, forwarding to HTTP backends
  forwarding_rules:
    - entry_port: 443
      entry_protocol: https
      target_port: 80
      target_protocol: http
      # Use Let's Encrypt certificate name (not ID to avoid renewal issues)
      certificate_name: "prod-example-com-cert"
  
  # HTTP health check with custom endpoint
  health_check:
    port: 80
    protocol: http
    path: "/healthz"
    check_interval_sec: 10
  
  # Enable sticky sessions for session affinity
  enable_sticky_sessions: true
```

**Prerequisites:**
- Upload or generate a certificate in DigitalOcean
- Use the certificate **name**, not ID (IDs change on renewal)

**Deploy:**
```bash
planton pulumi up --manifest prod-web-lb.yaml
```

---

## 3. Multi-Port Load Balancer (HTTP + HTTPS)

**Use Case:** Web application that accepts both HTTP and HTTPS traffic, with HTTP redirecting to HTTPS via application logic.

**Features:**
- Dual forwarding rules (ports 80 and 443)
- SSL termination on port 443
- Both forward to same backend port

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanLoadBalancer
metadata:
  name: dual-port-lb
spec:
  load_balancer_name: dual-port-lb
  region: fra1
  
  vpc:
    value: "vpc-eu-456"
  
  droplet_tag: "web-dual"
  
  forwarding_rules:
    # HTTP traffic (port 80 → 8080)
    - entry_port: 80
      entry_protocol: http
      target_port: 8080
      target_protocol: http
    
    # HTTPS traffic (port 443 → 8080 with SSL termination)
    - entry_port: 443
      entry_protocol: https
      target_port: 8080
      target_protocol: http
      certificate_name: "dual-cert"
  
  health_check:
    port: 8080
    protocol: http
    path: "/health"
```

**Note:** Application should handle HTTP-to-HTTPS redirect logic on port 8080 if needed.

---

## 4. TCP Load Balancer for Database

**Use Case:** Load balance TCP traffic for a MySQL/Galera cluster or PostgreSQL read replicas.

**Features:**
- Layer 4 (TCP) passthrough
- No SSL termination (handled by database)
- TCP health checks

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanLoadBalancer
metadata:
  name: db-lb
spec:
  load_balancer_name: prod-mysql-lb
  region: nyc1
  
  vpc:
    value: "vpc-db-789"
  
  droplet_tag: "mysql-galera"
  
  # TCP passthrough on MySQL port
  forwarding_rules:
    - entry_port: 3306
      entry_protocol: tcp
      target_port: 3306
      target_protocol: tcp
  
  # TCP health check (no path needed)
  health_check:
    port: 3306
    protocol: tcp
```

**Use Cases:**
- MySQL/MariaDB Galera clusters
- PostgreSQL read replicas
- Redis clusters
- Message queues (RabbitMQ, Kafka)

**Deploy:**
```bash
planton pulumi up --manifest db-lb.yaml
```

---

## 5. Load Balancer with Droplet IDs

**Use Case:** Static load balancer for specific Droplets (not recommended for production; use tags instead).

**Features:**
- Explicit Droplet ID targeting
- Useful for testing or manual setups

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanLoadBalancer
metadata:
  name: static-lb
spec:
  load_balancer_name: static-dev-lb
  region: nyc3
  
  vpc:
    value: "vpc-test-111"
  
  # Explicit Droplet IDs (mutually exclusive with droplet_tag)
  droplet_ids:
    - value: "386734086"
    - value: "386734087"
  
  forwarding_rules:
    - entry_port: 80
      entry_protocol: http
      target_port: 8080
      target_protocol: http
  
  health_check:
    port: 8080
    protocol: http
    path: "/"
```

**Drawbacks:**
- Manual management required when adding/removing Droplets
- Doesn't work with autoscaling
- Not suitable for blue-green deployments

**Recommendation:** Use `droplet_tag` instead for dynamic backend discovery.

---

## Advanced Patterns

### Blue-Green Deployment

Use tag-based targeting to enable zero-downtime deployments:

```yaml
# Production load balancer
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanLoadBalancer
metadata:
  name: prod-lb-blue
spec:
  load_balancer_name: prod-web-lb
  region: sfo3
  vpc:
    value: "vpc-prod"
  
  # Initially target "blue" environment
  droplet_tag: "blue"
  
  forwarding_rules:
    - entry_port: 443
      entry_protocol: https
      target_port: 80
      target_protocol: http
      certificate_name: "prod-cert"
  
  health_check:
    port: 80
    protocol: http
    path: "/healthz"
```

**Deployment Process:**
1. Deploy new version to Droplets tagged `green`
2. Test `green` Droplets independently
3. Update load balancer spec to change `droplet_tag` from `blue` to `green`
4. Apply change → instant traffic cutover
5. Keep `blue` Droplets for rollback capability

### Health Check Tuning for High Availability

```yaml
spec:
  health_check:
    port: 80
    protocol: http
    path: "/healthz"
    check_interval_sec: 5  # Check every 5 seconds for faster failover
```

**Best Practices:**
- Use application-level health checks (HTTP with `/health` endpoint)
- Implement health check logic that verifies database connectivity
- Return 200 OK only when application is truly healthy
- Avoid TCP-only checks for critical services (they only verify port is open)

---

## Best Practices Summary

1. **Use tag-based targeting** - Enables autoscaling and blue-green deployments
2. **Use certificate names, not IDs** - Prevents issues with Let's Encrypt renewals
3. **Implement proper health checks** - HTTP checks with application logic, not just TCP
4. **Place in VPC** - Use private network for Droplet communication (cost + security)
5. **Enable sticky sessions when needed** - For session-based applications
6. **Block direct Droplet access** - Use Cloud Firewalls to allow traffic only from LB

---

## Troubleshooting

### 503 Service Unavailable

**Cause:** No healthy backends

**Solutions:**
1. Check health check configuration (correct port, path, protocol)
2. Verify backend application responds to health check requests
3. Check Droplets are tagged correctly (if using tag-based targeting)
4. SSH to a Droplet and `curl http://localhost:<port><path>` to test

### SSL Handshake Failures

**Cause:** Incorrect certificate configuration or PROXY protocol mismatch

**Solutions:**
1. Verify `certificate_name` matches uploaded certificate
2. Don't enable PROXY protocol unless backends support it
3. Use certificate name, not ID

### Terraform State Drift

**Cause:** Droplet IDs change when using tag-based targeting

**Solution:** Use `ignore_changes` in lifecycle block (already configured in module)

---

## Next Steps

- Review [docs/README.md](docs/README.md) for architecture and best practices
- Check [iac/pulumi/README.md](iac/pulumi/README.md) for Pulumi deployment details
- Check [iac/tf/README.md](iac/tf/README.md) for Terraform deployment details
- See [hack/manifest.yaml](hack/manifest.yaml) for a test manifest

