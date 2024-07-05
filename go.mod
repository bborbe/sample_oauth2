module github.com/bborbe/sample_oauth2

go 1.22.4

exclude (
	k8s.io/api v0.29.0
	k8s.io/api v0.29.1
	k8s.io/api v0.29.2
	k8s.io/api v0.29.3
	k8s.io/api v0.29.4
	k8s.io/api v0.29.5
	k8s.io/api v0.29.6
	k8s.io/api v0.30.0
	k8s.io/api v0.30.1
	k8s.io/api v0.30.2
	k8s.io/apiextensions-apiserver v0.29.0
	k8s.io/apiextensions-apiserver v0.29.1
	k8s.io/apiextensions-apiserver v0.29.2
	k8s.io/apiextensions-apiserver v0.29.3
	k8s.io/apiextensions-apiserver v0.29.4
	k8s.io/apiextensions-apiserver v0.29.5
	k8s.io/apiextensions-apiserver v0.29.6
	k8s.io/apiextensions-apiserver v0.30.0
	k8s.io/apiextensions-apiserver v0.30.1
	k8s.io/apiextensions-apiserver v0.30.2
	k8s.io/apimachinery v0.29.0
	k8s.io/apimachinery v0.29.1
	k8s.io/apimachinery v0.29.2
	k8s.io/apimachinery v0.29.3
	k8s.io/apimachinery v0.29.4
	k8s.io/apimachinery v0.29.5
	k8s.io/apimachinery v0.29.6
	k8s.io/apimachinery v0.30.0
	k8s.io/apimachinery v0.30.1
	k8s.io/apimachinery v0.30.2
	k8s.io/client-go v0.29.0
	k8s.io/client-go v0.29.1
	k8s.io/client-go v0.29.2
	k8s.io/client-go v0.29.3
	k8s.io/client-go v0.29.4
	k8s.io/client-go v0.29.5
	k8s.io/client-go v0.29.6
	k8s.io/client-go v0.30.0
	k8s.io/client-go v0.30.1
	k8s.io/client-go v0.30.2
	k8s.io/code-generator v0.29.0
	k8s.io/code-generator v0.29.1
	k8s.io/code-generator v0.29.2
	k8s.io/code-generator v0.29.3
	k8s.io/code-generator v0.29.4
	k8s.io/code-generator v0.29.5
	k8s.io/code-generator v0.29.6
	k8s.io/code-generator v0.30.0
	k8s.io/code-generator v0.30.1
	k8s.io/code-generator v0.30.2
)

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.16

replace github.com/antlr/antlr4/runtime/Go/antlr/v4 => github.com/antlr4-go/antlr/v4 v4.13.0

require (
	github.com/bborbe/errors v1.2.0
	github.com/bborbe/http v1.2.0
	github.com/bborbe/log v1.0.0
	github.com/bborbe/sentry v1.7.0
	github.com/bborbe/service v1.3.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/glog v1.2.1
	github.com/google/addlicense v1.1.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/incu6us/goimports-reviser v0.1.6
	github.com/kisielk/errcheck v1.7.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.8.1
	github.com/onsi/ginkgo/v2 v2.19.0
	github.com/onsi/gomega v1.33.1
	github.com/prometheus/client_golang v1.19.1
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/oauth2 v0.19.0
	golang.org/x/vuln v1.1.2
	k8s.io/code-generator v0.28.11
)

require (
	cloud.google.com/go/compute v1.20.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/bborbe/argument/v2 v2.0.4 // indirect
	github.com/bborbe/collection v1.4.0 // indirect
	github.com/bborbe/math v1.0.0 // indirect
	github.com/bborbe/run v1.5.3 // indirect
	github.com/bborbe/time v1.2.0 // indirect
	github.com/bborbe/validation v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/emicklei/go-restful/v3 v3.12.1 // indirect
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20240521024322-9665fa269a30 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.54.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/telemetry v0.0.0-20240621194115-a740542b267c // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/gengo v0.0.0-20240404160639-a0386bf69313 // indirect
	k8s.io/gengo/v2 v2.0.0-20240404160639-a0386bf69313 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240620174524-b456828f718b // indirect
	k8s.io/utils v0.0.0-20240502163921-fe8a2dddb1d0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
