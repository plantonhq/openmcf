# Multi-Region Active-Active Half

One half of an active-active pair: this us-east-1 cluster names its us-west-2 peer by reference and us-east-2 as the witness (a third region that arbitrates failover — it runs no queries). Deploy the mirror-image instance in us-west-2 naming this one; both write, both read, synchronously. Pair fresh clusters only — the pairing window is create-time, and the customer-managed key keeps encryption under your governance in each region.
