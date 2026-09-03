package stackinput

import (
	"fmt"
	"os"

	"buf.build/go/protovalidate"
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput/fieldsextractor"
	"github.com/plantonhq/planton/pkg/protobufyaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	PulumiConfigKey   = "planton:stack-input"
	FilePathEnvVar    = "STACK_INPUT_YAML_FILE"
	YamlContentEnvVar = "STACK_INPUT_YAML"

	// OperationEnvVar tells the module program which engine operation spawned
	// it. `pulumi destroy --run-program` re-runs the program only so delete
	// hooks can fire; every resource in state is deleted whatever the program
	// registers. A module step that can FAIL for reasons unrelated to what is
	// being deleted (deriving CRDs from a chart version that was never
	// published, a repository that is unreachable) must therefore not fail
	// the program during a destroy, or the user could never tear down a stack
	// whose manifest is currently wrong. Set by every Pulumi spawner in this
	// repository; absent means an update or preview.
	OperationEnvVar  = "PLANTON_IAC_OPERATION"
	OperationDestroy = "destroy"
)

// IsDestroy reports whether the program runs for a `pulumi destroy`.
func IsDestroy() bool {
	return os.Getenv(OperationEnvVar) == OperationDestroy
}

func LoadStackInput(ctx *pulumi.Context, stackInput proto.Message) error {
	stackInputString, ok := ctx.GetConfig(PulumiConfigKey)
	var jsonBytes, stackInputYamlBytes []byte
	var err error

	if !ok {
		yamlContent := os.Getenv(YamlContentEnvVar)
		if yamlContent != "" {
			stackInputYamlBytes = []byte(yamlContent)
		} else {
			stackInputFilePath := os.Getenv(FilePathEnvVar)
			if stackInputFilePath == "" {
				return errors.Errorf("stack-input not found in pulumi config %s or in %s environment variable",
					PulumiConfigKey, FilePathEnvVar)
			}
			stackInputYamlBytes, err = os.ReadFile(stackInputFilePath)
			if err != nil {
				return errors.Wrap(err, "failed to read input file")
			}
		}

		jsonBytes, err = protobufyaml.YAMLToJSON(stackInputYamlBytes)
		if err != nil {
			return errors.Wrap(err, "failed to load yaml to json")
		}
	} else {
		jsonBytes, err = protobufyaml.YAMLToJSON([]byte(stackInputString))
		if err != nil {
			return errors.Wrap(err, "failed to load yaml to json")
		}
	}

	if err := protojson.Unmarshal(jsonBytes, stackInput); err != nil {
		return errors.Wrap(err, "failed to load json into proto message")
	}

	targetSpec, err := fieldsextractor.ExtractApiResourceSpecField(stackInput)
	if err != nil {
		return errors.Wrap(err, "failed to extract api resource spec field")
	}

	v, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages((*targetSpec).Interface()),
	)
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	if err = v.Validate((*targetSpec).Interface()); err != nil {
		return errors.Errorf("%s", err)
	}
	return nil
}
