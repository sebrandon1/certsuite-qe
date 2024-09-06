package operator

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/parameters"
)

const (
	ErrorDeployOperatorStr   = "Error deploying operator "
	ErrorLabelingOperatorStr = "Error labeling operator "
	ErrorRemovingLabelStr    = "Error removing label from operator "
)

var _ = Describe("Operator install-source,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.OperatorNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			tsparams.CertsuiteTargetCrdFilters, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		// Install 3 separate operators for testing
		By("Deploy operator group")
		err = tshelper.DeployTestOperatorGroup(randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	It("deploy cluster-wide cluster-logging operator with subscription in another namespace", func() {
		const (
			clusterLoggingOperatorName  = "cluster-logging"
			clusterLoggingTestNamespace = "fake-openshift-logging"
			openshiftLoggingNamespace   = "openshift-logging"
			subscriptionChannel         = "stable-5.9"
		)

		By("Create openshift-logging namespace")
		err := globalhelper.CreateNamespace(openshiftLoggingNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create fake operator group for cluster-logging operator")
		err = tshelper.DeployTestOperatorGroup(openshiftLoggingNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		DeferCleanup(func() {
			By("Delete fake namespace for cluster-logging operator")
			err := globalhelper.DeleteNamespaceAndWait(openshiftLoggingNamespace, tsparams.Timeout)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Query the packagemanifest for the " + clusterLoggingOperatorName)
		version, err := globalhelper.QueryPackageManifestForVersion(clusterLoggingOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for nginx-ingress-operator")

		By("Deploy cluster-logging operator for testing")
		err = tshelper.DeployOperatorSubscription(
			clusterLoggingOperatorName,
			subscriptionChannel,
			openshiftLoggingNamespace,
			tsparams.RedhatOperatorGroup,
			tsparams.OperatorSourceNamespace,
			clusterLoggingOperatorName+".v"+version,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+clusterLoggingOperatorName)

		By("Wait until operator is ready")
		err = tshelper.WaitUntilOperatorIsReady(clusterLoggingOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+clusterLoggingOperatorName+" is not ready")

		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				clusterLoggingOperatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+clusterLoggingOperatorName)

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66142
	It("one operator installed with OLM", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixCloudbees)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66143
	It("one operator not installed with OLM [negative]", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixOpenvino,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixOpenvino)

		By("Delete operator's subscription")
		err := globalhelper.DeleteSubscription(randomNamespace,
			tsparams.SubscriptionNameOpenvino)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66144
	It("two operators, both installed with OLM", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixCloudbees,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixCloudbees)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixAnchore)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66145
	It("two operators, one not installed with OLM [negative]", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixAnchore,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixAnchore)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixOpenvino,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixOpenvino)

		By("Delete operator's subscription")
		err := globalhelper.DeleteSubscription(randomNamespace,
			tsparams.SubscriptionNameOpenvino)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallSource,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallSource,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})
