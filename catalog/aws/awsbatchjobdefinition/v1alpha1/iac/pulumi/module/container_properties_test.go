package module

import (
	"encoding/json"
	"testing"

	awsbatchjobdefinitionv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbatchjobdefinition/v1alpha1"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func literal(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func decode(t *testing.T, document string) map[string]interface{} {
	t.Helper()
	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(document), &decoded); err != nil {
		t.Fatalf("container properties is not valid JSON: %v", err)
	}
	return decoded
}

// The minimal document must carry the image and the modern
// resourceRequirements sizing (string-typed values, per the Batch API), and
// nothing else -- absent optionals must be absent, not null/zero, so the
// provider's semantic-JSON comparison never sees phantom fields.
func TestContainerProperties_Minimal(t *testing.T) {
	spec := &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionSpec{
		Region: "us-west-2",
		Container: &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionContainer{
			Image:     "public.ecr.aws/amazonlinux/amazonlinux:2023",
			Vcpus:     0.25,
			MemoryMib: 512,
		},
	}

	document, err := buildContainerProperties(spec)
	if err != nil {
		t.Fatal(err)
	}
	decoded := decode(t, document)

	if decoded["image"] != "public.ecr.aws/amazonlinux/amazonlinux:2023" {
		t.Errorf("unexpected image: %v", decoded["image"])
	}
	requirements := decoded["resourceRequirements"].([]interface{})
	if len(requirements) != 2 {
		t.Fatalf("expected 2 resource requirements, got %d", len(requirements))
	}
	vcpu := requirements[0].(map[string]interface{})
	if vcpu["type"] != "VCPU" || vcpu["value"] != "0.25" {
		t.Errorf("unexpected VCPU requirement: %v", vcpu)
	}
	memory := requirements[1].(map[string]interface{})
	if memory["type"] != "MEMORY" || memory["value"] != "512" {
		t.Errorf("unexpected MEMORY requirement: %v", memory)
	}
	for _, absent := range []string{"command", "environment", "secrets", "volumes", "linuxParameters", "networkConfiguration"} {
		if _, exists := decoded[absent]; exists {
			t.Errorf("field %q must be absent from a minimal document", absent)
		}
	}
}

