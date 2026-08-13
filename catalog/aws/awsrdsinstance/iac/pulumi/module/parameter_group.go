package module

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// parameterGroup provisions the DB parameter group managed for inline
// spec.parameters -- pure glue (a named parameter list), which is why
// it stays inside this module instead of being its own node. The
// family is derived from engine + engine_version (both CEL-required
// with inline parameters); a family change replaces the group and
// re-associates the instance in the same update. Returns nil when the
// spec brings no inline parameters.
func parameterGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*rds.ParameterGroup, error) {
	spec := locals.AwsRdsInstance.Spec
	if len(spec.Parameters) == 0 {
		return nil, nil
	}

	parameters := rds.ParameterGroupParameterArray{}
	for _, parameter := range spec.Parameters {
		parameterArgs := &rds.ParameterGroupParameterArgs{
			Name:  pulumi.String(parameter.Name),
			Value: pulumi.String(parameter.Value),
		}
		if parameter.ApplyMethod != "" {
			parameterArgs.ApplyMethod = pulumi.String(parameter.ApplyMethod)
		}
		parameters = append(parameters, parameterArgs)
	}

	createdParameterGroup, err := rds.NewParameterGroup(ctx, "parameter-group",
		&rds.ParameterGroupArgs{
			Name:       pulumi.String(locals.InstanceIdentifier),
			Family:     pulumi.String(parameterGroupFamily(spec.Engine, spec.EngineVersion)),
			Parameters: parameters,
			Tags:       pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create DB parameter group")
	}
	return createdParameterGroup, nil
}

// parameterGroupFamily derives the parameter-group family from
// engine + engine_version, per AWS's own family naming (mirrors the
// Terraform module's locals key-for-key):
//
//	postgres     16.4          -> postgres16        (major)
//	mysql        8.0.39        -> mysql8.0          (major.minor)
//	mariadb      10.11.8       -> mariadb10.11      (major.minor)
//	oracle-ee    19.0.0.0.ru.. -> oracle-ee-19      (engine-major)
//	sqlserver-se 16.00.4085..  -> sqlserver-se-16.0 (engine-major.minor;
//	                              the numeric parse collapses "00" to 0)
//
// db2-* engines take the sqlserver arm (db2-ae-11.5).
func parameterGroupFamily(engine, engineVersion string) string {
	segments := strings.Split(engineVersion, ".")
	major := segments[0]
	minor := 0
	if len(segments) > 1 {
		if parsed, err := strconv.Atoi(segments[1]); err == nil {
			minor = parsed
		}
	}

	switch {
	case engine == "postgres":
		return "postgres" + major
	case engine == "mysql", engine == "mariadb":
		return engine + major + "." + strconv.Itoa(minor)
	case strings.HasPrefix(engine, "oracle-"):
		return engine + "-" + major
	default:
		return engine + "-" + major + "." + strconv.Itoa(minor)
	}
}
