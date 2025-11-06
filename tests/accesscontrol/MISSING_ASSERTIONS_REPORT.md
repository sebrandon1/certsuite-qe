# Access Control Suite - Missing Assertions Report

## Overview
This report identifies test cases (It blocks) in the accesscontrol suite that do NOT have assertions verifying the scenario setup before calling `LaunchTests()`.

## Purpose
Before calling `LaunchTests()`, test cases should have assertions (`Expect()` statements) that verify the test environment is properly configured. This ensures:
1. The test scenario is actually set up as intended
2. Failures can be traced to the actual test vs setup issues
3. Test cases are more maintainable and self-documenting

---

## ✅ PROGRESS UPDATE

**Files Fixed (Assertions Added):**
1. ✅ **access_control_bpf_capability_check.go** - All 4 tests now have assertions
2. ✅ **access_control_one_process_per_container.go** - All 4 tests now have assertions
3. ✅ **access_control_pod_service_account.go** - All 3 tests now have assertions
4. ✅ **access_control_ssh_daemons.go** - Both tests now have assertions

**Total Fixed: 15 test cases**

---

## Files with Missing Assertions

### 1. access_control_bpf_capability_check.go
**ALL test cases missing assertions** (4/4 tests)

- **Line 40-60**: "one deployment, one pod, one container, does not have BPF capability"
  - Creates deployment but doesn't assert BPF is NOT present
  
- **Line 62-84**: "one deployment, one pod, one container, does have BPF capability [negative]"
  - Creates deployment with BPF but doesn't assert it's present
  
- **Line 86-114**: "two deployments, one pod each, one container each, does not have BPF capability"
  - Creates two deployments but doesn't assert BPF is NOT present
  
- **Line 116-146**: "two deployments, one pod each, one container each, one does have BPF capability [negative]"
  - Creates deployments with varying BPF but doesn't assert the configuration

---

### 2. access_control_ssh_daemons.go
**1/2 test cases missing assertions**

- **Line 77-101**: "one pod with ssh daemon running [negative]"
  - Creates pod with SSH daemon command but doesn't verify it's actually running
  - **Compare to**: Line 42-75 which DOES have assertions checking for absence of SSH processes

---

### 3. access_control_namespace_resource_quota.go
**ALL test cases missing assertions** (4/4 tests)

- **Line 51-79**: "one deployment, one pod in a namespace with resource quota"
  - Creates deployment and resource quota but doesn't assert quota exists
  
- **Line 82-102**: "one deployment, one pod in a namespace without resource quota [negative]"
  - Creates deployment but doesn't assert absence of resource quota
  
- **Line 105-150**: "two deployments, one pod each, both in a namespace with resource quota"
  - Creates deployments and quotas but doesn't verify they exist
  
- **Line 153-190**: "two deployments, one pod each, one in a namespace without resource quota [negative]"
  - Creates deployments with mixed quota scenarios but doesn't verify

---

### 4. access_control_one_process_per_container.go
**ALL test cases missing assertions** (4/4 tests)

- **Line 46-66**: "one deployment, one pod, one container, only one process"
  - Creates deployment but doesn't verify process count
  
- **Line 68-90**: "one deployment, one pod, one container, two processes [negative]"
  - Defines command to launch two processes but doesn't verify they're running
  
- **Line 92-114**: "one deployment, one pod, two containers, one process each"
  - Creates deployment with two containers but doesn't verify process counts
  
- **Line 116-141**: "one deployment, one pod, two containers, the second one with two processes [negative]"
  - Creates deployment but doesn't verify process configuration

---

### 5. access_control_security_context.go
**1/4 test cases missing assertions**

- **Line 104-135**: "two deployments, one pod each, one container each, both have allowed security context"
  - Creates two deployments but doesn't assert their security contexts

---

### 6. access_control_namespace.go
**5/10 test cases missing assertions**

- **Line 36-58**: "one namespace, no invalid prefixes"
  - Only defines config, no assertions about namespace validity
  
- **Line 96-128**: "two namespaces, no invalid prefixes"
  - Creates namespace but doesn't assert it's valid
  
- **Line 176-205**: "one custom resource in a valid namespace"
  - Creates InstallPlan but doesn't assert it exists in correct namespace
  
- **Line 250-296**: "two custom resources, both in valid namespaces"
  - Creates two InstallPlans but doesn't verify namespaces
  
- **Line 348-387**: "two custom resources of different CRDs, both in valid namespace"
  - Creates InstallPlan and Subscription but no assertions about namespaces

---

### 7. access_control_crd_roles.go
**ALL test cases missing assertions** (4/4 tests)

- **Line 57-84**: "Custom resource is deployed, proper role defined"
  - Creates custom resource and role but doesn't verify role configuration
  
- **Line 86-113**: "Custom resource is deployed, one role defined with multiple api groups [negative]"
  - Creates role with multiple API groups but doesn't assert configuration
  
