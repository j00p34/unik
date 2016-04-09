package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/emc-advanced-dev/unik/pkg/config"
	"github.com/emc-advanced-dev/unik/pkg/state"
	"github.com/layer-x/layerx-commons/lxlog"
)

const awsStateFile = "/var/unik/aws_state.json"

type AwsProvider struct {
	config  config.Aws  `json:"Config"`
	State   state.State `json:"State"`
}

func NewAwsProvier(config config.Aws) *AwsProvider {
	return &AwsProvider{
		config: config,
		State: state.NewMemoryState(awsStateFile),
	}
}

func (p *AwsProvider) newEC2(logger lxlog.Logger) *ec2.EC2 {
	sess := session.New(&aws.Config{
		Region:      aws.String(p.config.Region),
		Credentials: credentials.NewStaticCredentials(p.config.AwsAccessKeyID, p.config.AwsSecretAcessKey, ""),
	})
	sess.Handlers.Send.PushFront(func(r *request.Request) {
		logger.WithFields(
			lxlog.Fields{"request": r}).Debugf("request sent to aws")
	})
	return ec2.New(sess)
}

func (p *AwsProvider) newMetadata() *ec2metadata.EC2Metadata {
	return ec2metadata.New(session.New(&aws.Config{
		Region:      aws.String(p.config.Region),
		Credentials: credentials.NewStaticCredentials(p.config.AwsAccessKeyID, p.config.AwsSecretAcessKey, ""),
	}))
}