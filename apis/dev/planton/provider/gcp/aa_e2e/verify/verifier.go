// Package verify checks that GCP resources created by an E2E scenario exist
// after DEPLOY and are gone after DESTROY. Each component has its own verifier
// because GCP verification is service-specific (iam serviceAccounts.get for a
// service account, cloudresourcemanager getIamPolicy for a grant, ...). All
// verifiers run against the same ambient ADC chain the deploy used, so a
// verification failure reflects real cloud state.
//
// Unlike providers whose resources carry one opaque id, GCP identifiers are
// frequently compound (an IAM grant is a project+role+member tuple), so
// verifiers receive the component's full string-ified stack outputs.
package verify

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/api/alloydb/v1"
	"google.golang.org/api/bigquery/v2"
	bigtableadmin "google.golang.org/api/bigtableadmin/v2"
	cloudfunctions "google.golang.org/api/cloudfunctions/v2"
	cloudkms "google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	composer "google.golang.org/api/composer/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	dataproc "google.golang.org/api/dataproc/v1"
	"google.golang.org/api/dns/v1"
	firestore "google.golang.org/api/firestore/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/networkconnectivity/v1"
	pubsub "google.golang.org/api/pubsub/v1"
	"google.golang.org/api/redis/v1"
	run "google.golang.org/api/run/v2"
	"google.golang.org/api/spanner/v1"
	"google.golang.org/api/sqladmin/v1"
	"google.golang.org/api/storage/v1"
	"google.golang.org/api/vpcaccess/v1"
)

// Services carries the resolved test project and the shared API clients every
// verifier probes through.
type Services struct {
	// Project is the resolved E2E test project id (exported as GOOGLE_PROJECT
	// for the IaC subprocesses); the fallback when outputs omit a project.
	Project             string
	Crm                 *cloudresourcemanager.Service
	Iam                 *iam.Service
	Compute             *compute.Service
	Storage             *storage.Service
	SqlAdmin            *sqladmin.Service
	Redis               *redis.Service
	Container           *container.Service
	Run                 *run.Service
	AlloyDB             *alloydb.Service
	DNS                 *dns.Service
	Spanner             *spanner.Service
	BigQuery            *bigquery.Service
	VpcAccess           *vpcaccess.Service
	Functions           *cloudfunctions.Service
	NetworkConnectivity *networkconnectivity.Service
	BigtableAdmin       *bigtableadmin.Service
	Firestore           *firestore.Service
	Dataproc            *dataproc.Service
	Composer            *composer.Service
	PubSub              *pubsub.Service
	CloudKms            *cloudkms.Service

	// RestClient is an ADC-authenticated HTTP client for GCP services whose
	// typed Go client is not yet in the pinned google.golang.org/api line
	// (e.g. the Memorystore for Valkey API). Verifiers use it for plain
	// REST GET probes only.
	RestClient *http.Client
}

// Verifier checks a single component's GCP resource for existence/absence.
type Verifier interface {
	// IDOutputKey is the stack-output key carrying the primary identifier —
	// used to confirm the deploy produced a verifiable handle.
	IDOutputKey() string
	// VerifyExists returns an error unless the resource exists.
	VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error
	// VerifyAbsent returns an error unless the resource is gone.
	VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error
}