func TestContainerProperties_FullSurface(t *testing.T) {
	spec := &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionSpec{
		Region:               "us-west-2",
		PlatformCapabilities: []string{"EC2"},
		Container: &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionContainer{
			Image:         "123456789012.dkr.ecr.us-west-2.amazonaws.com/etl:1.4.2",
			Command:       []string{"python", "run.py", "Ref::dataset"},
			Vcpus:         4,
			MemoryMib:     8192,
			Gpus:          1,
			JobRole:       literal("arn:aws:iam::123456789012:role/etl-job"),
			ExecutionRole: literal("arn:aws:iam::123456789012:role/etl-exec"),
			// Two entries prove the name-sorted deterministic ordering.
			Environment: map[string]string{"STAGE": "prod", "MODE": "batch"},
			Secrets: map[string]string{
				"DB_PASSWORD": "arn:aws:secretsmanager:us-west-2:123456789012:secret:db-AbCdEf",
			},
			LogConfiguration: &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionLogConfiguration{
				LogDriver:     "awslogs",
				Options:       map[string]string{"awslogs-group": "/custom/batch"},
				SecretOptions: map[string]string{"token": "arn:aws:ssm:us-west-2:123456789012:parameter/tok"},
			},
			MountPoints: []*awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionMountPoint{
				{SourceVolume: "shared", ContainerPath: "/data", ReadOnly: true},
			},
			Volumes: []*awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionVolume{
				{
					Name: "shared",
					Efs: &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionEfsVolume{
						FileSystemId:     literal("fs-0123456789abcdef0"),
						AccessPointId:    literal("fsap-0123456789abcdef0"),
						IamAuthorization: true,
					},
				},
				{Name: "scratch", HostPath: "/mnt/scratch"},
			},
			Ulimits: []*awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionUlimit{
				{Name: "nofile", SoftLimit: 8192, HardLimit: 65535},
			},
			LinuxParameters: &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionLinuxParameters{
				InitProcessEnabled:  true,
				SharedMemorySizeMib: 256,
				Tmpfs: []*awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionTmpfs{
					{ContainerPath: "/tmp/scratch", SizeMib: 128, MountOptions: []string{"noexec"}},
				},
				Devices: []*awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionDevice{
					{HostPath: "/dev/xvdf", Permissions: []string{"READ"}},
				},
			},
			Privileged:                     true,
			User:                           "1000:1000",
			ReadonlyRootFilesystem:         true,
			RepositoryCredentialsSecretArn: "arn:aws:secretsmanager:us-west-2:123456789012:secret:reg-AbCdEf",
		},
	}

	document, err := buildContainerProperties(spec)
	if err != nil {
		t.Fatal(err)
	}
	decoded := decode(t, document)

	requirements := decoded["resourceRequirements"].([]interface{})
	if len(requirements) != 3 {
		t.Fatalf("expected VCPU+MEMORY+GPU requirements, got %d", len(requirements))
	}
	gpu := requirements[2].(map[string]interface{})
	if gpu["type"] != "GPU" || gpu["value"] != "1" {
		t.Errorf("unexpected GPU requirement: %v", gpu)
	}

	environment := decoded["environment"].([]interface{})
	first := environment[0].(map[string]interface{})
	if first["name"] != "MODE" {
		t.Errorf("environment must be sorted by name; got first entry %v", first)
	}

	volumes := decoded["volumes"].([]interface{})
	efsVolume := volumes[0].(map[string]interface{})["efsVolumeConfiguration"].(map[string]interface{})
	if efsVolume["transitEncryption"] != "ENABLED" {
		t.Errorf("transit encryption must always be ENABLED, got %v", efsVolume["transitEncryption"])
	}
	authorization := efsVolume["authorizationConfig"].(map[string]interface{})
	if authorization["iam"] != "ENABLED" || authorization["accessPointId"] != "fsap-0123456789abcdef0" {
		t.Errorf("unexpected authorization config: %v", authorization)
	}
	hostVolume := volumes[1].(map[string]interface{})["host"].(map[string]interface{})
	if hostVolume["sourcePath"] != "/mnt/scratch" {
		t.Errorf("unexpected host volume: %v", hostVolume)
	}

	credentials := decoded["repositoryCredentials"].(map[string]interface{})
	if credentials["credentialsParameter"] != "arn:aws:secretsmanager:us-west-2:123456789012:secret:reg-AbCdEf" {
		t.Errorf("unexpected repository credentials: %v", credentials)
	}
}

// Fargate-only knobs render as their dedicated API sub-objects.
func TestContainerProperties_FargateKnobs(t *testing.T) {
	spec := &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionSpec{
		Region:               "us-west-2",
		PlatformCapabilities: []string{"FARGATE"},
		Container: &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionContainer{
			Image:                  "public.ecr.aws/amazonlinux/amazonlinux:2023",
			Vcpus:                  1,
			MemoryMib:              2048,
			ExecutionRole:          literal("arn:aws:iam::123456789012:role/exec"),
			FargatePlatformVersion: "1.4.0",
			AssignPublicIp:         true,
			EphemeralStorageGib:    50,
			RuntimePlatform: &awsbatchjobdefinitionv1alpha1.AwsBatchJobDefinitionRuntimePlatform{
				CpuArchitecture: "ARM64",
			},
		},
	}

	document, err := buildContainerProperties(spec)
	if err != nil {
		t.Fatal(err)
	}
	decoded := decode(t, document)

	if decoded["fargatePlatformConfiguration"].(map[string]interface{})["platformVersion"] != "1.4.0" {
		t.Errorf("unexpected fargatePlatformConfiguration: %v", decoded["fargatePlatformConfiguration"])
	}
	if decoded["networkConfiguration"].(map[string]interface{})["assignPublicIp"] != "ENABLED" {
		t.Errorf("unexpected networkConfiguration: %v", decoded["networkConfiguration"])
	}
	if decoded["ephemeralStorage"].(map[string]interface{})["sizeInGiB"] != float64(50) {
		t.Errorf("unexpected ephemeralStorage: %v", decoded["ephemeralStorage"])
	}
	if decoded["runtimePlatform"].(map[string]interface{})["cpuArchitecture"] != "ARM64" {
		t.Errorf("unexpected runtimePlatform: %v", decoded["runtimePlatform"])
	}
}
