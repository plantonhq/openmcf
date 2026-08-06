# High throughput — production

A production-scale connector for serverless fleets that push serious traffic into the VPC. `e2-standard-4` instances carry roughly a 1 Gbps class of throughput each (versus ~200 Mbps for `e2-micro`), and the 3–10 instance band gives the autoscaler real room.

Two capacity levers, different costs to pull: `machineType` changes in place, but DECREASING `minInstances`/`maxInstances` REPLACES the connector — a brief egress outage for every workload using it. Size the band generously up front; increases apply in place.
