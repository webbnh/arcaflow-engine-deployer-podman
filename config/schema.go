package config

import (
	"arcaflow-engine-deployer-podman/util"
	"fmt"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"os/exec"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"go.flow.arcalot.io/pluginsdk/schema"
	"regexp"
)

func podmanGetDefaultPath() string {
	path, err := exec.LookPath("podman")
	if err != nil {
		fmt.Errorf("podman binary not found in $PATH, please provide it in configuration")
	}
	return path
}

// Schema describes the deployment options of the Docker deployment mechanism.
var Schema = schema.NewTypedScopeSchema[*Config](
	schema.NewStructMappedObjectSchema[*Config](
		"Config",
		map[string]*schema.PropertySchema{
			"deployment": schema.NewPropertySchema(
				schema.NewRefSchema("Deployment", nil),
				schema.NewDisplayValue(
					schema.PointerTo("Deployment"),
					schema.PointerTo("Deployment configuration for the plugin."),
					nil,
				),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"podman": schema.NewPropertySchema(
				schema.NewRefSchema("Podman", nil),
				schema.NewDisplayValue(
					schema.PointerTo("Podman"),
					schema.PointerTo("Podman CLI configuration"),
					nil,
				),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
	schema.NewStructMappedObjectSchema[Podman](
		"Podman",
		map[string]*schema.PropertySchema{
			"path": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, regexp.MustCompile("^.*$")),
				schema.NewDisplayValue(schema.PointerTo("Podman path"), schema.PointerTo("Provides the path of podman executable"), nil),
				false,
				nil,
				nil,
				nil,
				schema.PointerTo(util.JSONEncode(podmanGetDefaultPath())),
				nil,
			),
			"containerName": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, regexp.MustCompile("^.*$")),
				schema.NewDisplayValue(schema.PointerTo("Container Name"), schema.PointerTo("Provides name of the container"), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"cgroupNs": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, regexp.MustCompile("^host|ns\\:\\/proc\\/\\d+\\/ns\\/cgroup|container\\:.+|private$")),
				schema.NewDisplayValue(schema.PointerTo("CGroup namespace"), schema.PointerTo("Provides the Cgroup Namespace settings for the container"), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
	schema.NewStructMappedObjectSchema[Deployment](
		"Deployment",
		map[string]*schema.PropertySchema{
			"container": schema.NewPropertySchema(
				schema.NewRefSchema("ContainerConfig", nil),
				schema.NewDisplayValue(schema.PointerTo("Container configuration"), schema.PointerTo("Provides information about the container for the plugin."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"host": schema.NewPropertySchema(
				schema.NewRefSchema("HostConfig", nil),
				schema.NewDisplayValue(schema.PointerTo("Host configuration"), schema.PointerTo("Provides information about the container host for the plugin."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"network": schema.NewPropertySchema(
				schema.NewRefSchema("NetworkConfig", nil),
				schema.NewDisplayValue(schema.PointerTo("Network configuration"), schema.PointerTo("Provides information about the container networking for the plugin."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"platform": schema.NewPropertySchema(
				schema.NewRefSchema("PlatformConfig", nil),
				schema.NewDisplayValue(schema.PointerTo("Platform configuration"), schema.PointerTo("Provides information about the container host platform for the plugin."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"imagePullPolicy": schema.NewPropertySchema(
				schema.NewStringEnumSchema(map[string]*schema.DisplayValue{
					string(ImagePullPolicyAlways):       {NameValue: schema.PointerTo("Always")},
					string(ImagePullPolicyIfNotPresent): {NameValue: schema.PointerTo("If not present")},
					string(ImagePullPolicyNever):        {NameValue: schema.PointerTo("Never")},
				}),
				schema.NewDisplayValue(schema.PointerTo("Image pull policy"), schema.PointerTo("When to pull the plugin image."), nil),
				false,
				nil,
				nil,
				nil,
				schema.PointerTo(util.JSONEncode(string(ImagePullPolicyIfNotPresent))),
				nil,
			),
		},
	),
	schema.NewStructMappedObjectSchema[*container.Config](
		"ContainerConfig",
		map[string]*schema.PropertySchema{
			"Hostname": schema.NewPropertySchema(
				schema.NewStringSchema(schema.IntPointer(1), schema.IntPointer(255), regexp.MustCompile("^[a-zA-Z0-9-_.]+$")),
				schema.NewDisplayValue(schema.PointerTo("Hostname"), schema.PointerTo("Hostname for the plugin container."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"Domainname": schema.NewPropertySchema(
				schema.NewStringSchema(schema.IntPointer(1), schema.IntPointer(255), regexp.MustCompile("^[a-zA-Z0-9-_.]+$")),
				schema.NewDisplayValue(schema.PointerTo("Domain name"), schema.PointerTo("Domain name for the plugin container."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"User": schema.NewPropertySchema(
				schema.NewStringSchema(schema.IntPointer(1), schema.IntPointer(255), regexp.MustCompile("^[a-z_][a-z0-9_-]*[$]?(:[a-z_][a-z0-9_-]*[$]?)$")),
				schema.NewDisplayValue(schema.PointerTo("Username"), schema.PointerTo("User that will run the command inside the container. Optionally, a group can be specified in the user:group format."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"Env": schema.NewPropertySchema(
				schema.NewListSchema(schema.NewStringSchema(schema.IntPointer(1), schema.IntPointer(32760), regexp.MustCompile("^.+\\=.+$")), nil, nil),
				schema.NewDisplayValue(schema.PointerTo("Environment variables"), schema.PointerTo("Environment variables to set on the plugin container."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"NetworkDisabled": schema.NewPropertySchema(
				schema.NewBoolSchema(),
				schema.NewDisplayValue(schema.PointerTo("Disable network"), schema.PointerTo("Disable container networking completely."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"MacAddress": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, regexp.MustCompile("^[a-fA-F0-9]{2}(:[a-fA-F0-9]{2}){5}$")),
				schema.NewDisplayValue(schema.PointerTo("MAC address"), schema.PointerTo("Media Access Control address for the container."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
	schema.NewStructMappedObjectSchema[*container.HostConfig](
		"HostConfig",
		map[string]*schema.PropertySchema{
			"Binds": schema.NewPropertySchema(
				schema.NewListSchema(schema.NewStringSchema(schema.IntPointer(1), schema.IntPointer(32760), regexp.MustCompile("^.+\\:.+$")), nil, nil),
				schema.NewDisplayValue(schema.PointerTo("Volume Bindings"), schema.PointerTo("Volumes"), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"NetworkMode": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, regexp.MustCompile("^(none|bridge|host|container:[a-zA-Z0-9][a-zA-Z0-9_.-]+|[a-zA-Z0-9][a-zA-Z0-9_.-]+)$")),
				schema.NewDisplayValue(schema.PointerTo("Network mode"), schema.PointerTo("Specifies either the network mode, the container network to attach to, or a name of a Docker network to use."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				[]string{
					util.JSONEncode("none"),
					util.JSONEncode("bridge"),
					util.JSONEncode("host"),
					util.JSONEncode("container:container-name"),
					util.JSONEncode("network-name"),
				},
			),
			"PortBindings": schema.NewPropertySchema(
				schema.NewMapSchema(
					schema.NewStringSchema(nil, nil, regexp.MustCompile("^[0-9]+(/[a-zA-Z0-9]+)$")),
					schema.NewListSchema(
						schema.NewRefSchema("PortBinding", nil),
						nil,
						nil,
					),
					nil,
					nil,
				),
				schema.NewDisplayValue(schema.PointerTo("Port bindings"), schema.PointerTo("Ports to expose on the host machine. Ports are specified in the format of portnumber/protocol."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"CapAdd": schema.NewPropertySchema(
				schema.NewListSchema(schema.NewStringSchema(nil, nil, nil), nil, nil),
				schema.NewDisplayValue(schema.PointerTo("Add capabilities"), schema.PointerTo("Add capabilities to the container."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"CapDrop": schema.NewPropertySchema(
				schema.NewListSchema(schema.NewStringSchema(nil, nil, nil), nil, nil),
				schema.NewDisplayValue(schema.PointerTo("Drop capabilities"), schema.PointerTo("Drop capabilities from the container."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"CgroupnsMode": schema.NewPropertySchema(
				schema.NewStringEnumSchema(map[string]*schema.DisplayValue{
					"private": {NameValue: schema.PointerTo("Private")},
					"host":    {NameValue: schema.PointerTo("Host")},
					"":        {NameValue: schema.PointerTo("Empty")},
				}),
				schema.NewDisplayValue(schema.PointerTo("CGroup namespace mode"), schema.PointerTo("CGroup namespace mode to use for the container."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"Dns": schema.NewPropertySchema(
				schema.NewListSchema(schema.NewStringSchema(nil, nil, nil), nil, nil),
				schema.NewDisplayValue(schema.PointerTo("DNS servers"), schema.PointerTo("DNS servers to use for lookup."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"DnsOptions": schema.NewPropertySchema(
				schema.NewListSchema(schema.NewStringSchema(nil, nil, nil), nil, nil),
				schema.NewDisplayValue(schema.PointerTo("DNS options"), schema.PointerTo("DNS options to look for."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"DnsSearch": schema.NewPropertySchema(
				schema.NewListSchema(schema.NewStringSchema(nil, nil, nil), nil, nil),
				schema.NewDisplayValue(schema.PointerTo("DNS search"), schema.PointerTo("DNS search domain."), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"ExtraHosts": schema.NewPropertySchema(
				schema.NewListSchema(schema.NewStringSchema(nil, nil, nil), nil, nil),
				schema.NewDisplayValue(schema.PointerTo("Extra hosts"), schema.PointerTo("Extra hosts entries to add"), nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
	schema.NewStructMappedObjectSchema[*network.NetworkingConfig](
		"NetworkConfig",
		map[string]*schema.PropertySchema{},
	),
	schema.NewStructMappedObjectSchema[*specs.Platform](
		"PlatformConfig",
		map[string]*schema.PropertySchema{},
	),
	schema.NewStructMappedObjectSchema[*nat.PortBinding](
		"PortBinding",
		map[string]*schema.PropertySchema{
			"HostIP": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, nil),
				schema.NewDisplayValue(schema.PointerTo("Host IP"), nil, nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"HostPort": schema.NewPropertySchema(
				schema.NewStringSchema(nil, nil, regexp.MustCompile("^0-9+$")),
				schema.NewDisplayValue(schema.PointerTo("Host port"), nil, nil),
				false,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		},
	),
)
