# Isolated Domain

The canonical segmentation building block: a routing domain whose members see shared services but not the rest of the network, with a blackhole guaranteeing the isolation.

## When to Use

- One route table per isolation zone (prod, non-prod, partner) on a segmented hub

## What It Configures

- **`associations`** — the domain's own spokes look up their traffic here (each must have its default-table association off)
- **`propagations`** — only the shared-services attachment advertises into this domain, so shared services are reachable but sibling spokes are not
- **A blackhole route** — the other environment's CIDR can never be reached from this domain, even if a propagation is added by mistake later (statics beat propagated routes)

## What to Customize

- Replace `<aws-region>`, the hub, and the attachment references with your resources
- Add an egress default route (`0.0.0.0/0` via an egress VPC attachment) if this domain reaches the internet through a central egress VPC