- **Line 115-142**: "Custom resource is deployed, one role with multiple resources defined [negative]"
  - Creates role with multiple resources but doesn't verify
  
- **Line 144-171**: "Custom resource is deployed, with improper role [skip]"
  - Creates role with bad API group but doesn't assert

---

### 8. access_control_pod_service_account.go
**ALL test cases missing assertions** (3/3 tests)

- **Line 38-64**: "one pod with valid service account"
  - Creates service account and pod but doesn't assert service account is set
  
- **Line 66-86**: "one pod with empty service account [negative]"
  - Creates pod with empty SA but doesn't verify
  
- **Line 88-108**: "one pod with default service account [negative]"
  - Creates pod with default SA but doesn't verify

---

### 9. access_control_pod_role_bindings.go
**ALL test cases missing assertions** (4/4 tests)

- **Line 57-79**: "one pod with valid role binding"
  - Creates pod with service account but doesn't assert role binding exists
  
- **Line 81-101**: "one pod with no specified service account (default SA) [negative]"
  - Creates pod but doesn't verify it uses default SA
  
- **Line 103-143**: "one pod with service account in different namespace"
  - Complex setup with namespace changes but no assertions
  
- **Line 145-185**: "one pod with role binding in different namespace"
  - Creates role binding in different namespace but doesn't verify

---

### 10. access_control_pod_cluster_role_bindings.go
**1/2 test cases missing assertions**

- **Line 41-61**: "one deployment, one pod, does not have cluster role binding"
  - Creates deployment but doesn't assert absence of cluster role binding
  - **Compare to**: Line 63-110 which DOES have assertions verifying service account

---

## Summary Statistics

| File | Tests Missing Assertions | Total Tests | Percentage |
|------|-------------------------|-------------|------------|
| access_control_bpf_capability_check.go | 4 | 4 | 100% |
| access_control_namespace_resource_quota.go | 4 | 4 | 100% |
| access_control_one_process_per_container.go | 4 | 4 | 100% |
| access_control_crd_roles.go | 4 | 4 | 100% |
| access_control_pod_service_account.go | 3 | 3 | 100% |
| access_control_pod_role_bindings.go | 4 | 4 | 100% |
| access_control_namespace.go | 5 | 10 | 50% |
| access_control_ssh_daemons.go | 1 | 2 | 50% |
| access_control_security_context.go | 1 | 4 | 25% |
| access_control_pod_cluster_role_bindings.go | 1 | 2 | 50% |
| **TOTAL** | **31** | **41** | **75.6%** |

---

## Good Examples (Tests WITH Assertions)

These test files serve as good examples of proper assertion patterns:

### access_control_container_host_port.go
- ✅ All tests have assertions before LaunchTests
- Example pattern (lines 51-54):
```go
By("Assert deployment has no host port configured")
runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
Expect(err).ToNot(HaveOccurred())
Expect(runningDeployment.Spec.Template.Spec.Containers[0].Ports).To(BeEmpty())
```

### access_control_net_admin_capability_check.go
- ✅ All tests have assertions before LaunchTests
- Example pattern (lines 51-54):
```go
By("Assert deployment does not have NET_ADMIN capability")
runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
Expect(err).ToNot(HaveOccurred())
Expect(runningDeployment.Spec.Template.Spec.Containers[0].SecurityContext).To(BeNil())
```

### access_control_ssh_daemons.go (first test)
- ✅ Good example of verifying scenario (lines 51-62):
```go
By("Assert pod does not have ssh daemon running")
runningPod, err := globalhelper.GetRunningPod(randomNamespace, tsparams.TestPodName)
Expect(err).ToNot(HaveOccurred())

// Use the exec command to check if the ssh daemon is running
output, err := globalhelper.ExecCommand(*runningPod, []string{"ps", "-ef"})
Expect(err).ToNot(HaveOccurred())
Expect(output.String()).To(Not(ContainSubstring("ssh-agent")))
Expect(output.String()).To(Not(ContainSubstring("ssh-add")))
// ... more checks
```

---

## Recommended Action Items

1. **Priority 1 (Critical)**: Add assertions to files with 100% missing assertions:
   - access_control_bpf_capability_check.go
   - access_control_namespace_resource_quota.go
   - access_control_one_process_per_container.go
   - access_control_crd_roles.go
   - access_control_pod_service_account.go
   - access_control_pod_role_bindings.go

2. **Priority 2 (High)**: Add assertions to remaining test cases in:
   - access_control_namespace.go
   - access_control_ssh_daemons.go
   - access_control_security_context.go
   - access_control_pod_cluster_role_bindings.go

3. **Create Jira stories** to track this work systematically

4. **Use existing patterns**: Reference the "Good Examples" section for assertion patterns to follow

---

Generated: $(date)

