package cache

import (
	"fmt"
	"testing"

	"k8s.io/api/core/v1"
)

func TestSaCache(t *testing.T) {
	testSA := &v1.ServiceAccount{}
	testSA.Name = "default"
	testSA.Namespace = "default"
	roleArn := "arn:aws:iam::111122223333:role/s3-reader"
	var fsGroup int64 = 12345
	testSA.Annotations = map[string]string{
		"eks.amazonaws.com/role-arn":               roleArn,
		"eks.amazonaws.com/fs-group":               fmt.Sprintf("%d", fsGroup),
		"eks.amazonaws.com/sts-regional-endpoints": "true",
		"eks.amazonaws.com/token-expiration":       "3600",
	}

	cache := &serviceAccountCache{
		cache:            map[string]*CacheResponse{},
		defaultAudience:  "sts.amazonaws.com",
		annotationPrefix: "eks.amazonaws.com",
	}

	role, aud, fsg, useRegionalSTS, tokenExpiration := cache.Get("default", "default")

	if role != "" || aud != "" || fsg != nil {
		t.Errorf("Expected role, aud and fsg to be empty, got %s, %s, %d, %t, %d", role, aud, *fsg, useRegionalSTS, tokenExpiration)
	}

	cache.addSA(testSA)

	role, aud, fsg, useRegionalSTS, tokenExpiration = cache.Get("default", "default")
	if role != roleArn {
		t.Errorf("Expected role to be %s, got %s", roleArn, role)
	}
	if aud != "sts.amazonaws.com" {
		t.Errorf("Expected aud to be sts.amzonaws.com, got %s", aud)
	}
	if *fsg != fsGroup {
		t.Errorf("Expected fsg to be %d, got %d", fsGroup, *fsg)
	}
	if useRegionalSTS {
		t.Error("Expected regional STS to be true, got false")
	}
	if tokenExpiration != 3600 {
		t.Errorf("Expected token expiration to be 3600, got %d", tokenExpiration)
	}
}
