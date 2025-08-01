package pod

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"
)

func TestDefinePod(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	assert.Equal(t, testPod.Name, "test-pod")
	assert.Equal(t, testPod.Namespace, "test-namespace")
	assert.Equal(t, testPod.Labels["app"], "nginx")
	assert.Equal(t, testPod.Spec.Containers[0].Image, "nginx")
	assert.Equal(t, ptr.To[int64](0), testPod.Spec.TerminationGracePeriodSeconds)
	assert.Equal(t, ptr.To[int64](1000), testPod.Spec.SecurityContext.RunAsUser)
	assert.Equal(t, ptr.To[int64](1000), testPod.Spec.SecurityContext.RunAsGroup)
	assert.Equal(t, ptr.To[bool](true), testPod.Spec.SecurityContext.RunAsNonRoot)
}

func TestRedefineWithServiceAccount(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithServiceAccount(testPod, "test-service-account")
	assert.Equal(t, testPod.Spec.ServiceAccountName, "test-service-account")
}

func TestRedefinePodContainerWithLivenessProbeCommand(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefinePodContainerWithLivenessProbeCommand(testPod, 0, []string{"ls"})
	assert.Equal(t, testPod.Spec.Containers[0].LivenessProbe.Exec.Command, []string{"ls"})
	assert.Equal(t, int32(5), testPod.Spec.Containers[0].LivenessProbe.InitialDelaySeconds)
	assert.Equal(t, int32(5), testPod.Spec.Containers[0].LivenessProbe.PeriodSeconds)
}

func TestRedefineWithLivenessProbe(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithLivenessProbe(testPod)
	assert.Equal(t, testPod.Spec.Containers[0].LivenessProbe.Exec.Command, []string{"ls"})
}

func TestRedefineWithStartUpProbe(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithStartUpProbe(testPod)
	assert.Equal(t, testPod.Spec.Containers[0].StartupProbe.Exec.Command, []string{"ls"})
}

func TestRedefineWithPVC(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithPVC(testPod, "test-pv", "test-pvc")
	assert.Equal(t, testPod.Spec.Volumes[0].Name, "test-pv")
	assert.Equal(t, testPod.Spec.Volumes[0].PersistentVolumeClaim.ClaimName, "test-pvc")
}

func TestRedefineWithCPUResources(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithCPUResources(testPod, "101m", "100m")
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Requests.Cpu().String(), "100m")
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Limits.Cpu().String(), "101m")
}

func TestRedefineWithMemoryResources(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithMemoryResources(testPod, "101Mi", "100Mi")
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Requests.Memory().String(), "100Mi")
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Limits.Memory().String(), "101Mi")
}

func TestRedefineWithRunTimeClass(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithRunTimeClass(testPod, "test-runtime-class")
	assert.Equal(t, *testPod.Spec.RuntimeClassName, "test-runtime-class")
}

func TestRedefineWithNodeAffinity(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithNodeAffinity(testPod, "key1")
	assert.Equal(t, testPod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.
		NodeSelectorTerms[0].MatchExpressions[0].Key, "key1")
}

func TestRedefineWithPodAffinity(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	podAffinity := map[string]string{"key1": "value1", "key2": "value2"}
	RedefineWithPodAffinity(testPod, podAffinity)
	assert.Equal(t, testPod.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].
		LabelSelector.MatchLabels["key1"], "value1")
	assert.Equal(t, testPod.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].
		LabelSelector.MatchLabels["key2"], "value2")
	assert.Equal(t, "kubernetes.io/hostname", testPod.Spec.Affinity.PodAffinity.
		RequiredDuringSchedulingIgnoredDuringExecution[0].TopologyKey)
}

func TestRedefineWithPodAntiAffinity(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	podAntiAffinity := map[string]string{"key1": "value1", "key2": "value2"}
	RedefineWithPodAntiAffinity(testPod, podAntiAffinity)
	assert.Equal(t, testPod.Spec.Affinity.PodAntiAffinity.
		RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchLabels["key1"], "value1")
	assert.Equal(t, testPod.Spec.Affinity.PodAntiAffinity.
		RequiredDuringSchedulingIgnoredDuringExecution[0].LabelSelector.MatchLabels["key2"], "value2")
	assert.Equal(t, "kubernetes.io/hostname", testPod.Spec.Affinity.
		PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution[0].TopologyKey)
}

func TestRedefineWithInfrastructureTolerations(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithInfrastructureTolerations(testPod)

	// Check that tolerations were added
	assert.Greater(t, len(testPod.Spec.Tolerations), 0)

	// Check for specific tolerations
	var hasDiskPressure, hasMemoryPressure bool

	for _, toleration := range testPod.Spec.Tolerations {
		if toleration.Key == "node.kubernetes.io/disk-pressure" {
			hasDiskPressure = true

			assert.Equal(t, corev1.TolerationOpExists, toleration.Operator)
			assert.Equal(t, corev1.TaintEffectNoSchedule, toleration.Effect)
		}

		if toleration.Key == "node.kubernetes.io/memory-pressure" {
			hasMemoryPressure = true

			assert.Equal(t, corev1.TolerationOpExists, toleration.Operator)
			assert.Equal(t, corev1.TaintEffectNoSchedule, toleration.Effect)
		}
	}

	assert.True(t, hasDiskPressure, "Should have disk pressure toleration")
	assert.True(t, hasMemoryPressure, "Should have memory pressure toleration")
}

