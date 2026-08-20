# Scheduled Fargate Task

A rate schedule launching a Fargate task every 15 minutes, joined to an existing group by name, with a five-minute flexible window so a fleet of these never stampedes the cluster at once. The task definition, subnets, and security group all wire by reference.
