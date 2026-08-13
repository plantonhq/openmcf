# Azure Data Factory Integration Runtime -- Operational Guide

Judgment for operating integration runtimes in production, earned from how Azure actually behaves.

## Know which flavor bills, and when

The three flavors have three billing models. The data-flow compute bills per vCore-hour only while a cluster is up: every run, plus every warm minute `time_to_live_min` buys. The SSIS runtime bills per node-hour from START to STOP -- creating it is free, and this component only creates it; treat starting it as turning on a meter (`number_of_nodes` x node size, around the clock until stopped). The self-hosted flavor is free on Azure's side -- the compute is yours. When an SSIS bill surprises someone, the answer is almost always a runtime left started over a weekend.

## Warm pools trade money for latency, deliberately

A data flow against a cold runtime pays several minutes of cluster startup. `time_to_live_min` keeps the cluster warm between runs, and `cleanup_enabled: false` preserves the pool -- back-to-back flows then start in seconds. Set the TTL to your inter-run gap, not longer: warm minutes bill exactly like run minutes, and AutoResolve regions cannot share warm pools across regions.

## The managed virtual network switch is a create-time decision

`virtual_network_enabled` requires the factory itself to have its managed virtual network enabled, and Azure rejects the runtime -- with its own error -- when it does not. The switch is also ForceNew: flipping it later replaces the runtime. Decide network posture when the FACTORY is designed; a data flow that must reach a private endpoint needs the whole chain (factory managed VNet -> runtime inside it -> managed private endpoint) in place.

## One name, one runtime -- switching flavors replaces it

All three flavors live in one factory-scoped namespace ({factory_id}/integrationRuntimes/{name}). Changing the variant block replaces the object at the same ARM address, and every linked service or activity referencing the name picks up the new engine immediately. Rename rather than reshape when anything still depends on the old flavor.

## Treat the self-hosted authorization keys like passwords

Azure returns the keys readable; this component surfaces them as sensitive outputs, and they are full join credentials -- any machine holding one can register as YOUR runtime and see the data that flows through it. Wire them into installers by reference, rotate by re-keying (the secondary key exists exactly so agents migrate before the primary rotates), and never paste them into manifests. A linked registration (RBAC authorization) issues no keys of its own -- the primary runtime's keys govern.

## SSIS custom setup: prefer express, reference Key Vault

The express custom setup covers the common node preparation (environment variables, PowerShell, licensed components, cmdkey credentials) without maintaining a setup-script container, and its password and license fields each have a Key-Vault-reference alternative -- use those; when both forms are set Azure receives the INLINE value, silently winning over the reference. Reserve the script container for installers express setup cannot express, and remember every setup runs on every node start, so slow installers stretch every scale-out.

## The SSIS catalog outlives the runtime

`catalog_info` creates the SSISDB database on YOUR Azure SQL server when the runtime first starts; deleting the runtime does not delete SSISDB. The server endpoint and tier are the runtime's contract, but the database's lifecycle (backups, failover via `dual_standby_pair_name`, cost) is the SQL server's own story -- plan it there.
