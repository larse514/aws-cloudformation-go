package cf

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

const (
	stackName    = "STACKNAME"
	templateBody = "URL"
)

//CreateStack tests, mocks, and methods

// Define a mock to return a basic success
type mockGoodCloudFormationClient struct {
	cloudformationiface.CloudFormationAPI
}

func (m *mockGoodCloudFormationClient) CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return nil, nil
}

func (m *mockGoodCloudFormationClient) WaitUntilStackCreateComplete(*cloudformation.DescribeStacksInput) error {
	return nil
}
func (m *mockGoodCloudFormationClient) UpdateStack(*cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	return &cloudformation.UpdateStackOutput{}, nil
}

// Define a mock to return an error.
type mockBadCloudFormationClient struct {
	cloudformationiface.CloudFormationAPI
}

func (m *mockBadCloudFormationClient) CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return nil, errors.New("Bad Error")
}
func (m *mockBadCloudFormationClient) UpdateStack(*cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	return &cloudformation.UpdateStackOutput{}, errors.New("Bad Error")
}
func (m *mockBadCloudFormationClient) WaitUntilStackCreateComplete(*cloudformation.DescribeStacksInput) error {
	return errors.New("THIS IS AN ERROR")
}
func TestCloudformationCreateStackFromS3(t *testing.T) {
	executor := IaaSExecutor{Client: &mockGoodCloudFormationClient{}}
	m := map[string]string{
		"KEY": "VALUE",
	}
	err := executor.CreateStackFromS3(templateBody, stackName, &m, nil)
	if err != nil {
		t.Log("Successful stack request return error ", err.Error())
		t.Fail()
	}

}

func TestCloudformationCreateStackFromS3Fails(t *testing.T) {
	executor := IaaSExecutor{Client: &mockBadCloudFormationClient{}}
	m := map[string]string{
		"KEY": "VALUE",
	}
	err := executor.CreateStackFromS3(templateBody, stackName, &m, nil)
	if err == nil {
		t.Log("Error should have been returned")
		t.Fail()
	}

}
func TestCloudformationCreateStack(t *testing.T) {
	executor := IaaSExecutor{Client: &mockGoodCloudFormationClient{}}
	m := map[string]string{
		"KEY": "VALUE",
	}
	err := executor.CreateStack(templateBody, stackName, &m, nil)
	if err != nil {
		t.Log("Successful stack request return error ", err.Error())
		t.Fail()
	}

}

func TestCloudformationCreateStackFails(t *testing.T) {
	executor := IaaSExecutor{Client: &mockBadCloudFormationClient{}}
	m := map[string]string{
		"KEY": "VALUE",
	}
	err := executor.CreateStack(templateBody, stackName, &m, nil)
	if err == nil {
		t.Log("Error should have been returned")
		t.Fail()
	}

}
func TestCloudformationUpdateStack(t *testing.T) {
	executor := IaaSExecutor{Client: &mockGoodCloudFormationClient{}}
	m := map[string]string{
		"KEY": "VALUE",
	}
	err := executor.UpdateStack(templateBody, stackName, &m, nil)
	if err != nil {
		t.Log("Successful stack request return error ", err.Error())
		t.Fail()
	}

}

func TestCloudformationUpdateStackFails(t *testing.T) {
	executor := IaaSExecutor{Client: &mockBadCloudFormationClient{}}
	m := map[string]string{
		"KEY": "VALUE",
	}
	err := executor.UpdateStack(templateBody, stackName, &m, nil)
	if err == nil {
		t.Log("Error should have been returned")
		t.Fail()
	}

}

//PauseUntilCreateFinished tests, mocks, and methods

func TestCloudformationWaitUntilStackCreateComplete(t *testing.T) {
	executor := IaaSExecutor{Client: &mockGoodCloudFormationClient{}}

	err := executor.PauseUntilCreateFinished(stackName)
	if err != nil {
		t.Log("Successful stack request return error")
		t.Fail()
	}

}
func TestCloudformationWaitUntilStackCreateCompleteFails(t *testing.T) {
	executor := IaaSExecutor{Client: &mockBadCloudFormationClient{}}

	err := executor.PauseUntilCreateFinished(stackName)
	if err == nil {
		t.Log("Error should have been returned")
		t.Fail()
	}

}

func TestCreateTagsLength(t *testing.T) {
	m := map[string]string{
		"TAG": "VALUE",
	}
	tags := createTags(&m)

	if len(tags) != 1 {
		t.Log("invalid number of tags")
		t.Fail()
	}
}
func TestCreateTagsKey(t *testing.T) {
	m := map[string]string{
		"TAG": "VALUE",
	}
	tags := createTags(&m)
	key := *tags[0].Key
	if key != "TAG" {
		t.Log("invalid number key expected TAG", " found ", key)
		t.Fail()
	}
}
func TestCreateTagsValue(t *testing.T) {
	m := map[string]string{
		"TAG": "VALUE",
	}
	tags := createTags(&m)
	value := *tags[0].Value
	if value != "VALUE" {
		t.Log("invalid number value expected VALUE found ", value)
		t.Fail()
	}
}
