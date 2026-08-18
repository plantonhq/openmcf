# Organization-Shared Runtime Layer

A platform team's layer published for the whole AWS Organization: the wildcard principal scoped down by `organizationId` lets every member account attach it, and `skipDestroy` keeps replaced versions alive so consumer functions never break mid-rollout. The `sourceCodeHash` makes pipeline republishes declarative — set it from `base64(sha256(zip))` in your build.
