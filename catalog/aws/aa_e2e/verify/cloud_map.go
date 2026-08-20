package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	sdtypes "github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"
	pkgerrors "github.com/pkg/errors"
)

// cloudMapNamespaceVerifier verifies an AwsCloudMapNamespace with its
// folded services and statically registered instances: the namespace
// via GetNamespace, each service via GetService over the service_ids
// map, and each registration via GetInstance over the
// instance_service_ids map (keyed "{service_name}//{instance_id}" -
// the instance id is the key's second segment).
type cloudMapNamespaceVerifier struct{}

func cloudMapClientForRegion(cfg aws.Config, region string) *servicediscovery.Client {
	return servicediscovery.NewFromConfig(cfg, func(o *servicediscovery.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

func isCloudMapNotFound(err error) bool {
	var namespaceNotFound *sdtypes.NamespaceNotFound
	var serviceNotFound *sdtypes.ServiceNotFound
	var instanceNotFound *sdtypes.InstanceNotFound
	return errors.As(err, &namespaceNotFound) || errors.As(err, &serviceNotFound) || errors.As(err, &instanceNotFound)
}

func (*cloudMapNamespaceVerifier) IDOutputKey() string { return "namespace_id" }

func (v *cloudMapNamespaceVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := cloudMapClientForRegion(cfg, region).GetNamespace(ctx, &servicediscovery.GetNamespaceInput{
		Id: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudmapnamespace verify-exists failed for %q", id)
	}
	return nil
}

func (v *cloudMapNamespaceVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := cloudMapClientForRegion(cfg, region).GetNamespace(ctx, &servicediscovery.GetNamespaceInput{
		Id: aws.String(id),
	})
	if err == nil {
		return pkgerrors.Errorf("awscloudmapnamespace %q still exists after destroy", id)
	}
	if isCloudMapNotFound(err) {
		return nil
	}
	return pkgerrors.Wrapf(err, "awscloudmapnamespace verify-absent failed for %q", id)
}

func (v *cloudMapNamespaceVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	namespaceId := stringOutput(outputs, "namespace_id")
	if namespaceId == "" {
		return pkgerrors.New("awscloudmapnamespace verify-exists: no namespace_id in outputs")
	}
	if err := v.VerifyExists(ctx, cfg, namespaceId, region); err != nil {
		return err
	}
	client := cloudMapClientForRegion(cfg, region)
	for serviceName, serviceId := range stringMapOutput(outputs["service_ids"]) {
		if _, err := client.GetService(ctx, &servicediscovery.GetServiceInput{
			Id: aws.String(serviceId),
		}); err != nil {
			return pkgerrors.Wrapf(err, "awscloudmapnamespace service %q (%s) verify-exists failed", serviceName, serviceId)
		}
	}
	for pairKey, serviceId := range stringMapOutput(outputs["instance_service_ids"]) {
		instanceId := instanceIdFromPairKey(pairKey)
		if instanceId == "" {
			return pkgerrors.Errorf("awscloudmapnamespace registration key %q has no instance segment", pairKey)
		}
		if _, err := client.GetInstance(ctx, &servicediscovery.GetInstanceInput{
			ServiceId:  aws.String(serviceId),
			InstanceId: aws.String(instanceId),
		}); err != nil {
			return pkgerrors.Wrapf(err, "awscloudmapnamespace registration %q verify-exists failed", pairKey)
		}
	}
	return nil
}

func (v *cloudMapNamespaceVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	client := cloudMapClientForRegion(cfg, region)
	for pairKey, serviceId := range stringMapOutput(outputs["instance_service_ids"]) {
		instanceId := instanceIdFromPairKey(pairKey)
		if instanceId == "" {
			continue
		}
		_, err := client.GetInstance(ctx, &servicediscovery.GetInstanceInput{
			ServiceId:  aws.String(serviceId),
			InstanceId: aws.String(instanceId),
		})
		if err == nil {
			return pkgerrors.Errorf("awscloudmapnamespace registration %q still exists after destroy", pairKey)
		}
		if !isCloudMapNotFound(err) {
			return pkgerrors.Wrapf(err, "awscloudmapnamespace registration %q verify-absent failed", pairKey)
		}
	}
	for serviceName, serviceId := range stringMapOutput(outputs["service_ids"]) {
		_, err := client.GetService(ctx, &servicediscovery.GetServiceInput{
			Id: aws.String(serviceId),
		})
		if err == nil {
			return pkgerrors.Errorf("awscloudmapnamespace service %q still exists after destroy", serviceName)
		}
		if !isCloudMapNotFound(err) {
			return pkgerrors.Wrapf(err, "awscloudmapnamespace service %q verify-absent failed", serviceName)
		}
	}
	if namespaceId := stringOutput(outputs, "namespace_id"); namespaceId != "" {
		return v.VerifyAbsent(ctx, cfg, namespaceId, region)
	}
	return nil
}

// instanceIdFromPairKey extracts the instance id from a
// "{service_name}//{instance_id}" map key.
func instanceIdFromPairKey(pairKey string) string {
	parts := strings.SplitN(pairKey, "//", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
