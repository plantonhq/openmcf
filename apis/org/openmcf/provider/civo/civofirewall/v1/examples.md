# Civo Firewall Examples

This document provides real-world examples of firewall configurations using the `CivoFirewall` resource. Each example includes the complete YAML manifest, explanation, and security considerations.

## Table of Contents

1. [Web Server Firewall (HTTP/HTTPS)](#1-web-server-firewall-httphttps)
2. [Database Firewall (PostgreSQL)](#2-database-firewall-postgresql)
3. [Kubernetes Cluster Firewall](#3-kubernetes-cluster-firewall)
4. [Bastion Host Firewall](#4-bastion-host-firewall)
5. [Multi-Tier Application Firewall](#5-multi-tier-application-firewall)

---

## 1. Web Server Firewall (HTTP/HTTPS)

**Use Case:** Public-facing web server that needs to accept HTTP and HTTPS traffic from anywhere, but restrict SSH access to office IP.

**Requirements:**
- HTTP (port 80) from anywhere
- HTTPS (port 443) from anywhere
- SSH (port 22) from office IP only
- Allow ICMP (ping) for monitoring

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoFirewall
metadata:
  name: web-server-firewall
spec:
  name: web-server-fw
  networkId:
    value: "your-network-id-here"
  inboundRules:
    # SSH - Restricted to office IP
    - protocol: tcp
      portRange: "22"
      cidrs:
        - "203.0.113.10/32"  # Replace with your office IP
      action: allow
      label: SSH from office
    
    # HTTP - Public access
    - protocol: tcp
      portRange: "80"
      cidrs:
        - "0.0.0.0/0"
      action: allow
      label: HTTP from anywhere
    
    # HTTPS - Public access
    - protocol: tcp
      portRange: "443"
      cidrs:
        - "0.0.0.0/0"
      action: allow
      label: HTTPS from anywhere
    
    # ICMP - Allow ping for monitoring
    - protocol: icmp
      portRange: ""
      cidrs:
        - "0.0.0.0/0"
      action: allow
      label: Allow ping
  
  tags:
    - web-server
    - public
```

**After deployment:**

```bash
# Apply the firewall
planton apply -f web-server-firewall.yaml

# Test HTTP access
curl http://your-server-ip

# Test HTTPS access
curl https://your-server-ip

# Test SSH (should only work from office IP)
ssh user@your-server-ip
```

**Security notes:**
- SSH is restricted to a single office IP (`/32` CIDR)
- If your office IP changes, update the firewall before losing access
- Consider adding a VPN IP as a backup SSH access point
- Web traffic (80/443) is intentionally public

---

## 2. Database Firewall (PostgreSQL)

**Use Case:** PostgreSQL database that should only accept connections from application tier servers, with SSH access from a bastion host.

**Requirements:**
- PostgreSQL (port 5432) from app tier subnet only
- SSH (port 22) from bastion host only
- No public internet access

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoFirewall
metadata:
  name: database-firewall
spec:
  name: postgres-db-fw
  networkId:
    value: "your-network-id-here"
  inboundRules:
    # PostgreSQL - Only from application tier
    - protocol: tcp
      portRange: "5432"
      cidrs:
        - "10.0.1.0/24"  # App tier subnet
      action: allow
      label: PostgreSQL from app tier
    
    # SSH - Only from bastion host
    - protocol: tcp
      portRange: "22"
      cidrs:
        - "10.0.0.10/32"  # Bastion host IP
      action: allow
      label: SSH from bastion
    
    # Monitoring agent (optional)
    - protocol: tcp
      portRange: "9100"
      cidrs:
        - "10.0.2.0/24"  # Monitoring subnet
      action: allow
      label: Node exporter for Prometheus
  
  tags:
    - database
    - postgresql
    - internal
```

**Connection workflow:**

```bash
# 1. SSH to bastion host
ssh user@bastion-public-ip

# 2. From bastion, SSH to database
ssh user@10.0.0.20  # DB private IP

# 3. Application connects directly (no SSH tunnel needed)
# App tier can connect: psql -h 10.0.0.20 -U dbuser -d mydb
```

**Security notes:**
- Database port is never exposed to the internet
- Only the app tier subnet can reach PostgreSQL
- SSH requires going through bastion host first (defense in depth)
- No outbound rules = all outbound traffic allowed by default (for OS updates, etc.)

**Multi-database variant:**

For multiple database types on different ports:

```yaml
inboundRules:
  - protocol: tcp
    portRange: "5432"
    cidrs: ["10.0.1.0/24"]
    label: PostgreSQL from app tier
  
  - protocol: tcp
    portRange: "6379"
    cidrs: ["10.0.1.0/24"]
    label: Redis from app tier
  
  - protocol: tcp
    portRange: "27017"
    cidrs: ["10.0.1.0/24"]
    label: MongoDB from app tier
```

---

## 3. Kubernetes Cluster Firewall

**Use Case:** Firewall for Kubernetes cluster nodes, allowing API server access, NodePort services, and ingress traffic.

**Requirements:**
- Kubernetes API (port 6443) from anywhere (for kubectl access)
- NodePort range (30000-32767) from anywhere (for NodePort services)
- HTTP/HTTPS (80/443) from anywhere (for ingress controller)
- SSH from office IP for node maintenance

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoFirewall
metadata:
  name: k8s-cluster-firewall
spec:
  name: k8s-cluster-fw
  networkId:
    value: "your-network-id-here"
  inboundRules:
    # Kubernetes API server
    - protocol: tcp
      portRange: "6443"
      cidrs:
        - "0.0.0.0/0"
      action: allow
      label: Kubernetes API server
    
    # Kubernetes NodePort range (for LoadBalancer/NodePort services)
    - protocol: tcp
      portRange: "30000-32767"
      cidrs:
        - "0.0.0.0/0"
      action: allow
      label: Kubernetes NodePort range
    
    # HTTP ingress
    - protocol: tcp
      portRange: "80"
      cidrs:
        - "0.0.0.0/0"
      action: allow
      label: HTTP ingress controller
    
    # HTTPS ingress
    - protocol: tcp
      portRange: "443"
      cidrs:
        - "0.0.0.0/0"
      action: allow
      label: HTTPS ingress controller
    
    # SSH for node maintenance
    - protocol: tcp
      portRange: "22"
      cidrs:
        - "203.0.113.10/32"  # Office IP
      action: allow
      label: SSH for node maintenance
    
    # ICMP for health checks
    - protocol: icmp
      portRange: ""
      cidrs:
        - "0.0.0.0/0"
      action: allow
      label: Allow ping
  
  tags:
    - kubernetes
    - cluster
```

**Testing the cluster:**

```bash
# 1. Get cluster kubeconfig
civo kubernetes config k8s-cluster --save

# 2. Test API access (should work from anywhere)
kubectl get nodes

# 3. Deploy a test service
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --type=NodePort --port=80

# 4. Access via NodePort
curl http://<node-ip>:<node-port>
```

**Security notes:**
- API server is publicly accessible (common for managed K8s)
- Consider restricting API to VPN/office IPs for private clusters
- NodePort range is intentionally broad (required for service discovery)
- Ingress controller ports (80/443) route to backend services via K8s networking

**Production hardening:**

For production, restrict API access:

```yaml
- protocol: tcp
  portRange: "6443"
  cidrs:
    - "203.0.113.0/24"  # Office network
    - "198.51.100.0/24"  # VPN range
    - "10.0.0.0/16"      # Internal cluster network
  label: Kubernetes API (restricted)
```

---

## 4. Bastion Host Firewall

**Use Case:** Jump server (bastion host) for accessing private instances securely.

**Requirements:**
- SSH (port 22) from trusted IPs only
- All outbound traffic allowed (to reach private instances)
- No other inbound services

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoFirewall
metadata:
  name: bastion-firewall
spec:
  name: bastion-fw
  networkId:
    value: "your-network-id-here"
  inboundRules:
    # SSH - Multiple trusted IPs
    - protocol: tcp
      portRange: "22"
      cidrs:
        - "203.0.113.10/32"  # Office IP 1
        - "198.51.100.20/32"  # Office IP 2
        - "192.0.2.0/24"      # VPN range
      action: allow
      label: SSH from trusted sources
    
    # ICMP for monitoring
    - protocol: icmp
      portRange: ""
      cidrs:
        - "203.0.113.0/24"  # Office network for ping tests
      action: allow
      label: Ping from office
  
  # No outbound rules = all outbound allowed (default)
  # This allows SSH from bastion to private instances
  
  tags:
    - bastion
    - jump-host
```

**Usage workflow:**

```bash
# 1. SSH to bastion (public IP)
ssh -i ~/.ssh/bastion-key user@bastion-public-ip

# 2. From bastion, SSH to private instances
ssh -i ~/.ssh/private-key user@10.0.1.10  # Private DB
ssh -i ~/.ssh/private-key user@10.0.1.20  # Private app server
```

**Advanced: SSH Agent Forwarding**

To avoid storing private keys on the bastion:

```bash
# On local machine, add key to agent
ssh-add ~/.ssh/private-key

# SSH with agent forwarding
ssh -A user@bastion-public-ip

# From bastion, SSH without needing key on bastion
ssh user@10.0.1.10  # Uses forwarded agent
```

**Security notes:**
- Bastion is the **only** publicly accessible SSH gateway
- All private instances have firewalls that only allow SSH from bastion
- Use SSH key authentication, never passwords
- Enable SSH session logging for audit trails
- Consider fail2ban to block brute force attempts

**Hardening tips:**

1. Change SSH port to non-standard (security through obscurity):

```yaml
- protocol: tcp
  portRange: "2222"  # Non-standard SSH port
  cidrs: ["203.0.113.10/32"]
  label: SSH on non-standard port
```

2. Disable SSH password authentication in `/etc/ssh/sshd_config`:

```bash
PasswordAuthentication no
PubkeyAuthentication yes
```

3. Use MFA (multi-factor authentication) with Google Authenticator

---

## 5. Multi-Tier Application Firewall

**Use Case:** Complex application with separate web, app, and database tiers, each with appropriate firewall rules.

### Web Tier Firewall

Public-facing load balancers:

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoFirewall
metadata:
  name: web-tier-firewall
spec:
  name: web-tier-fw
  networkId:
    value: "your-network-id-here"
  inboundRules:
    - protocol: tcp
      portRange: "80"
      cidrs: ["0.0.0.0/0"]
      label: HTTP from internet
    
    - protocol: tcp
      portRange: "443"
      cidrs: ["0.0.0.0/0"]
      label: HTTPS from internet
    
    - protocol: tcp
      portRange: "22"
      cidrs: ["10.0.0.10/32"]  # Bastion only
      label: SSH from bastion
  
  tags:
    - web-tier
    - frontend
```

### App Tier Firewall

Application servers (not public):

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoFirewall
metadata:
  name: app-tier-firewall
spec:
  name: app-tier-fw
  networkId:
    value: "your-network-id-here"
  inboundRules:
    # Application port - Only from web tier
    - protocol: tcp
      portRange: "8080"
      cidrs:
        - "10.0.1.0/24"  # Web tier subnet
      action: allow
      label: App server from web tier
    
    # Redis - Only from same app tier
    - protocol: tcp
      portRange: "6379"
      cidrs:
        - "10.0.2.0/24"  # App tier subnet (self)
      action: allow
      label: Redis within app tier
    
    # SSH - Only from bastion
    - protocol: tcp
      portRange: "22"
      cidrs:
        - "10.0.0.10/32"  # Bastion host
      action: allow
      label: SSH from bastion
  
  tags:
    - app-tier
    - backend
```

### Database Tier Firewall

Database servers (most restricted):

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoFirewall
metadata:
  name: db-tier-firewall
spec:
  name: db-tier-fw
  networkId:
    value: "your-network-id-here"
  inboundRules:
    # PostgreSQL - Only from app tier
    - protocol: tcp
      portRange: "5432"
      cidrs:
        - "10.0.2.0/24"  # App tier subnet
      action: allow
      label: PostgreSQL from app tier
    
    # SSH - Only from bastion
    - protocol: tcp
      portRange: "22"
      cidrs:
        - "10.0.0.10/32"  # Bastion host
      action: allow
      label: SSH from bastion
  
  tags:
    - db-tier
    - database
```

**Network diagram:**

```
Internet
   ↓ (80/443)
Web Tier (10.0.1.0/24) - web-tier-fw
   ↓ (8080)
App Tier (10.0.2.0/24) - app-tier-fw
   ↓ (5432)
DB Tier (10.0.3.0/24) - db-tier-fw
   ↑ (SSH)
Bastion (10.0.0.10) - bastion-fw
```

**Deployment order:**

```bash
# 1. Create bastion firewall first
planton apply -f bastion-firewall.yaml

# 2. Create tier firewalls
planton apply -f web-tier-firewall.yaml
planton apply -f app-tier-firewall.yaml
planton apply -f db-tier-firewall.yaml

# 3. Verify firewalls
planton get civofirewalls
```

**Security benefits:**

- **Defense in depth**: Compromise of web tier doesn't give direct DB access
- **Segmentation**: Each tier isolated from others
- **Least privilege**: Only necessary ports open between tiers
- **Audit trail**: All firewall changes tracked in Git

---

## Advanced Patterns

### Deny Rules (Blacklisting)

While less common, you can explicitly deny traffic:

```yaml
inboundRules:
  # Block specific bad actor IP
  - protocol: tcp
    portRange: "80"
    cidrs:
      - "192.0.2.50/32"
    action: deny
    label: Block malicious IP
  
  # Allow all other HTTP
  - protocol: tcp
    portRange: "80"
    cidrs:
      - "0.0.0.0/0"
    action: allow
    label: HTTP from everyone else
```

**Note**: Civo processes rules in order, so place deny rules before allow rules.

### Temporary Access (Maintenance Window)

For temporary access during maintenance:

```yaml
inboundRules:
  # Temporary contractor access - Remove after project
  - protocol: tcp
    portRange: "22"
    cidrs:
      - "203.0.113.99/32"
    action: allow
    label: "TEMPORARY: Contractor access - Remove after 2025-12-31"
```

Add a calendar reminder to remove the rule.

### Monitoring and Logging Ports

Common monitoring agent ports:

```yaml
inboundRules:
  # Prometheus Node Exporter
  - protocol: tcp
    portRange: "9100"
    cidrs: ["10.0.10.0/24"]
    label: Node exporter for Prometheus
  
  # Grafana Agent
  - protocol: tcp
    portRange: "12345"
    cidrs: ["10.0.10.0/24"]
    label: Grafana agent
  
  # Elastic Filebeat
  - protocol: tcp
    portRange: "5044"
    cidrs: ["10.0.10.0/24"]
    label: Filebeat for ELK
```

---

## Testing Your Firewalls

### Connection Testing

```bash
# Test TCP port (should timeout if blocked)
nc -zv <instance-ip> 80

# Test with timeout
timeout 5 nc -zv <instance-ip> 22

# Test ICMP
ping -c 3 <instance-ip>
```

### Port Scanning (Authorized Testing Only)

```bash
# Scan common ports
nmap -Pn <instance-ip>

# Scan specific port range
nmap -p 80,443,22 <instance-ip>
```

**Warning**: Only scan systems you own or have permission to test.

### Verify Firewall Assignment

```bash
# Via Civo CLI
civo instance show <instance-name> | grep firewall

# Via OpenMCF
planton get civofirewalls/your-firewall -o yaml
```

---

## Troubleshooting

### Rule Not Working

1. **Check syntax**:
   - Protocol must be lowercase: `tcp`, not `TCP`
   - Port range format: `80` or `8000-9000`, not `80-`
   - CIDR notation: `203.0.113.10/32`, not `203.0.113.10`

2. **Verify rule order** (for deny rules):
   - Deny rules must come before allow rules
   - Civo processes rules sequentially

3. **Wait for propagation**:
   - Changes can take 10-30 seconds
   - Test again after a brief wait

### Locked Out of Instance

If you accidentally remove SSH access:

1. Use Civo web console for emergency access
2. Fix firewall via CLI or dashboard
3. Test SSH from allowed IP before closing console

**Prevention**: Always test from a second SSH session before closing your current one.

### Performance Issues

Too many firewall rules can cause issues:

- **Limit**: Aim for <50 rules per firewall
- **Consolidate**: Use port ranges instead of individual rules
- **Split**: Create multiple firewalls for different functions

---

## Additional Resources

- [Main README](README.md) - Component overview and quick start
- [Research Documentation](docs/README.md) - Deep dive into deployment methods and design decisions
- [Pulumi Module](iac/pulumi/README.md) - Direct Pulumi usage
- [Civo Firewall API](https://www.civo.com/api/firewalls) - Official API documentation

---

## Need Help?

- Check the [Troubleshooting section](README.md#troubleshooting) in the main README
- Open an issue on [GitHub](https://github.com/plantonhq/openmcf/issues)
- Contact Civo support: support@civo.com

