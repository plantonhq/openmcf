package module

// Keys exported by the gcp_cloud_run Pulumi module.
const (
	OpUrl         = "url"          // Canonical serving URL of the service
	OpServiceName = "service_name" // Name of the Cloud Run service in GCP
	OpRevision    = "revision"     // Latest ready revision name
	OpLocation    = "location"     // Region the service is deployed in
	OpUid         = "uid"          // Server-assigned unique identifier
	OpUrls        = "urls"         // Every URL serving this service
)
