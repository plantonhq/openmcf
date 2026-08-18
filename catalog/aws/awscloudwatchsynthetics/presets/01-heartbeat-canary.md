# Heartbeat Canary

A five-minute health probe with one retry: the script reads `TARGET_URL` and fails the run when the endpoint does. Trimmed success retention keeps artifact cost down; `deleteLambda` keeps teardown clean. Stage the zip (`nodejs/node_modules/heartbeat.js`) through an AwsS3ObjectSet or your pipeline.
