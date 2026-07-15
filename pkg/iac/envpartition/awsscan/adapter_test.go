package awsscan

import (
	"reflect"
	"testing"
)

// The property documents here mirror the exact decoded-JSON shapes a Cloud
// Control scan delivers (tags as [{Key, Value}], attachments as enriched
// lists) -- the same shapes the recorded eval fixtures pin.
func TestAdaptReadsNameTagAndContainers(t *testing.T) {
	resource := Adapt("AWS::EC2::Subnet", "subnet-1", map[string]any{
		"VpcId":     "vpc-1",
		"CidrBlock": "10.0.1.0/24",
		"Tags": []any{
			map[string]any{"Key": "Name", "Value": "orders-prod-app-subnet"},
			map[string]any{"Key": "team", "Value": "orders"},
		},
	})
	if resource.Name != "orders-prod-app-subnet" {
		t.Fatalf("Name = %q, want the Name tag", resource.Name)
	}
	if resource.Tags["team"] != "orders" {
		t.Fatalf("Tags = %v, want the flattened tag map", resource.Tags)
	}
	if !reflect.DeepEqual(resource.Containers, []string{"vpc-1"}) {
		t.Fatalf("Containers = %v, want [vpc-1]", resource.Containers)
	}
}

func TestAdaptFallsBackToNameShapedProperties(t *testing.T) {
	queue := Adapt("AWS::SQS::Queue", "https://sqs.us-west-2.amazonaws.com/1/orders-prod-events", map[string]any{
		"QueueName": "orders-prod-events",
	})
	if queue.Name != "orders-prod-events" {
		t.Fatalf("queue Name = %q, want QueueName", queue.Name)
	}

	// No Name tag and no name-shaped property: Name stays empty. The
	// opaque identifier must NEVER become the name surface -- a random id
	// can embed a misleading environment token.
	nat := Adapt("AWS::EC2::NatGateway", "nat-0prodfix123", map[string]any{
		"SubnetId": "subnet-1", "VpcId": "vpc-1",
	})
	if nat.Name != "" {
		t.Fatalf("NAT Name = %q, want empty (identifiers are not name surfaces)", nat.Name)
	}
	if !reflect.DeepEqual(nat.Containers, []string{"vpc-1", "subnet-1"}) {
		t.Fatalf("NAT Containers = %v, want [vpc-1 subnet-1]", nat.Containers)
	}
}

func TestAdaptReadsInternetGatewayAttachments(t *testing.T) {
	igw := Adapt("AWS::EC2::InternetGateway", "igw-1", map[string]any{
		"Attachments": []any{
			map[string]any{"VpcId": "vpc-1", "State": "attached"},
		},
	})
	if !reflect.DeepEqual(igw.Containers, []string{"vpc-1"}) {
		t.Fatalf("IGW Containers = %v, want [vpc-1]", igw.Containers)
	}
}

func TestAdaptTagBeatsNameShapedProperty(t *testing.T) {
	// When the owners tagged a Name, that is the human-authored surface;
	// the type's own name property is the fallback, not a competitor.
	table := Adapt("AWS::DynamoDB::Table", "orders-prod-table", map[string]any{
		"TableName": "orders-prod-table",
		"Tags": []any{
			map[string]any{"Key": "Name", "Value": "orders-prod-table-display"},
		},
	})
	if table.Name != "orders-prod-table-display" {
		t.Fatalf("Name = %q, want the Name tag to win", table.Name)
	}
}
