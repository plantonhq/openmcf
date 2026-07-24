package verify

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscredentials "github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

// SeaweedFsVerifier checks a SeaweedFS object store to the point clients
// can rely on it: master and filer StatefulSets ready (the dedicated S3
// Deployment too when that shape is deployed), the `-s3` Service present,
// and — on lanes with S3 auth — a LIVE S3 put/get round-trip through the
// gateway with the chart-materialized credentials (an object store that
// cannot store an object is not an object store). Declared buckets are
// asserted through the S3 ListBuckets API, which proves the chart's
// bucket-creation hook actually ran.
//
// The behavioral-durability scenario (recognized by name) additionally
// DELETES the volume-server pod after the put and re-reads the object once
// it returns — object bytes surviving pod loss through the PVC is the
// proof.
type SeaweedFsVerifier struct {
	Namespace string
	Name      string
	// DedicatedS3 marks the shape where the gateway runs as its own
	// Deployment (`<name>-s3`) instead of embedded on the filer.
	DedicatedS3 bool
	// CredentialsSecretName is the chart-generated `<name>-s3-secret` (or
	// the referenced existing config secret). Empty = auth disabled; the
	// S3 proof runs unauthenticated.
	CredentialsSecretName string
	// Buckets declared in the spec — each must exist after install (the
	// chart's post-install hook creates them).
	Buckets []string
	// AdminEnabled asserts the admin console Deployment and Service.
	AdminEnabled bool
	Durability   bool
}

func (v *SeaweedFsVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] seaweedfs store %q in namespace %q\n", v.Name, v.Namespace)

	if err := v.waitForStatefulSetReady(ctx, kubeconfig, v.Name+"-master", 10*time.Minute); err != nil {
		return err
	}
	if err := v.waitForStatefulSetReady(ctx, kubeconfig, v.Name+"-filer", 10*time.Minute); err != nil {
		return err
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name+"-s3", v.Namespace); err != nil {
		return errors.Wrap(err, "s3 service not found")
	}
	if v.DedicatedS3 {
		if err := KubectlResourceExists(ctx, kubeconfig, "deployment", v.Name+"-s3", v.Namespace); err != nil {
			return errors.Wrap(err, "dedicated s3 gateway deployment not found")
		}
	}
	if v.AdminEnabled {
		if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name+"-admin", v.Namespace); err != nil {
			return errors.Wrap(err, "admin console service not found")
		}
	}
	return v.proveS3RoundTrip(ctx, kubeconfig)
}

func (v *SeaweedFsVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "statefulset", v.Name+"-master", v.Namespace)
}

