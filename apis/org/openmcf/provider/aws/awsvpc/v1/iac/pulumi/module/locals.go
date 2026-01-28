package module

import (
	"fmt"
	"net"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"

	awsvpcv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/aws/awsvpc/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SubnetName string
type SubnetCidr string
type AvailabilityZone string

type Locals struct {
	AwsVpc             *awsvpcv1.AwsVpc
	AwsTags            map[string]string
	PrivateAzSubnetMap map[AvailabilityZone]map[SubnetName]SubnetCidr
	PublicAzSubnetMap  map[AvailabilityZone]map[SubnetName]SubnetCidr
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsvpcv1.AwsVpcStackInput) (*Locals, error) {
	locals := &Locals{}
	locals.AwsVpc = stackInput.Target

	_, vpcCidrBlock, err := net.ParseCIDR(locals.AwsVpc.Spec.VpcCidr)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid vpc cidr %s", locals.AwsVpc.Spec.VpcCidr)
	}
	vpcMaskSize, _ := vpcCidrBlock.Mask.Size()
	if vpcMaskSize > int(locals.AwsVpc.Spec.SubnetSize) {
		return nil, errors.Errorf("spec.subnetSize /%d cannot be bigger than the VPC /%d",
			locals.AwsVpc.Spec.SubnetSize, vpcMaskSize)
	}

	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsVpc.Metadata.Org,
		awstagkeys.Environment:  locals.AwsVpc.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsVpc.String(),
		awstagkeys.ResourceId:   locals.AwsVpc.Metadata.Id,
	}

	locals.PrivateAzSubnetMap = GetPrivateAzSubnetMap(locals.AwsVpc)
	locals.PublicAzSubnetMap = getPublicAzSubnetMap(locals.AwsVpc)

	return locals, nil
}

func GetPrivateAzSubnetMap(awsVpc *awsvpcv1.AwsVpc) map[AvailabilityZone]map[SubnetName]SubnetCidr {
	privateAzSubnetMap := make(map[AvailabilityZone]map[SubnetName]SubnetCidr, 0)

	// Calculate total number of public subnets (they come first in CIDR allocation)
	publicSubnetsCount := len(awsVpc.Spec.AvailabilityZones) * int(awsVpc.Spec.SubnetsPerAvailabilityZone)

	for azIndex, az := range awsVpc.Spec.AvailabilityZones {
		for subnetIndex := 0; subnetIndex < int(awsVpc.Spec.SubnetsPerAvailabilityZone); subnetIndex++ {
			privateSubnetName := fmt.Sprintf("private-subnet-%s-%d", az, subnetIndex)
			// Calculate subnet index: offset by public subnets + position within private subnets
			globalSubnetIndex := publicSubnetsCount + azIndex*int(awsVpc.Spec.SubnetsPerAvailabilityZone) + subnetIndex
			privateSubnetCidr := calculateSubnetCidr(awsVpc.Spec.VpcCidr, int(awsVpc.Spec.SubnetSize), globalSubnetIndex)

			if privateAzSubnetMap[AvailabilityZone(az)] == nil {
				privateAzSubnetMap[AvailabilityZone(az)] = make(map[SubnetName]SubnetCidr)
			}

			privateAzSubnetMap[AvailabilityZone(az)][SubnetName(privateSubnetName)] = SubnetCidr(privateSubnetCidr)
		}
	}
	return privateAzSubnetMap
}

func getPublicAzSubnetMap(awsVpc *awsvpcv1.AwsVpc) map[AvailabilityZone]map[SubnetName]SubnetCidr {
	publicAzSubnetMap := make(map[AvailabilityZone]map[SubnetName]SubnetCidr, 0)

	for azIndex, az := range awsVpc.Spec.AvailabilityZones {
		for subnetIndex := 0; subnetIndex < int(awsVpc.Spec.SubnetsPerAvailabilityZone); subnetIndex++ {
			publicSubnetName := fmt.Sprintf("public-subnet-%s-%d", az, subnetIndex)
			// Calculate global subnet index for public subnets
			globalSubnetIndex := azIndex*int(awsVpc.Spec.SubnetsPerAvailabilityZone) + subnetIndex
			publicSubnetCidr := calculateSubnetCidr(awsVpc.Spec.VpcCidr, int(awsVpc.Spec.SubnetSize), globalSubnetIndex)

			if publicAzSubnetMap[AvailabilityZone(az)] == nil {
				publicAzSubnetMap[AvailabilityZone(az)] = make(map[SubnetName]SubnetCidr)
			}
			publicAzSubnetMap[AvailabilityZone(az)][SubnetName(publicSubnetName)] = SubnetCidr(publicSubnetCidr)
		}
	}
	return publicAzSubnetMap
}

// calculateSubnetCidr calculates a subnet CIDR based on the VPC CIDR, subnet mask size, and subnet index
// For example: VPC CIDR "10.0.0.0/16" with subnet size /24 and index 0 returns "10.0.0.0/24"
func calculateSubnetCidr(vpcCidr string, subnetMaskSize int, subnetIndex int) string {
	_, vpcCidrBlock, err := net.ParseCIDR(vpcCidr)
	if err != nil {
		// This should have been validated in initializeLocals, but handle gracefully
		return fmt.Sprintf("10.0.%d.0/%d", subnetIndex, subnetMaskSize)
	}

	// Get the base IP address as a 32-bit integer
	baseIP := vpcCidrBlock.IP.To4()
	if baseIP == nil {
		// IPv6 not supported yet
		return fmt.Sprintf("10.0.%d.0/%d", subnetIndex, subnetMaskSize)
	}

	// Convert base IP to uint32
	baseIPInt := uint32(baseIP[0])<<24 | uint32(baseIP[1])<<16 | uint32(baseIP[2])<<8 | uint32(baseIP[3])

	// Calculate the number of IP addresses per subnet
	// For /24 subnet, this is 2^(32-24) = 256 addresses
	vpcMaskSize, _ := vpcCidrBlock.Mask.Size()
	ipsPerSubnet := uint32(1) << (uint(subnetMaskSize) - uint(vpcMaskSize))

	// Calculate the subnet IP by adding the offset
	subnetIPInt := baseIPInt + uint32(subnetIndex)*ipsPerSubnet

	// Convert back to IP address
	subnetIP := net.IPv4(
		byte(subnetIPInt>>24),
		byte(subnetIPInt>>16),
		byte(subnetIPInt>>8),
		byte(subnetIPInt),
	)

	return fmt.Sprintf("%s/%d", subnetIP.String(), subnetMaskSize)
}

func getSortedAzKeys(azSubnetMap map[AvailabilityZone]map[SubnetName]SubnetCidr) []string {
	keys := make([]string, 0, len(azSubnetMap))
	for k := range azSubnetMap {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	return keys
}

func getSortedSubnetNameKeys(subnetMap map[SubnetName]SubnetCidr) []string {
	keys := make([]string, 0, len(subnetMap))
	for k := range subnetMap {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	return keys
}
