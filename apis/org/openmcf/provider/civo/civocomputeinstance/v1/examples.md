# CivoComputeInstance Examples

Real-world examples for common Civo compute instance deployment scenarios.

## Table of Contents

1. [Minimal Development Instance](#1-minimal-development-instance)
2. [Production Web Server with Firewall](#2-production-web-server-with-firewall)
3. [Database Server with Persistent Storage](#3-database-server-with-persistent-storage)
4. [Bastion Host with Reserved IP](#4-bastion-host-with-reserved-ip)
5. [Application Server with Cloud-Init](#5-application-server-with-cloud-init)
6. [Multi-Instance Setup (Web + Database)](#6-multi-instance-setup-web--database)
7. [High-Availability Load Balanced Setup](#7-high-availability-load-balanced-setup)

---

## 1. Minimal Development Instance

**Use Case**: Quick dev environment for testing.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoComputeInstance
metadata:
  name: dev-sandbox
spec:
  instanceName: dev-sandbox
  region: NYC1
  size: g3.small
  image: ubuntu-jammy
  network:
    value: default-network-id
```

**Key Points**:
- 1 vCPU, 2GB RAM (~$10.86/month)
- No firewall (uses Civo default - all ports open, dev only!)
- No reserved IP (ephemeral)
- Boots in < 60 seconds

**Post-Deploy**:
```bash
# SSH to instance (use Civo dashboard to get IP)
ssh root@<public-ip>
```

---

## 2. Production Web Server with Firewall

**Use Case**: HTTPS-enabled web server.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoComputeInstance
metadata:
  name: web-prod-01
  description: Production web server
spec:
  instanceName: web-prod-01
  region: LON1
  size: g3.large
  image: ubuntu-jammy
  network:
    value: prod-network-id
  sshKeyIds:
    - team-ssh-key-id
  firewallIds:
    - value: web-firewall-id  # Allows 80/443 from 0.0.0.0/0, SSH from office
  reservedIpId:
    value: web-static-ip-id
  tags:
    - env:prod
    - service:web
    - criticality:high
  userData: |
    #!/bin/bash
    apt-get update
    apt-get upgrade -y
    apt-get install -y nginx certbot python3-certbot-nginx
    systemctl enable nginx
    systemctl start nginx
```

**Key Points**:
- 4 vCPU, 8GB RAM (~$55/month)
- Custom firewall for security
- Reserved IP for stable DNS
- Cloud-init installs Nginx + Certbot

---

## 3. Database Server with Persistent Storage

**Use Case**: PostgreSQL database with data volume.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoComputeInstance
metadata:
  name: postgres-db-01
spec:
  instanceName: postgres-db-01
  region: NYC1
  size: g3.xlarge
  image: ubuntu-jammy
  network:
    value: db-network-id
  sshKeyIds:
    - admin-ssh-key-id
  firewallIds:
    - value: postgres-firewall-id  # Allows 5432 from app network only
  volumeIds:
    - value: postgres-data-volume-id  # 100GB persistent volume
  reservedIpId:
    value: db-static-ip-id
  tags:
    - env:prod
    - service:database
    - db-type:postgresql
  userData: |
    #!/bin/bash
    # Mount persistent volume
    mkfs.ext4 /dev/vdb
    mkdir -p /var/lib/postgresql
    mount /dev/vdb /var/lib/postgresql
    echo "/dev/vdb /var/lib/postgresql ext4 defaults 0 0" >> /etc/fstab
    
    # Install PostgreSQL
    apt-get update
    apt-get install -y postgresql-14
    systemctl enable postgresql
    systemctl start postgresql
```

**Key Points**:
- 6 vCPU, 16GB RAM (~$87/month)
- Attached 100GB volume for data persistence
- Firewall restricts access to app network only
- Reserved IP for stable connection string

---

## 4. Bastion Host with Reserved IP

**Use Case**: Secure jump host for accessing private instances.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoComputeInstance
metadata:
  name: bastion-prod
spec:
  instanceName: bastion-prod
  region: FRA1
  size: g3.small
  image: ubuntu-jammy
  network:
    value: prod-network-id
  sshKeyIds:
    - admin-key-1
    - admin-key-2
  firewallIds:
    - value: bastion-firewall-id  # SSH only from office IPs
  reservedIpId:
    value: bastion-static-ip
  tags:
    - env:prod
    - role:bastion
  userData: |
    #!/bin/bash
    # Harden SSH
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
    sed -i 's/#PermitRootLogin yes/PermitRootLogin prohibit-password/' /etc/ssh/sshd_config
    systemctl restart sshd
    
    # Install fail2ban
    apt-get update
    apt-get install -y fail2ban
    systemctl enable fail2ban
```

**Key Points**:
- Small instance (minimal cost)
- Multiple admin SSH keys
- Firewall restricts to office IPs
- Reserved IP for stable access

---

## 5. Application Server with Cloud-Init

**Use Case**: Dockerized application with automated deployment.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoComputeInstance
metadata:
  name: app-prod-01
spec:
  instanceName: app-prod-01
  region: NYC1
  size: g3.medium
  image: ubuntu-jammy
  network:
    value: app-network-id
  sshKeyIds:
    - deploy-key-id
  firewallIds:
    - value: app-firewall-id  # 8080 from LB, SSH from bastion
  tags:
    - env:prod
    - service:api
  userData: |
    #!/bin/bash
    set -euo pipefail
    
    # Install Docker
    apt-get update
    apt-get install -y docker.io docker-compose
    systemctl enable docker
    systemctl start docker
    
    # Deploy application
    docker login -u deploy -p ${REGISTRY_PASSWORD}
    docker pull myregistry/api:v1.2.3
    docker run -d --name api \
      --restart=always \
      -p 8080:8080 \
      -e DATABASE_URL=${DATABASE_URL} \
      myregistry/api:v1.2.3
    
    # Setup monitoring
    curl -s https://www.civo.com/civostatsd.sh | bash
```

**Key Points**:
- Docker-based deployment
- Automated application pull and start
- Environment variables from secrets
- Civo monitoring agent

---

## 6. Multi-Instance Setup (Web + Database)

**Use Case**: Web tier + database tier in isolated networks.

### Database Instance

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoComputeInstance
metadata:
  name: db-backend
spec:
  instanceName: db-backend
  region: LON1
  size: g3.xlarge
  image: ubuntu-jammy
  network:
    value: backend-network-id
  sshKeyIds:
    - admin-key
  firewallIds:
    - value: db-firewall-id  # 5432 from web network only
  volumeIds:
    - value: db-volume-id
  tags:
    - env:prod
    - tier:database
```

### Web Instance

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoComputeInstance
metadata:
  name: web-frontend
spec:
  instanceName: web-frontend
  region: LON1
  size: g3.large
  image: ubuntu-jammy
  network:
    value: frontend-network-id
  sshKeyIds:
    - admin-key
  firewallIds:
    - value: web-firewall-id  # 80/443 from internet, SSH from bastion
  reservedIpId:
    value: web-public-ip
  tags:
    - env:prod
    - tier:web
  userData: |
    #!/bin/bash
    apt-get update && apt-get install -y nginx
    # Configure Nginx to proxy to backend
    cat > /etc/nginx/sites-available/default <<'NGINX'
    server {
      listen 80;
      location / {
        proxy_pass http://db-backend.internal:8080;
      }
    }
    NGINX
    systemctl reload nginx
```

**Key Points**:
- Isolated networks (frontend/backend)
- Database firewall only allows web network
- Web tier has public IP
- Both use reserved IPs for stability

---

## 7. High-Availability Multi-Instance Setup

**Use Case**: Multiple web servers for redundancy and availability (load balancing via Kubernetes or DNS).

### Web Instance 1

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoComputeInstance
metadata:
  name: web-01
spec:
  instanceName: web-01
  region: FRA1
  size: g3.large
  image: ubuntu-jammy
  network:
    value: web-network-id
  sshKeyIds:
    - deploy-key
  firewallIds:
    - value: web-firewall  # HTTP/HTTPS access
  tags:
    - env:prod
    - service:web
    - instance-group:web-cluster
```

### Web Instance 2

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoComputeInstance
metadata:
  name: web-02
spec:
  instanceName: web-02
  region: FRA1
  size: g3.large
  image: ubuntu-jammy
  network:
    value: web-network-id
  sshKeyIds:
    - deploy-key
  firewallIds:
    - value: web-firewall
  tags:
    - env:prod
    - service:web
    - instance-group:web-cluster
```

**Key Points**:
- Two identical instances for redundancy
- Use DNS round-robin or Kubernetes Services for traffic distribution
- Tags enable grouping for management
- Each instance has public IP for direct access if needed

**Note**: For production load balancing, consider:
- Civo Kubernetes with Service type LoadBalancer (automatically provisions LB)
- External load balancer solutions (HAProxy, Nginx)
- DNS-based failover with health checks

---

## Best Practices Summary

1. **Naming**: Use consistent patterns (`<service>-<env>-<number>`)
2. **Sizes**: Start small, scale up based on metrics
3. **Images**: Use LTS releases (ubuntu-jammy, debian-11)
4. **Networks**: Isolate tiers (web, app, database)
5. **Firewalls**: Always create custom (never use default)
6. **SSH**: Keys only, never passwords
7. **Reserved IPs**: For production services only
8. **Volumes**: For stateful data (databases, uploads)
9. **Cloud-Init**: Install monitoring, security updates, application
10. **Tags**: For cost tracking, access control, organization

---

## Testing Your Instance

After provisioning:

```bash
# Get instance details
civo instance show <instance-id>

# Test SSH connectivity
ssh -i ~/.ssh/civo_key root@<public-ip>

# Check cloud-init completed
ssh root@<public-ip> "cloud-init status --wait"

# View cloud-init logs
ssh root@<public-ip> "tail -100 /var/log/cloud-init-output.log"

# Check services running
ssh root@<public-ip> "systemctl status nginx docker postgresql"
```

---

## Related Documentation

- **API Reference**: [README.md](./README.md)
- **Research**: [docs/README.md](./docs/README.md)
- **Pulumi Module**: [iac/pulumi/](./iac/pulumi/)
- **Civo Compute Docs**: [civo.com/docs/compute](https://www.civo.com/docs/compute)
