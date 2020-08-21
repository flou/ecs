package aws

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/fatih/color"
)

// DetailedInstanceOutput prints a container instance's attributes and capabilities
func DetailedInstanceOutput(containerInstance *ecs.ContainerInstance) {
	var line string
	instanceAttributes := make([]string, 0)
	capabilities := make([]string, 0)
	for _, attr := range containerInstance.Attributes {
		if strings.Contains(*attr.Name, "ecs.capability.") {
			capability := strings.SplitAfter(*attr.Name, "ecs.capability.")[1]
			if strings.HasPrefix(capability, "docker-remote-api.") {
				continue
			}
			if attr.Value == nil {
				line = fmt.Sprintf(" - %s", capability)
			} else {
				line = fmt.Sprintf(" - %-22s %s", capability, color.YellowString(*attr.Value))
			}
			capabilities = append(capabilities, line)
		} else {
			if attr.Value == nil {
				line = fmt.Sprintf(" - %s", *attr.Name)
			} else {
				line = fmt.Sprintf(" - %-22s %s", *attr.Name, color.YellowString(*attr.Value))
			}
			instanceAttributes = append(instanceAttributes, line)
		}
	}
	fmt.Println("Attributes:")
	sort.Strings(instanceAttributes)
	for _, attr := range instanceAttributes {
		fmt.Println(attr)
	}
	fmt.Println("Capabilities:")
	sort.Strings(capabilities)
	for _, attr := range capabilities {
		fmt.Println(attr)
	}
	fmt.Println()
}

// FindResource finds a specific resource in a list of ecs.Resource
func FindResource(resources []ecs.Resource, name string) ecs.Resource {
	var resource ecs.Resource
	for _, res := range resources {
		if *res.Name == name {
			resource = res
			break
		}
	}
	return resource
}

// FindAttribute finds a specific attribute in a list of ecs.Attribute
func FindAttribute(attributes []ecs.Attribute, name string) ecs.Attribute {
	var attribute ecs.Attribute
	for _, attr := range attributes {
		if *attr.Name == name {
			attribute = attr
			break
		}
	}
	return attribute
}

// FindTag finds a specific tag in a list of EC2 tags
func FindTag(tags []ec2.Tag, name string) ec2.Tag {
	var tag ec2.Tag
	for _, t := range tags {
		if *t.Key == name {
			tag = t
			break
		}
	}
	return tag
}
