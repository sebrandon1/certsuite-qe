package helper

import (
	"errors"
	"testing"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

const (
	testDeploymentName = "test-deployment"
)

func assertDeployment(t *testing.T, deployment *appsv1.Deployment) {
	t.Helper()
	assert.Equal(t, testDeploymentName, deployment.Name)
	assert.Equal(t, int32(1), *deployment.Spec.Replicas)
	assert.Equal(t, 1, len(deployment.Spec.Template.Spec.Containers))
	assert.Equal(t, parameters.TestDeploymentLabels, deployment.Spec.Template.Labels)
	assert.Equal(t, parameters.SampleWorkloadImage, deployment.Spec.Template.Spec.Containers[0].Image)
}

func TestDefineDeployment(t *testing.T) {
	// Define a deployment with 1 replica and 1 container
	deployment, err := DefineDeployment(1, 1, testDeploymentName, parameters.TestAccessControlNameSpace)
	assert.Nil(t, err)

	assertDeployment(t, deployment)
	assert.Equal(t, parameters.TestAccessControlNameSpace, deployment.Namespace)
	assert.Equal(t, parameters.TestAccessControlNameSpace, deployment.Spec.Template.Namespace)
}

func TestDefineDeploymentWithClusterRoleBindingWithServiceAccount(t *testing.T) {
	// Define a deployment with 1 replica and 1 container
	deployment, err := DefineDeploymentWithClusterRoleBindingWithServiceAccount(1, 1, testDeploymentName,
		parameters.TestAccessControlNameSpace, "my-service-account")
	assert.Nil(t, err)

	assertDeployment(t, deployment)
	assert.Equal(t, parameters.TestAccessControlNameSpace, deployment.Namespace)
	assert.Equal(t, parameters.TestAccessControlNameSpace, deployment.Spec.Template.Namespace)
}

func TestDefineDeploymentWithNamespace(t *testing.T) {
	// Define a deployment with 1 replica and 1 container
	deployment, err := DefineDeploymentWithNamespace(1, 1, testDeploymentName, "test-namespace")
	assert.Nil(t, err)

	assertDeployment(t, deployment)
	assert.Equal(t, "test-namespace", deployment.Namespace)
	assert.Equal(t, "test-namespace", deployment.Spec.Template.Namespace)
}

func TestDefineDeploymentWithContainerPorts(t *testing.T) {
	// Define a deployment with 1 replica and 1 container
	deployment, err := DefineDeploymentWithContainerPorts(testDeploymentName,
		parameters.TestAccessControlNameSpace, 1, []corev1.ContainerPort{
			{
				Name:          "test-port",
				ContainerPort: 80,
			},
		})
	assert.Nil(t, err)

	assertDeployment(t, deployment)
	assert.Equal(t, parameters.TestAccessControlNameSpace, deployment.Namespace)
	assert.Equal(t, parameters.TestAccessControlNameSpace, deployment.Spec.Template.Namespace)
	assert.Equal(t, 1, len(deployment.Spec.Template.Spec.Containers[0].Ports))
	assert.Equal(t, "test-port", deployment.Spec.Template.Spec.Containers[0].Ports[0].Name)
	assert.Equal(t, int32(80), deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort)
}

func TestSetServiceAccountAutomountServiceAccountToken(t *testing.T) {
	// generate test serviceaccount objects
	generateServiceAccount := func() *corev1.ServiceAccount {
		return &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-service-account",
				Namespace: parameters.TestAccessControlNameSpace,
			},
		}
	}

	testCases := []struct {
		testValue   string
		expectedErr error
	}{
		{
			testValue:   "true",
			expectedErr: nil,
		},
		{
			testValue:   "false",
			expectedErr: nil,
		},
		{
			testValue:   "invalid",
			expectedErr: errors.New("invalid value for token value"),
		},
		{
			testValue:   "nil",
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		client := k8sfake.NewSimpleClientset(generateServiceAccount())

		// Set the globalhelper client to the fake client
		globalhelper.SetTestK8sAPIClient(client)

		// Set the automountServiceAccountToken to the test value
		err := SetServiceAccountAutomountServiceAccountToken(parameters.TestAccessControlNameSpace, "my-service-account", testCase.testValue)

		// Check if the error is as expected
		assert.Equal(t, testCase.expectedErr, err)

		// Get the serviceaccount from the fake client and check if the automountServiceAccountToken is set to the test value
		serviceAccount, err := client.CoreV1().
			ServiceAccounts(parameters.TestAccessControlNameSpace).Get(t.Context(), "my-service-account", metav1.GetOptions{})
		assert.Nil(t, err)

		if err == nil {
			switch testCase.testValue {
			case "true":
				assert.Equal(t, true, *serviceAccount.AutomountServiceAccountToken)
			case "false":
				assert.Equal(t, false, *serviceAccount.AutomountServiceAccountToken)
			case "nil":
				assert.Nil(t, serviceAccount.AutomountServiceAccountToken)
			}
		}

		globalhelper.UnsetTestK8sAPIClient()
	}
}

func TestDefineAndCreateServiceOnCluster(t *testing.T) {
	testCases := []struct {
		ipFamilyDesired string
		ipFams          []corev1.IPFamily
		withNodePort    bool
	}{
		{
			ipFamilyDesired: "IPv4",
			ipFams:          []corev1.IPFamily{corev1.IPv4Protocol},
			withNodePort:    false,
		},
		{
			ipFamilyDesired: "IPv6",
			ipFams:          []corev1.IPFamily{corev1.IPv6Protocol},
			withNodePort:    false,
		},
		{
			ipFamilyDesired: "IPv4",
			ipFams:          []corev1.IPFamily{corev1.IPv4Protocol},
			withNodePort:    true,
		},
		{
			ipFamilyDesired: "IPv6",
			ipFams:          []corev1.IPFamily{corev1.IPv6Protocol},
			withNodePort:    true,
		},
		{
			ipFamilyDesired: "",
			ipFams:          []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol},
			withNodePort:    false,
		},
	}

	for _, testCase := range testCases {
		// Create a fake clientset
		client := k8sfake.NewSimpleClientset()

		// Set the globalhelper client to the fake client
		globalhelper.SetTestK8sAPIClient(client)

		// Define a service
		err := DefineAndCreateServiceOnCluster("service-name", parameters.TestAccessControlNameSpace,
			8080, 8080, testCase.withNodePort, testCase.ipFams, testCase.ipFamilyDesired)
		assert.Nil(t, err)

		// Get the service from the fake client and check if it exists
		_, err = client.CoreV1().Services(parameters.TestAccessControlNameSpace).Get(t.Context(), "service-name", metav1.GetOptions{})
		assert.Nil(t, err)

		globalhelper.UnsetTestK8sAPIClient()
	}
}
