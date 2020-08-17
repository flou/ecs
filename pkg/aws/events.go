package aws

import "github.com/aws/aws-sdk-go-v2/service/ecs"

type EventsByCreatedAt []ecs.ServiceEvent

func (c EventsByCreatedAt) Len() int           { return len(c) }
func (c EventsByCreatedAt) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c EventsByCreatedAt) Less(i, j int) bool { return c[i].CreatedAt.Before(*c[j].CreatedAt) }
