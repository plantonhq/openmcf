# Web Tier

The public-subnet posture done statelessly-correct: HTTPS and HTTP in, ephemeral ports in (responses to outbound calls), HTTPS and ephemeral ports out (responses to inbound clients). Both public subnets associate atomically; everything unmatched hits AWS's catch-all deny.
