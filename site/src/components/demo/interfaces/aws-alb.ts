import { ApiResourceMetadata } from './common';

export interface StringValueOrRef {
  value?: string;
  valueFrom?: {
    kind?: string;
    env?: string;
    name?: string;
    fieldPath?: string;
  };
}

export interface AwsAlbDns {
  enabled: boolean;
  route53ZoneId?: StringValueOrRef;
  hostnames: string[];
}

export interface AwsAlbSsl {
  enabled: boolean;
  certificateArn?: StringValueOrRef;
}

export interface AwsAlbSpec {
  subnets: StringValueOrRef[];
  securityGroups: StringValueOrRef[];
  internal: boolean;
  deleteProtectionEnabled: boolean;
  idleTimeoutSeconds: number;
  dns?: AwsAlbDns;
  ssl?: AwsAlbSsl;
}

export interface AwsAlbStackOutputs {
  loadBalancerArn: string;
  loadBalancerName: string;
  loadBalancerDnsName: string;
  loadBalancerHostedZoneId: string;
}

export interface AwsAlbStatus {
  outputs?: AwsAlbStackOutputs;
}

export interface AwsAlb {
  apiVersion: string;
  kind: string;
  metadata: ApiResourceMetadata;
  spec: AwsAlbSpec;
  status?: AwsAlbStatus;
}
