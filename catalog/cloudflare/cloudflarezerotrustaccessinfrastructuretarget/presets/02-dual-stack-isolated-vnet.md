# Dual-stack target in an isolated virtual network

The overlapping-CIDR shape: both address families pinned into a specific virtual network (a literal UUID here; reference a `CloudflareZeroTrustTunnelVirtualNetwork` resource's output in real estates). Use this when two sites reuse the same private ranges -- the virtual network disambiguates which 10.0.10.5 this target means.
