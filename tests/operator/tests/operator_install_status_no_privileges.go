package operator

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/parameters"
)

var _ = Describe("Operator install-status-no-privileges,", Serial, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string
	var operatorName string
	var catalogSource string

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

		By("Deploy operator group")
		err = tshelper.DeployTestOperatorGroup(randomNamespace, false)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operator group")

		// grafana operator has clusterPermissions but no resourceNames
		By("Query the packagemanifest for grafana operator package name and catalog source")
		operatorName, catalogSource = globalhelper.CheckOperatorExistsOrFail("grafana", randomNamespace)

		By("Query the packagemanifest for available channel, version and CSV for " + operatorName)
		channel, _, csvName := globalhelper.CheckOperatorChannelAndVersionOrFail(operatorName, randomNamespace)

		By("Deploy grafana operator for testing")
		err = tshelper.DeployOperatorSubscription(
			operatorName,
			operatorName,
			channel,
			randomNamespace,
			catalogSource,
			tsparams.OperatorSourceNamespace,
			csvName,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			operatorName)

		err = tshelper.WaitUntilOperatorIsReady(operatorName,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName+
			" is not ready")

		// postgresql operator has no clusterPermissions
		By("Query the packagemanifest for postgresql operator package name and catalog source")
		postgresqlOperatorName, catalogSource2, err := globalhelper.QueryPackageManifestForOperatorNameAndCatalogSource(
			tsparams.OperatorPackageNamePrefixLightweight, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for postgresql operator")
		Expect(postgresqlOperatorName).ToNot(Equal("not found"), "postgresql operator package not found")
		Expect(catalogSource2).ToNot(Equal("not found"), "postgresql operator catalog source not found")

		By("Query the packagemanifest for available channel, version and CSV for " + postgresqlOperatorName)
		channel2, version2, csvName2, err := globalhelper.QueryPackageManifestForAvailableChannelVersionAndCSV(
			postgresqlOperatorName, randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Error querying package manifest for "+postgresqlOperatorName)
		Expect(channel2).ToNot(Equal("not found"), "Channel not found")
		Expect(version2).ToNot(Equal("not found"), "Version not found")
		Expect(csvName2).ToNot(Equal("not found"), "CSV name not found")

		By(fmt.Sprintf("Deploy postgresql operator (channel %s, version %s) for testing", channel2, version2))
		err = tshelper.DeployOperatorSubscription(
			postgresqlOperatorName,
			postgresqlOperatorName,
			channel2,
			randomNamespace,
			catalogSource2,
			tsparams.OperatorSourceNamespace,
			csvName2,
			v1alpha1.ApprovalAutomatic,
		)
		Expect(err).ToNot(HaveOccurred(), ErrorDeployOperatorStr+
			postgresqlOperatorName)

		err = tshelper.WaitUntilOperatorIsReady(tsparams.OperatorPrefixLightweight,
			randomNamespace)
		Expect(err).ToNot(HaveOccurred(), "Operator "+csvName2+
			" is not ready")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// 66381
	It("one operator with no clusterPermissions", func() {
		By("Label operator")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixLightweight,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixLightweight)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66383
	It("one operator with clusterPermissions [negative]", func() {
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 66384
	It("two operators, one with no clusterPermissions and one with clusterPermissions", func() {
		By("Label operators")
		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				tsparams.OperatorPrefixLightweight,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+tsparams.OperatorPrefixLightweight)

		Eventually(func() error {
			return tshelper.AddLabelToInstalledCSV(
				operatorName,
				randomNamespace,
				tsparams.OperatorLabel)
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(Not(HaveOccurred()),
			ErrorLabelingOperatorStr+operatorName)

		By("Start test")
		err := globalhelper.LaunchTests(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
			randomReportDir,
			randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteOperatorInstallStatusNoPrivileges,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

})
