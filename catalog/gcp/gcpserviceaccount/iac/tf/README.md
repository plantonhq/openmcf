# Terraform Module to Deploy AWS DynamoDB table

```shell
planton tofu init --manifest e2e/manifest.yaml --backend-type s3 \
  --backend-config="bucket=planton-tf-state-backend" \
  --backend-config="dynamodb_table=planton-tf-state-backend-lock" \
  --backend-config="region=ap-south-2" \
  --backend-config="key=planton/gcp-stacks/test-gcp-service-account.tfstate"
```

```shell
planton tofu plan --manifest e2e/manifest.yaml
```

```shell
planton tofu apply --manifest e2e/manifest.yaml --auto-approve
```

```shell
planton tofu destroy --manifest e2e/manifest.yaml --auto-approve
```