func (v *SeaweedFsVerifier) waitForStatefulSetReady(ctx context.Context, kubeconfig, name string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var lastReady string
	for time.Now().Before(deadline) {
		ready, _ := kubectlGetJSONPath(ctx, kubeconfig, "statefulset", name, v.Namespace, "{.status.readyReplicas}")
		replicas, _ := kubectlGetJSONPath(ctx, kubeconfig, "statefulset", name, v.Namespace, "{.spec.replicas}")
		lastReady = ready
		if ready != "" && ready == replicas {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("statefulset %s never became ready (last readyReplicas %q)", name, lastReady)
}

// credentials reads the admin access/secret key pair from the
// chart-materialized credentials Secret (its stringData keys).
func (v *SeaweedFsVerifier) credentials(ctx context.Context, kubeconfig string) (string, string, error) {
	accessB64, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", v.CredentialsSecretName, v.Namespace, "{.data.admin_access_key_id}")
	if err != nil {
		return "", "", errors.Wrapf(err, "reading secret %q admin_access_key_id", v.CredentialsSecretName)
	}
	secretB64, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", v.CredentialsSecretName, v.Namespace, "{.data.admin_secret_access_key}")
	if err != nil {
		return "", "", errors.Wrapf(err, "reading secret %q admin_secret_access_key", v.CredentialsSecretName)
	}
	access, err := base64.StdEncoding.DecodeString(accessB64)
	if err != nil {
		return "", "", err
	}
	secret, err := base64.StdEncoding.DecodeString(secretB64)
	if err != nil {
		return "", "", err
	}
	return strings.TrimSpace(string(access)), strings.TrimSpace(string(secret)), nil
}

// proveS3RoundTrip drives the S3 API over a port-forward to the `-s3`
// Service: assert every declared bucket exists, put a run-unique object,
// read it back and compare bytes. The durability variant deletes the
// volume-server pod between put and get.
func (v *SeaweedFsVerifier) proveS3RoundTrip(ctx context.Context, kubeconfig string) error {
	const localPort = "18333"

	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+v.Name+"-s3", localPort+":8333", "-n", v.Namespace)
	var pfOut strings.Builder
	pf.Stdout = &pfOut
	pf.Stderr = &pfOut
	if err := pf.Start(); err != nil {
		cancel()
		return errors.Wrap(err, "starting port-forward to the s3 service")
	}
	// ONE deferred func, cancel FIRST — Wait blocks forever on a
	// port-forward that is never told to exit.
	defer func() {
		cancel()
		_ = pf.Wait()
	}()

	accessKey, secretKey := "any", "any"
	if v.CredentialsSecretName != "" {
		var err error
		accessKey, secretKey, err = v.credentials(ctx, kubeconfig)
		if err != nil {
			return err
		}
	}

	// SeaweedFS serves path-style S3; region is accepted but ignored.
	client := awss3.New(awss3.Options{
		Region:       "us-east-1",
		BaseEndpoint: aws.String("http://127.0.0.1:" + localPort),
		UsePathStyle: true,
		Credentials:  awscredentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	})

	// The gateway answers after the filer finishes wiring; retry the
	// first call across the warm-up window.
	if err := v.waitForBuckets(ctx, client, 4*time.Minute); err != nil {
		return err
	}

	bucket := "e2e-proof"
	if len(v.Buckets) > 0 {
		bucket = v.Buckets[0]
	} else {
		if _, err := client.CreateBucket(ctx, &awss3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
			return errors.Wrap(err, "creating the proof bucket")
		}
	}

	key := fmt.Sprintf("e2e-marker-%d.txt", time.Now().Unix())
	payload := []byte("seaweedfs-e2e-proof " + key)
	if _, err := client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(payload),
	}); err != nil {
		return errors.Wrap(err, "the S3 put never succeeded")
	}
	fmt.Printf("  [verify] S3: object %q written to bucket %q\n", key, bucket)

	if v.Durability {
		pod := v.Name + "-volume-0"
		fmt.Printf("  [verify] DURABILITY: deleting volume-server pod %q\n", pod)
		if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"delete", "pod", pod, "-n", v.Namespace, "--wait=false").CombinedOutput(); err != nil {
			return errors.Wrapf(err, "deleting the volume-server pod: %s", string(out))
		}
		if err := v.waitForStatefulSetReady(ctx, kubeconfig, v.Name+"-volume", 10*time.Minute); err != nil {
			return errors.Wrap(err, "the volume server never returned after deletion")
		}
	}

	got, err := v.getObjectWithRetry(ctx, client, bucket, key, 4*time.Minute)
	if err != nil {
		return errors.Wrap(err, "the S3 get never succeeded")
	}
	if !bytes.Equal(got, payload) {
		return errors.Errorf("the object came back different: got %q", string(got))
	}
	if v.Durability {
		fmt.Printf("  [verify] DURABILITY: object read back AFTER volume-server loss — bytes survived on the PVC\n")
	} else {
		fmt.Printf("  [verify] S3: object read back byte-identical\n")
	}
	return nil
}

// waitForBuckets retries ListBuckets until the gateway answers, then
// asserts every declared bucket is present (the chart's post-install hook
// creates them).
func (v *SeaweedFsVerifier) waitForBuckets(ctx context.Context, client *awss3.Client, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		out, err := client.ListBuckets(ctx, &awss3.ListBucketsInput{})
		if err == nil {
			present := map[string]bool{}
			for _, b := range out.Buckets {
				present[aws.ToString(b.Name)] = true
			}
			for _, want := range v.Buckets {
				if !present[want] {
					return errors.Errorf("declared bucket %q was not created by the install hook (found: %v)", want, bucketNames(out))
				}
			}
			if len(v.Buckets) > 0 {
				fmt.Printf("  [verify] S3: all %d declared buckets present\n", len(v.Buckets))
			}
			return nil
		}
		lastErr = err
		time.Sleep(10 * time.Second)
	}
	return errors.Wrap(lastErr, "the S3 gateway never answered ListBuckets")
}

func (v *SeaweedFsVerifier) getObjectWithRetry(ctx context.Context, client *awss3.Client, bucket, key string, budget time.Duration) ([]byte, error) {
	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		out, err := client.GetObject(ctx, &awss3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err == nil {
			defer out.Body.Close()
			return io.ReadAll(out.Body)
		}
		lastErr = err
		time.Sleep(10 * time.Second)
	}
	return nil, lastErr
}

func bucketNames(out *awss3.ListBucketsOutput) []string {
	names := make([]string, 0, len(out.Buckets))
	for _, b := range out.Buckets {
		names = append(names, aws.ToString(b.Name))
	}
	return names
}
