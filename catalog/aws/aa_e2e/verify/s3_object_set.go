package verify

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	pkgerrors "github.com/pkg/errors"
)

// s3ObjectSetVerifier verifies an AwsS3ObjectSet by HeadObject on each key
// named in the object_etags output map. The bucket_id output carries the
// target bucket. Destroy removes the objects (not the bucket), so absence is
// NoSuchKey on every staged key.
type s3ObjectSetVerifier struct{}

func (*s3ObjectSetVerifier) IDOutputKey() string { return "bucket_id" }

func (v *s3ObjectSetVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	return pkgerrors.Errorf("awss3objectset verify-exists requires full outputs (object_etags); use OutputsVerifier path")
}

func (v *s3ObjectSetVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	return pkgerrors.Errorf("awss3objectset verify-absent requires full outputs (object_etags); use OutputsVerifier path")
}

func (*s3ObjectSetVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	bucket, keys, err := s3ObjectSetTargets(outputs)
	if err != nil {
		return pkgerrors.Wrap(err, "awss3objectset verify-exists")
	}
	for _, key := range keys {
		exists, err := headObject(ctx, cfg, bucket, key, region)
		if err != nil {
			return pkgerrors.Wrapf(err, "awss3objectset verify-exists failed for s3://%s/%s", bucket, key)
		}
		if !exists {
			return pkgerrors.Errorf("awss3objectset object s3://%s/%s not found after deploy", bucket, key)
		}
	}
	return nil
}

func (*s3ObjectSetVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	bucket, keys, err := s3ObjectSetTargets(outputs)
	if err != nil {
		return pkgerrors.Wrap(err, "awss3objectset verify-absent")
	}
	for _, key := range keys {
		exists, err := headObject(ctx, cfg, bucket, key, region)
		if err != nil {
			return pkgerrors.Wrapf(err, "awss3objectset verify-absent failed for s3://%s/%s", bucket, key)
		}
		if exists {
			return pkgerrors.Errorf("awss3objectset object s3://%s/%s still exists after destroy", bucket, key)
		}
	}
	return nil
}

func s3ObjectSetTargets(outputs map[string]interface{}) (bucket string, keys []string, err error) {
	bucket = stringOutputMap(outputs, "bucket_id")
	if bucket == "" {
		return "", nil, pkgerrors.New("no bucket_id in outputs -- cannot verify")
	}
	etags, ok := outputs["object_etags"]
	if !ok || etags == nil {
		return "", nil, pkgerrors.New("no object_etags in outputs -- cannot verify")
	}
	switch m := etags.(type) {
	case map[string]interface{}:
		for k := range m {
			keys = append(keys, k)
		}
	case map[string]string:
		for k := range m {
			keys = append(keys, k)
		}
	default:
		return "", nil, pkgerrors.Errorf("object_etags has unexpected type %T", etags)
	}
	if len(keys) == 0 {
		return "", nil, pkgerrors.New("object_etags is empty -- nothing to verify")
	}
	return bucket, keys, nil
}

func stringOutputMap(outputs map[string]interface{}, key string) string {
	if outputs == nil {
		return ""
	}
	if v, ok := outputs[key]; ok {
		if s, isStr := v.(string); isStr {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func headObject(ctx context.Context, cfg aws.Config, bucket, key, region string) (bool, error) {
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		return true, nil
	}
	var notFound *s3types.NotFound
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}
