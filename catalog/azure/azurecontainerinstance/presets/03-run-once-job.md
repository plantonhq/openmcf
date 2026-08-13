# Run-Once Job

This preset runs a batch job to completion: an init container fetches configuration into a shared scratch volume, the main container consumes it, and the group terminates when the work is done -- billed only for the seconds it ran.

## When to Use

- Data loads, schema migrations, report generation -- anything with a beginning and an end
- Pipeline steps that need more CPU/memory than the pipeline runner offers
- One-off maintenance commands you want isolated, logged, and disposable

## Key Configuration Choices

- **`restartPolicy: Never`** -- the container's exit ENDS the job (the group shows Terminated); use "OnFailure" instead when non-zero exits should retry
- **`ipAddressType: None` + `priority: Spot`** -- a job needs no group IP, and "None" is the only posture Spot accepts; Spot trades evictability for a steep discount, which rerunnable jobs are the right home for (delete the priority line if yours cannot tolerate eviction)
- **The shared `emptyDir`** -- the same volume name in the init and main containers is ONE scratch volume; the init step seeds it, the job reads it
- **`secureEnvironmentVariables`** -- hidden from the portal and API reads; Azure never returns the values, so imports re-supply them from the manifest
- **Read the exit state before deleting** -- there is no rerun; a "Never" group that failed stays Terminated with its logs until you delete it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console |
| `<your-registry>` | Your registry's login-server prefix | The Azure Container Registry's overview page |
| `<your-database-password>` | The job's database credential | Your secret store (prefer a pipeline secret reference over a literal) |

## Related Presets

- **Public Web Container** -- the always-on service shape
- **Private VNet Worker** -- the subnet-joined worker shape
