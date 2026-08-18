# ML Training Scratch

Hot shard storage next to the GPUs: a One Zone bucket in the training cluster's zone, force-destroyable because everything in it is reconstructible from the regional source-of-truth bucket. The derived name lands as `training-scratch--use1-az4--x-s3`.
