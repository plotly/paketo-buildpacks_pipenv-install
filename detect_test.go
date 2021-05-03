package pipenvinstall_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit"
	pipenvinstall "github.com/paketo-community/pipenv-install"
	"github.com/paketo-community/pipenv-install/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		detect     packit.DetectFunc
		lockParser *fakes.Parser
		workingDir string
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(workingDir, "Pipfile"), []byte{}, 0644)
		Expect(err).NotTo(HaveOccurred())

		lockParser = &fakes.Parser{}

		detect = pipenvinstall.Detect(lockParser)
	})

	context("detection", func() {
		it("returns a build plan that provides site-packages", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: pipenvinstall.SitePackages},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: pipenvinstall.CPython,
						Metadata: pipenvinstall.BuildPlanMetadata{
							Build: true,
						},
					},
					{
						Name: pipenvinstall.Pipenv,
						Metadata: pipenvinstall.BuildPlanMetadata{
							Build: true,
						},
					},
				},
			}))
			Expect(lockParser.ParseVersionCall.Receives.Path).To(Equal(workingDir))
		})

		context("when there is no Pipfile", func() {
			it.Before(func() {
				Expect(os.Remove(filepath.Join(workingDir, "Pipfile"))).To(Succeed())
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(packit.Fail))
			})
		})

		context("failure cases", func() {
			context("when the Pipfile cannot be read", func() {
				it.Before(func() {
					Expect(os.Chmod(workingDir, 0000)).To(Succeed())
				})

				it.After(func() {
					Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := detect(packit.DetectContext{
						WorkingDir: workingDir,
					})
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})
		})
	})
}
