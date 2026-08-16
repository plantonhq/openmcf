# Organization Log Sink

The central-account half of cross-account log aggregation: a named Kinesis-backed destination whose access policy lists the producer accounts allowed to point subscription filters at it. Producers then subscribe their log groups to this destination's ARN — every account's logs land in one stream.
