package cf

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"

	"fmt"
)

//Executor is an interface to execute and create stacks
type Executor interface {
	CreateStack(templateBody string, stackName string, parameterMap *map[string]string, tags *map[string]string) error
	UpdateStack(templateBody string, stackName string, parameterMap *map[string]string, tags *map[string]string) error

	PauseUntilCreateFinished(stackName string) error
	PauseUntilUpdateFinished(stackName string) error
}

//IaaSExecutor struct used to create cloudformation stacks
type IaaSExecutor struct {
	Client cloudformationiface.CloudFormationAPI
}

//UpdateStack is a method to update Cloudformation stack
func (executor IaaSExecutor) UpdateStack(templateBody string, stackName string,
	parameterMap *map[string]string, tags *map[string]string) error {

	//generate cloudformation CreateStackInput to be used to create stack
	input := &cloudformation.UpdateStackInput{}

	input.SetTemplateBody(*aws.String(templateBody))
	input.SetStackName(*aws.String(stackName))
	input.SetParameters(CreateCloudformationParameters(parameterMap))
	input.SetCapabilities(createCapability())
	if tags != nil {
		input.SetTags(createTags(tags))
	}
	//todo-refactor to return output
	_, err := executor.Client.UpdateStack(input)
	//if there's an error return it
	if err != nil {
		fmt.Println("Got error creating stack: ", err.Error())
		return errors.New("Error creating stack")

	}
	return nil
}

//CreateStack is a general method to create aws cloudformation stacks
func (executor IaaSExecutor) CreateStack(templateBody string, stackName string,
	parameterMap *map[string]string, tags *map[string]string) error {
	//generate cloudformation CreateStackInput to be used to create stack
	input := &cloudformation.CreateChangeSetInput{}

	input.SetTemplateBody(*aws.String(templateBody))
	input.SetStackName(*aws.String(stackName))
	input.SetParameters(CreateCloudformationParameters(parameterMap))
	input.SetCapabilities(createCapability())
	input.SetChangeSetType("CREATE")

	if tags != nil {
		input.SetTags(createTags(tags))
	}
	//todo-refactor to return output
	_, err := executor.Client.CreateChangeSet(input)
	//if there's an error return it
	if err != nil {
		fmt.Println("Got error creating stack: ", err.Error())
		return errors.New("Error creating stack")

	}
	return nil

}

//PauseUntilCreateFinished is a method to wait on the status of a cloudformation stack until it finishes
func (executor IaaSExecutor) PauseUntilCreateFinished(stackName string) error {
	fmt.Println("Waiting for stack to be created")

	// Wait until stack is created
	desInput := &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}
	err := executor.Client.WaitUntilStackCreateComplete(desInput)
	if err != nil {
		fmt.Println("Got error waiting for stack to be created")
		fmt.Println(err)
	}
	return err
}

//PauseUntilUpdateFinished is a method to wait on the status of a cloudformation stack until it finishes
func (executor IaaSExecutor) PauseUntilUpdateFinished(stackName string) error {
	fmt.Println("Waiting for stack to be updated")

	// Wait until stack is created
	desInput := &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}
	err := executor.Client.WaitUntilStackUpdateComplete(desInput)
	if err != nil {
		fmt.Println("Got error waiting for stack to be updated")
		fmt.Println(err)
	}
	return err
}

//helper method which statically generates CAPABILITY_IAM (a requirement for CloudFormation)
func createCapability() []*string {
	capabilities := make([]*string, 0)
	capIAM := "CAPABILITY_IAM"
	capabilities = append(capabilities, &capIAM)

	return capabilities
}

//helper method to convert a map of tags into cloudformation tags
func createTags(tagMap *map[string]string) []*cloudformation.Tag {
	tags := make([]*cloudformation.Tag, 0)
	for key, value := range *tagMap {
		tags = append(tags, &cloudformation.Tag{Key: &key, Value: &value})
	}
	return tags
}