// verifiers maps a component name to its verifier. New GCP components register
// here as they are forged.
var verifiers = map[string]Verifier{
	"gcpserviceaccount":                      &serviceAccountVerifier{},
	"gcpiamcustomrole":                       &iamCustomRoleVerifier{},
	"gcpprojectiammember":                    &projectIamMemberVerifier{},
	"gcpworkloadidentitypool":                &workloadIdentityPoolVerifier{},
	"gcpworkloadidentitypoolprovider":        &workloadIdentityPoolProviderVerifier{},
	"gcphealthcheck":                         &healthCheckVerifier{},
	"gcpbackendbucket":                       &backendBucketVerifier{},
	"gcpbackendservice":                      &backendServiceVerifier{},
	"gcpregionnetworkendpointgroup":          &regionNetworkEndpointGroupVerifier{},
	"gcpurlmap":                              &urlMapVerifier{},
	"gcpmanagedsslcertificate":               &managedSslCertificateVerifier{},
	"gcpsubnetwork":                          &subnetworkVerifier{},
	"gcpvpcnetwork":                          &vpcVerifier{},
	"gcpgcsbucket":                           &gcsBucketVerifier{},
	"gcptargethttpproxy":                     &targetHttpProxyVerifier{},
	"gcptargethttpsproxy":                    &targetHttpsProxyVerifier{},
	"gcpglobalforwardingrule":                &globalForwardingRuleVerifier{},
	"gcpglobaladdress":                       &globalAddressVerifier{},
	"gcpsslpolicy":                           &sslPolicyVerifier{},
	"gcpsslcertificate":                      &sslCertificateVerifier{},
	"gcpservicenetworkingconnection":         &serviceNetworkingConnectionVerifier{},
	"gcpaddress":                             &addressVerifier{},
	"gcpcloudsql":                            &cloudSqlInstanceVerifier{},
	"gcpcloudsqldatabase":                    &cloudSqlDatabaseVerifier{},
	"gcpcloudsqluser":                        &cloudSqlUserVerifier{},
	"gcpredisinstance":                       &redisInstanceVerifier{},
	"gcprouternat":                           &routerNatVerifier{},
	"gcpgkecluster":                          &gkeClusterVerifier{},
	"gcpgkenodepool":                         &gkeNodePoolVerifier{},
	"gcpcloudrun":                            &cloudRunVerifier{},
	"gcpalloydbcluster":                      &alloydbClusterVerifier{},
	"gcpalloydbinstance":                     &alloydbInstanceVerifier{},
	"gcpalloydbuser":                         &alloydbUserVerifier{},
	"gcpdnszone":                             &dnsZoneVerifier{},
	"gcpcloudrunjob":                         &cloudRunJobVerifier{},
	"gcpspannerinstance":                     &spannerInstanceVerifier{},
	"gcpspannerdatabase":                     &spannerDatabaseVerifier{},
	"gcpspannerbackupschedule":               &spannerBackupScheduleVerifier{},
	"gcpbigquerydataset":                     &bigQueryDatasetVerifier{},
	"gcpbigquerytable":                       &bigQueryTableVerifier{},
	"gcpserverlessvpcconnector":              &serverlessVpcConnectorVerifier{},
	"gcpcloudfunction":                       &cloudFunctionVerifier{},
	"gcpserviceconnectionpolicy":             &serviceConnectionPolicyVerifier{},
	"gcpmemorystoreinstance":                 &memorystoreInstanceVerifier{},
	"gcpfirestoredatabase":                   &firestoreDatabaseVerifier{},
	"gcpfirestorebackupschedule":             &firestoreBackupScheduleVerifier{},
	"gcpfirestoreindex":                      &firestoreIndexVerifier{},
	"gcpbigtableinstance":                    &bigtableInstanceVerifier{},
	"gcpbigtabletable":                       &bigtableTableVerifier{},
	"gcpfirewallrule":                        &firewallRuleVerifier{},
	"gcpdataproccluster":                     &dataprocClusterVerifier{},
	"gcpdataprocautoscalingpolicy":           &dataprocAutoscalingPolicyVerifier{},
	"gcpcloudcomposerenvironment":            &composerEnvironmentVerifier{},
	"gcpcloudcomposeruserworkloadssecret":    &composerUserWorkloadsSecretVerifier{},
	"gcpcloudcomposeruserworkloadsconfigmap": &composerUserWorkloadsConfigMapVerifier{},
	"gcppubsubschema":                        &pubSubSchemaVerifier{},
	"gcppubsubtopic":                         &pubSubTopicVerifier{},
	"gcppubsubsubscription":                  &pubSubSubscriptionVerifier{},
	"gcpkmskeyring":                          &kmsKeyRingVerifier{},
	"gcpkmskey":                              &kmsKeyVerifier{},
}

// GetVerifier returns the verifier for a component, or an error if none is registered.
func GetVerifier(component string) (Verifier, error) {
	v, ok := verifiers[component]
	if !ok {
		return nil, errors.Errorf("no GCP verifier registered for component %q", component)
	}
	return v, nil
}