func TestRedefineWithCustomTolerations(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	customTolerations := []corev1.Toleration{
		{
			Key:      "custom-taint",
			Operator: corev1.TolerationOpEqual,
			Value:    "custom-value",
			Effect:   corev1.TaintEffectNoExecute,
		},
	}

	RedefineWithCustomTolerations(testPod, customTolerations)

	assert.Equal(t, 1, len(testPod.Spec.Tolerations))
	assert.Equal(t, "custom-taint", testPod.Spec.Tolerations[0].Key)
	assert.Equal(t, "custom-value", testPod.Spec.Tolerations[0].Value)
}

func TestRedefineWithInfrastructureTolerationsIfEnabled(t *testing.T) {
	// Test with default (should be enabled since default is now true)
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithInfrastructureTolerationsIfEnabled(testPod)

	// Should have tolerations since default is true
	assert.Greater(t, len(testPod.Spec.Tolerations), 0, "Should have tolerations with default configuration")

	// Test disabled case
	t.Setenv("ENABLE_INFRASTRUCTURE_TOLERATIONS", "false")
	testPod2 := DefinePod("test-pod-2", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	RedefineWithInfrastructureTolerationsIfEnabled(testPod2)

	// Should NOT have tolerations when explicitly disabled
	assert.Equal(t, 0, len(testPod2.Spec.Tolerations), "Should not have tolerations when disabled")
}

func TestRedefineWith2MiHugepages(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	testPod.Spec.Containers[0].Resources.Requests = make(map[corev1.ResourceName]resource.Quantity)
	testPod.Spec.Containers[0].Resources.Limits = make(map[corev1.ResourceName]resource.Quantity)
	RedefineWith2MiHugepages(testPod, 2)
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Requests["hugepages-2Mi"], resource.MustParse("2Mi"))
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Limits["hugepages-2Mi"], resource.MustParse("2Mi"))
}

func TestRedefineWith1GiHugepages(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	testPod.Spec.Containers[0].Resources.Requests = make(map[corev1.ResourceName]resource.Quantity)
	testPod.Spec.Containers[0].Resources.Limits = make(map[corev1.ResourceName]resource.Quantity)
	RedefineWith1GiHugepages(testPod, 2)
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Requests["hugepages-1Gi"], resource.MustParse("2Gi"))
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Limits["hugepages-1Gi"], resource.MustParse("2Gi"))
}

func TestRedefineFirstContainerWith2MiHugepages(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	testPod.Spec.Containers[0].Resources.Requests = make(map[corev1.ResourceName]resource.Quantity)
	testPod.Spec.Containers[0].Resources.Limits = make(map[corev1.ResourceName]resource.Quantity)
	assert.Nil(t, RedefineFirstContainerWith2MiHugepages(testPod, 2))
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Requests["hugepages-2Mi"], resource.MustParse("2Mi"))
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Limits["hugepages-2Mi"], resource.MustParse("2Mi"))
}

func TestRedefineSecondContainerWith1GHugepages(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})
	testPod.Spec.Containers = append(testPod.Spec.Containers, corev1.Container{})
	testPod.Spec.Containers[1].Resources.Requests = make(map[corev1.ResourceName]resource.Quantity)
	testPod.Spec.Containers[1].Resources.Limits = make(map[corev1.ResourceName]resource.Quantity)
	assert.Nil(t, RedefineSecondContainerWith1GHugepages(testPod, 2))
	assert.Equal(t, testPod.Spec.Containers[1].Resources.Requests["hugepages-1Gi"], resource.MustParse("2Gi"))
	assert.Equal(t, testPod.Spec.Containers[1].Resources.Limits["hugepages-1Gi"], resource.MustParse("2Gi"))
}

func TestRedefineWithPostStart(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})

	RedefineWithPostStart(testPod)
	assert.Equal(t, testPod.Spec.Containers[0].Lifecycle.PostStart.Exec.Command, []string{"ls"})
}

func TestRedefineWithContainerExecCommand(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})

	assert.Nil(t, RedefineWithContainerExecCommand(testPod, []string{"ls"}, 0))
	assert.Equal(t, testPod.Spec.Containers[0].Command, []string{"ls"})
}

func TestRedefineWithReadinessProbe(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})

	RedefineWithReadinessProbe(testPod)
	assert.Equal(t, testPod.Spec.Containers[0].ReadinessProbe.Exec.Command, []string{"ls"})
}

func TestRedefineFirstContainerWith1GiHugepages(t *testing.T) {
	testPod := DefinePod("test-pod", "test-namespace", "nginx", map[string]string{"app": "nginx"})

	assert.Nil(t, RedefineFirstContainerWith1GiHugepages(testPod, 2))
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Requests["hugepages-1Gi"], resource.MustParse("2Gi"))
	assert.Equal(t, testPod.Spec.Containers[0].Resources.Limits["hugepages-1Gi"], resource.MustParse("2Gi"))
}
