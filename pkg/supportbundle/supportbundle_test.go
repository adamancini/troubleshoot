package supportbundle

import (
	"reflect"
	"testing"

	troubleshootv1beta2 "github.com/replicatedhq/troubleshoot/pkg/apis/troubleshoot/v1beta2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_LoadAndConcatSpec(t *testing.T) {

	bundle1doc, err := LoadSupportBundleSpec("test/supportbundle1.yaml")
	if err != nil {
		t.Error("couldn't load bundle1 from file")
	}

	bundle2doc, err := LoadSupportBundleSpec("test/supportbundle2.yaml")
	if err != nil {
		t.Error("couldn't load bundle2 from file")
	}

	bundle1, err := ParseSupportBundleFromDoc(bundle1doc)
	if err != nil {
		t.Error("couldn't parse bundle 1")
	}

	bundle2, err := ParseSupportBundleFromDoc(bundle2doc)
	if err != nil {
		t.Error("couldn't parse bundle 2")
	}

	fulldoc, err := LoadSupportBundleSpec("test/completebundle.yaml")
	if err != nil {
		t.Error("couldn't load full bundle from file")
	}

	fullbundle, err := ParseSupportBundleFromDoc(fulldoc)
	if err != nil {
		t.Error("couldn't parse full bundle")
	}

	bundle3 := ConcatSpec(bundle1, bundle2)

	if reflect.DeepEqual(fullbundle, bundle3) == false {
		t.Error("Full bundle and concatenated bundle are not the same.")
	}

}

func Test_ParseSupportBundle(t *testing.T) {

	upstreamSupportBundleDoc := []byte(`
apiVersion: troubleshoot/v1beta2
kind: SupportBundle
spec:
  uri: https://go-test/upstream.yaml
	collectors:
	  cluster-info: {}
		cluster-resources: {}
	analyzers: {}`)

	tests := []struct {
		name    string
		given   []byte
		expect  *troubleshootv1beta2.SupportBundle
		wantErr bool
	}{
		{
			name: "given a spec with a uri should return the upstream spec",
			given: []byte(`
apiVersion: troubleshoot/v1beta2
kind: SupportBundle
spec:
  uri: https://go-test/upstream.yaml
	collectors:
	  cluster-info: {}
	analyzers: {}`),
			expect: &troubleshootv1beta2.SupportBundle{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "troubleshoot/v1beta2",
					Kind:       "SupportBundle",
				},
				Spec: troubleshootv1beta2.SupportBundleSpec{
					Uri: "https://go-test/upstream.yaml",
					Collectors: []*troubleshootv1beta2.Collect{
						{
							ClusterInfo: &troubleshootv1beta2.ClusterInfo{},
						},
						{
							ClusterResources: &troubleshootv1beta2.ClusterResources{},
						},
					},
					Analyzers: []*troubleshootv1beta2.Analyze{},
				},
			},
			wantErr: false,
		},
		// 		{
		// 			name: "if the uri is unreachable, return the given spec",
		// 			upstream: []byte(`
		// apiVersion: troubleshoot/v1beta2
		// kind: SupportBundle
		// spec:
		//   uri: https://go-test/upstream.yaml
		// 	collectors:
		// 	  cluster-info: {}
		// 		cluster-resources: {}
		// 	analyzers: {}`),
		// 			given:   &troubleshootv1beta2.SupportBundle{},
		// 			expect:  given,
		// 			wantErr: false,
		// 		},
		// 		{
		// 			name: "if the upstream spec is unparseable, return the given spec",
		// 			upstream: []byte(`
		// apiVersion: troubleshoot/v1beta2
		// kind: SupportBundle
		// spec:
		//   uri: https://go-test/upstream.yaml
		// 	collectors:
		// 	  cluster-info: {}
		// 		cluster-resources: {}
		// 	analyzers: {}`),
		// 			given:   &troubleshootv1beta2.SupportBundle{},
		// 			expect:  given,
		// 			wantErr: false,
		// 		},
		// 		{
		// 			name: "if the given spec is only a uri, and if it is unreachable, return an error",
		// 			upstream: []byte(`
		// apiVersion: troubleshoot/v1beta2
		// kind: SupportBundle
		// spec:
		//   uri: https://go-test/upstream.yaml
		// 	collectors:
		// 	  cluster-info: {}
		// 		cluster-resources: {}
		// 	analyzers: {}`),
		// 			given:   &troubleshootv1beta2.SupportBundle{},
		// 			expect:  nil,
		// 			wantErr: true,
		// 		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			givenBundle, err := ParseSupportBundle(test.given, true)
			req.NoError(err)

			assert.Equal(t, test.expect, givenBundle)
		})
	}
}
