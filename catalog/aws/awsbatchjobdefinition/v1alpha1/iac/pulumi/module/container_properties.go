package module

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	awsbatchjobdefinitionv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbatchjobdefinition/v1alpha1"
)

// buildContainerProperties renders the spec's structured container block as
// the Batch RegisterJobDefinition containerProperties JSON document.
//
// The provider takes containerProperties as an opaque JSON string in both
// engines and compares it semantically, so the ONE requirement beyond
// correctness is that this document and the Terraform module's jsonencode()
// output are semantically identical for the same spec -- field names here
// are the Batch API's camelCase names, and map-derived lists are sorted so
// the document is deterministic across runs.
func buildContainerProperties(spec *awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionSpec) (string, error) {
	container := spec.Container

	properties := map[string]interface{}{
		"image": container.Image,
	}

	if len(container.Command) > 0 {
		properties["command"] = container.Command
	}

	// Sizing goes through resourceRequirements (the modern shape; the
	// top-level vcpus/memory API fields are deprecated). Values are API
	// strings, which also keeps 64-bit-int JSON hazards out of the
	// document by construction.
	requirements := []map[string]string{
		{"type": "VCPU", "value": strconv.FormatFloat(container.Vcpus, 'f', -1, 64)},
		{"type": "MEMORY", "value": strconv.FormatInt(int64(container.MemoryMib), 10)},
	}
	if container.Gpus > 0 {
		requirements = append(requirements, map[string]string{
			"type": "GPU", "value": strconv.FormatInt(int64(container.Gpus), 10),
		})
	}
	properties["resourceRequirements"] = requirements

	if container.JobRole.GetValue() != "" {
		properties["jobRoleArn"] = container.JobRole.GetValue()
	}
	if container.ExecutionRole.GetValue() != "" {
		properties["executionRoleArn"] = container.ExecutionRole.GetValue()
	}

	if len(container.Environment) > 0 {
		properties["environment"] = sortedNameValueList(container.Environment, "value")
	}
	if len(container.Secrets) > 0 {
		properties["secrets"] = sortedNameValueList(container.Secrets, "valueFrom")
	}

	if container.LogConfiguration != nil {
		logConfiguration := map[string]interface{}{
			"logDriver": container.LogConfiguration.LogDriver,
		}
		if len(container.LogConfiguration.Options) > 0 {
			logConfiguration["options"] = sortedStringMap(container.LogConfiguration.Options)
		}
		if len(container.LogConfiguration.SecretOptions) > 0 {
			logConfiguration["secretOptions"] = sortedNameValueList(container.LogConfiguration.SecretOptions, "valueFrom")
		}
		properties["logConfiguration"] = logConfiguration
	}

	if len(container.MountPoints) > 0 {
		mountPoints := make([]map[string]interface{}, 0, len(container.MountPoints))
		for _, mount := range container.MountPoints {
			mountPoints = append(mountPoints, map[string]interface{}{
				"sourceVolume":  mount.SourceVolume,
				"containerPath": mount.ContainerPath,
				"readOnly":      mount.ReadOnly,
			})
		}
		properties["mountPoints"] = mountPoints
	}

	if len(container.Volumes) > 0 {
		volumes := make([]map[string]interface{}, 0, len(container.Volumes))
		for _, volume := range container.Volumes {
			entry := map[string]interface{}{"name": volume.Name}
			if volume.Efs != nil {
				efsConfiguration := map[string]interface{}{
					"fileSystemId": volume.Efs.FileSystemId.GetValue(),
					// Transit encryption is always on: AWS requires it for
					// access points and IAM auth, and there is no good
					// reason to mount EFS unencrypted in transit.
					"transitEncryption": "ENABLED",
				}
				if volume.Efs.RootDirectory != "" {
					efsConfiguration["rootDirectory"] = volume.Efs.RootDirectory
				}
				if volume.Efs.AccessPointId.GetValue() != "" || volume.Efs.IamAuthorization {
					authorization := map[string]interface{}{}
					if volume.Efs.AccessPointId.GetValue() != "" {
						authorization["accessPointId"] = volume.Efs.AccessPointId.GetValue()
					}
					if volume.Efs.IamAuthorization {
						authorization["iam"] = "ENABLED"
					}
					efsConfiguration["authorizationConfig"] = authorization
				}
				entry["efsVolumeConfiguration"] = efsConfiguration
			}
			if volume.HostPath != "" {
				entry["host"] = map[string]interface{}{"sourcePath": volume.HostPath}
			}
			volumes = append(volumes, entry)
		}
		properties["volumes"] = volumes
	}

	if len(container.Ulimits) > 0 {
		ulimits := make([]map[string]interface{}, 0, len(container.Ulimits))
		for _, ulimit := range container.Ulimits {
			ulimits = append(ulimits, map[string]interface{}{
				"name":      ulimit.Name,
				"softLimit": ulimit.SoftLimit,
				"hardLimit": ulimit.HardLimit,
			})
		}
		properties["ulimits"] = ulimits
	}

	if container.LinuxParameters != nil {
		linux := map[string]interface{}{}
		if container.LinuxParameters.InitProcessEnabled {
			linux["initProcessEnabled"] = true
		}
		if container.LinuxParameters.SharedMemorySizeMib > 0 {
			linux["sharedMemorySize"] = container.LinuxParameters.SharedMemorySizeMib
		}
		if container.LinuxParameters.MaxSwapMib > 0 {
			linux["maxSwap"] = container.LinuxParameters.MaxSwapMib
		}
		if container.LinuxParameters.Swappiness > 0 {
			linux["swappiness"] = container.LinuxParameters.Swappiness
		}
		if len(container.LinuxParameters.Tmpfs) > 0 {
			tmpfsList := make([]map[string]interface{}, 0, len(container.LinuxParameters.Tmpfs))
			for _, tmpfs := range container.LinuxParameters.Tmpfs {
				tmpfsEntry := map[string]interface{}{
					"containerPath": tmpfs.ContainerPath,
					"size":          tmpfs.SizeMib,
				}
				if len(tmpfs.MountOptions) > 0 {
					tmpfsEntry["mountOptions"] = tmpfs.MountOptions
				}
				tmpfsList = append(tmpfsList, tmpfsEntry)
			}
			linux["tmpfs"] = tmpfsList
		}
		if len(container.LinuxParameters.Devices) > 0 {
			devices := make([]map[string]interface{}, 0, len(container.LinuxParameters.Devices))
			for _, device := range container.LinuxParameters.Devices {
				deviceEntry := map[string]interface{}{"hostPath": device.HostPath}
				if device.ContainerPath != "" {
					deviceEntry["containerPath"] = device.ContainerPath
				}
				if len(device.Permissions) > 0 {
					deviceEntry["permissions"] = device.Permissions
				}
				devices = append(devices, deviceEntry)
			}
			linux["devices"] = devices
		}
		if len(linux) > 0 {
			properties["linuxParameters"] = linux
		}
	}

	if container.Privileged {
		properties["privileged"] = true
	}
	if container.User != "" {
		properties["user"] = container.User
	}
	if container.ReadonlyRootFilesystem {
		properties["readonlyRootFilesystem"] = true
	}
	if container.RepositoryCredentialsSecretArn != "" {
		properties["repositoryCredentials"] = map[string]interface{}{
			"credentialsParameter": container.RepositoryCredentialsSecretArn,
		}
	}

	if container.RuntimePlatform != nil {
		runtimePlatform := map[string]interface{}{}
		if container.RuntimePlatform.CpuArchitecture != "" {
			runtimePlatform["cpuArchitecture"] = container.RuntimePlatform.CpuArchitecture
		}
		if container.RuntimePlatform.OperatingSystemFamily != "" {
			runtimePlatform["operatingSystemFamily"] = container.RuntimePlatform.OperatingSystemFamily
		}
		if len(runtimePlatform) > 0 {
			properties["runtimePlatform"] = runtimePlatform
		}
	}
	if container.FargatePlatformVersion != "" {
		properties["fargatePlatformConfiguration"] = map[string]interface{}{
			"platformVersion": container.FargatePlatformVersion,
		}
	}
	if container.AssignPublicIp {
		properties["networkConfiguration"] = map[string]interface{}{
			"assignPublicIp": "ENABLED",
		}
	}
	if container.EphemeralStorageGib > 0 {
		properties["ephemeralStorage"] = map[string]interface{}{
			"sizeInGiB": container.EphemeralStorageGib,
		}
	}

	encoded, err := json.Marshal(properties)
	if err != nil {
		return "", errors.Wrap(err, "encode container properties")
	}
	return string(encoded), nil
}

// sortedNameValueList renders a proto map as the Batch API's list-of-objects
// shape ({"name": k, "<valueKey>": v}), sorted by name for a deterministic
// document.
func sortedNameValueList(entries map[string]string, valueKey string) []map[string]string {
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	list := make([]map[string]string, 0, len(keys))
	for _, key := range keys {
		list = append(list, map[string]string{
			"name":   key,
			valueKey: entries[key],
		})
	}
	return list
}

// sortedStringMap copies a map so json.Marshal emits it with sorted keys
// (Go's encoder sorts map keys, so a plain copy suffices; the helper exists
// to make the determinism intent explicit at call sites).
func sortedStringMap(entries map[string]string) map[string]string {
	copied := make(map[string]string, len(entries))
	for key, value := range entries {
		copied[key] = value
	}
	return copied
}
