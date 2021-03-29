package cmd

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	cobradoc "github.com/spf13/cobra/doc"
	"sigs.k8s.io/yaml"

	"github.com/linkerd/linkerd2/pkg/k8s"
)

type references struct {
	CLIReference         []cmdDoc
	AnnotationsReference []annotationDoc
}

type cmdOption struct {
	Name         string
	Shorthand    string
	DefaultValue string
	Usage        string
}

type cmdDoc struct {
	Name             string
	Synopsis         string
	Description      string
	Options          []cmdOption
	InheritedOptions []cmdOption
	Example          string
	SeeAlso          []string
}

type annotationDoc struct {
	Name        string
	Description string
}

func newCmdDoc() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "doc",
		Hidden: true,
		Short:  "Generate YAML documentation for the Linkerd CLI & Proxy annotations",
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdList, err := generateCLIDocs(RootCmd)
			if err != nil {
				return err
			}

			annotations := generateAnnotationsDocs()

			ref := references{
				CLIReference:         cmdList,
				AnnotationsReference: annotations,
			}
			out, err := yaml.Marshal(ref)
			if err != nil {
				return err
			}

			warn := "# Automatically generated by the linkerd doc command, do not manually edit"
			fmt.Printf("%s\n\n%s\n", warn, out)

			return nil
		},
	}

	return cmd
}

// generateCLIDocs takes a command and recursively walks the tree of commands,
// adding each as an item to cmdList.
func generateCLIDocs(cmd *cobra.Command) ([]cmdDoc, error) {
	var cmdList []cmdDoc

	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		subList, err := generateCLIDocs(c)
		if err != nil {
			return nil, err
		}

		cmdList = append(cmdList, subList...)
	}

	var buf bytes.Buffer

	if err := cobradoc.GenYaml(cmd, io.Writer(&buf)); err != nil {
		return nil, err
	}

	var doc cmdDoc
	if err := yaml.Unmarshal(buf.Bytes(), &doc); err != nil {
		return nil, err
	}

	// Cobra names start with linkerd, strip that off for the docs.
	doc.Name = strings.TrimPrefix(doc.Name, "linkerd ")

	// Don't include the root command.
	if doc.Name != "linkerd" {
		cmdList = append(cmdList, doc)
	}

	return cmdList, nil
}

// generateAnnotationsDocs make list of annotations and its docs
func generateAnnotationsDocs() []annotationDoc {
	return []annotationDoc{
		{
			Name:        k8s.ProxyInjectAnnotation,
			Description: "Controls whether or not a pod should be injected; accepted values are `enabled`, `disabled` and `ingress`",
		},
		{
			Name:        k8s.ProxyImageAnnotation,
			Description: "Linkerd proxy container image name",
		},
		{
			Name:        k8s.ProxyImagePullPolicyAnnotation,
			Description: "Docker image pull policy",
		},
		{
			Name:        k8s.ProxyInitImageAnnotation,
			Description: "Linkerd init container image name",
		},
		{
			Name:        k8s.ProxyInitImageVersionAnnotation,
			Description: "Linkerd init container image version",
		},
		{
			Name:        k8s.DebugImageAnnotation,
			Description: "Linkerd debug container image name",
		},
		{
			Name:        k8s.DebugImageVersionAnnotation,
			Description: "Linkerd debug container image version",
		},
		{
			Name:        k8s.DebugImagePullPolicyAnnotation,
			Description: "Docker image pull policy for debug image",
		},
		{
			Name:        k8s.ProxyControlPortAnnotation,
			Description: "Proxy port to use for control",
		},
		{
			Name:        k8s.ProxyIgnoreInboundPortsAnnotation,
			Description: "Ports that should skip the proxy and send directly to the application. Comma-separated list of values, where each value can be a port number or a range `a-b`.",
		},
		{
			Name:        k8s.ProxyOpaquePortsAnnotation,
			Description: "Ports that skip the proxy's protocol detection mechanism and are proxied opaquely. Comma-separated list of values, where each value can be a port number or a range `a-b`.",
		},
		{
			Name:        k8s.ProxyIgnoreOutboundPortsAnnotation,
			Description: "Outbound ports that should skip the proxy. Comma-separated list of values, where each value can be a port number or a range `a-b`.",
		},
		{
			Name:        k8s.ProxyInboundPortAnnotation,
			Description: "Proxy port to use for inbound traffic",
		},
		{
			Name:        k8s.ProxyAdminPortAnnotation,
			Description: "Proxy port to serve metrics on",
		},
		{
			Name:        k8s.ProxyOutboundPortAnnotation,
			Description: "Proxy port to use for outbound traffic",
		},
		{
			Name:        k8s.ProxyCPURequestAnnotation,
			Description: "Amount of CPU units that the proxy sidecar requests",
		},
		{
			Name:        k8s.ProxyMemoryRequestAnnotation,
			Description: "Amount of Memory that the proxy sidecar requests",
		},
		{
			Name:        k8s.ProxyCPULimitAnnotation,
			Description: "Maximum amount of CPU units that the proxy sidecar can use",
		},
		{
			Name:        k8s.ProxyMemoryLimitAnnotation,
			Description: "Maximum amount of Memory that the proxy sidecar can use",
		},
		{
			Name:        k8s.ProxyUIDAnnotation,
			Description: "Run the proxy under this user ID",
		},
		{
			Name:        k8s.ProxyLogLevelAnnotation,
			Description: "Log level for the proxy",
		},
		{
			Name:        k8s.ProxyLogFormatAnnotation,
			Description: "Log format (plain or json) for the proxy",
		},
		{
			Name:        k8s.ProxyEnableExternalProfilesAnnotation,
			Description: "Enable service profiles for non-Kubernetes services",
		},
		{
			Name:        k8s.ProxyVersionOverrideAnnotation,
			Description: "Tag to be used for the Linkerd proxy images",
		},
		{
			Name:        k8s.ProxyDisableIdentityAnnotation,
			Description: "Disables resources from participating in TLS identity",
		},
		{
			Name:        k8s.ProxyEnableDebugAnnotation,
			Description: "Inject a debug sidecar for data plane debugging",
		},
		{
			Name:        k8s.ProxyOutboundConnectTimeout,
			Description: "Used to configure the outbound TCP connection timeout in the proxy",
		},
		{
			Name:        k8s.ProxyWaitBeforeExitSecondsAnnotation,
			Description: "The proxy sidecar will stay alive for at least the given period before receiving SIGTERM signal from Kubernetes but no longer than pod's `terminationGracePeriodSeconds`. If not provided, it will be defaulted to `0`",
		},
		{
			Name:        k8s.AwaitProxy,
			Description: "The application container will not start until the proxy is ready",
		},
	}
}
