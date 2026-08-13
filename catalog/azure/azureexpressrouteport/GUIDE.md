# Azure ExpressRoute Port -- Operational Guide

Judgment that does not fit in field references.

## This is the expensive one

An ExpressRoute Port bills its full monthly rate -- thousands of
dollars a month at 10 Gbps, several times that at 100 -- from the
moment ARM creates it, whether or not a single cross-connect exists.
There is no cheap tier and no pause switch. Create a port when the
facility order is genuinely in motion, and delete abandoned ones the
day the plan changes.

## Enrollment may precede everything

Some subscriptions must be enrolled for ExpressRoute Direct before ARM
will accept a port create (the create fails with an authorization-
class error, not a quota error). If the create rejects on a fresh
subscription, that is a Microsoft-side enablement conversation --
raise it through your account team or a support request; no manifest
change will fix it.

## The port is step one of a physical workflow

Creating the ARM object provisions nothing physical by itself:

1. Create the port (this component). The per-link outputs (router,
   interface, patch panel, rack, connector) are the letter-of-
   authorization facts.
2. Hand those facts to the colocation facility to order the two
   cross-connects (one per link).
3. Enable the links (`link1.adminEnabled` / `link2.adminEnabled` --
   they start DISABLED) once the facility completes the physical work.
4. Carve circuits from the port (AzureExpressRouteCircuit in Direct
   mode, referencing `express_route_port_id`), then configure peerings
   on the circuits.

Expect facility lead time in days-to-weeks; the ARM object itself
provisions in minutes.

## Encapsulation is a one-way door

DOT1Q vs QINQ is fixed at creation and decides how circuits share the
port: DOT1Q gives each circuit one VLAN tag (fine when you control the
VLAN plan); QINQ stacks an Azure-managed outer tag so overlapping
customer VLAN ranges coexist (the multi-tenant/provider-ish shape).
Changing your mind replaces the port AND every circuit on it.

## MACsec needs three parties aligned

Layer-2 encryption on a link requires: both Key Vault secret IDs (CKN
names the key, CAK is the key -- they travel together), a
USER_ASSIGNED identity on the port with secret GET on both, and the
facility side configured with the same key material and SCI setting.
The spec enforces the first two upfront; the third is out-of-band
coordination -- a MACsec mismatch reads as a dead link, not an error
message.

## Oversubscription is a feature, within reason

Azure lets the aggregate provisioned bandwidth of the port's circuits
exceed the port size (up to 2x) -- ingress above the physical rate is
dropped, egress is shaped. Carving 2x10 Gbps circuits from a 10 Gbps
port for failover topologies is a designed pattern, not a mistake --
just never plan sustained simultaneous saturation.

## Authorizations are credentials, not metadata

Each `authorizations` entry generates a key that lets a circuit in
ANOTHER subscription consume this port's capacity. Treat the
`authorization_keys` output like any credential store: issue one per
consuming team/subscription (named for the consumer), and delete the
entry to revoke access when the consumer is decommissioned.
