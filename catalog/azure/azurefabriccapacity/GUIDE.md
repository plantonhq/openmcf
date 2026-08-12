# Microsoft Fabric Capacity -- Operational Guide

Judgment calls that matter when you run Fabric capacities in production.

## The meter runs from minute one

A Fabric capacity is the rare catalog kind that costs real money at REST: the hour meter starts when the deploy finishes and stops only at delete, whether or not any workspace uses it. Pay-as-you-go F2 is roughly a third of a dollar per hour; F2048 scales that by a thousand. Treat capacity lifetime as an operational decision -- create development capacities when needed and delete them after, and never let a proof-of-concept capacity outlive its meeting.

## Size with the dial, not with foresight

The F-SKU moves up AND down in place with no downtime and no replacement, so the right starting size is the smallest one that works (F2 for almost all development). Watch the Fabric Capacity Metrics app for throttling and smoothing-debt, and step the SKU when evidence -- not anticipation -- says so. The two thresholds that matter on the way up: F64 unlocks Copilot and lets Power BI free-license users consume shared content.

## The capacity is ARM's whole Fabric story

Everything inside Fabric -- workspaces, lakehouses, warehouses, pipelines -- is managed OUTSIDE azurerm: Microsoft ships a dedicated `fabric` Terraform provider for it, and the Fabric portal/APIs own the rest. Plan the operational boundary accordingly: this kind anchors billing and compute in your Azure estate; workspace governance is a Fabric-side discipline with its own tooling and its own administrators (the `administration_members` you declare here).

## Administrators are declared here, exercised there

The administration members list is the bridge between the two worlds: identities declared on the ARM resource, exercised in the Fabric admin experience (assigning workspaces, managing settings). Keep it to a small platform group -- the spec requires at least one at all times, because a capacity with no administrator is unmanageable from the Fabric side without an Azure-side edit.
