# Queue-to-Queue with a Filter

The classic pipe: drain one SQS queue into another, keeping only order events and reshaping each message with an input template. Wire both queues and a pipes-trusting role by reference — no polling Lambda, no glue code.
