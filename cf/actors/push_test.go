package actors_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/cli/cf/actors"
	fakeBits "github.com/cloudfoundry/cli/cf/api/application_bits/fakes"
	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/cloudfoundry/cli/cf/app_files/fakes"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/gofileutils/fileutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Push Actor", func() {
	var (
		appBitsRepo  *fakeBits.FakeApplicationBitsRepository
		appFiles     *fakes.FakeAppFiles
		zipper       *fakes.FakeZipper
		actor        actors.PushActor
		fixturesDir  string
		appDir       string
		allFiles     []models.AppFileFields
		presentFiles []resources.AppFileResource
	)

	BeforeEach(func() {
		appBitsRepo = &fakeBits.FakeApplicationBitsRepository{}
		appFiles = &fakes.FakeAppFiles{}
		zipper = &fakes.FakeZipper{}
		actor = actors.NewPushActor(appBitsRepo, zipper, appFiles)
		fixturesDir = filepath.Join("..", "..", "fixtures", "applications")
	})

	Describe("GatherFiles", func() {
		BeforeEach(func() {
			allFiles = []models.AppFileFields{
				models.AppFileFields{Path: "example-app/.cfignore"},
				models.AppFileFields{Path: "example-app/app.rb"},
				models.AppFileFields{Path: "example-app/config.ru"},
				models.AppFileFields{Path: "example-app/Gemfile"},
				models.AppFileFields{Path: "example-app/Gemfile.lock"},
				models.AppFileFields{Path: "example-app/ignore-me"},
				models.AppFileFields{Path: "example-app/manifest.yml"},
			}

			presentFiles = []resources.AppFileResource{
				resources.AppFileResource{Path: "example-app/ignore-me"},
			}

			appDir = filepath.Join(fixturesDir, "example-app.zip")
			zipper.UnzipReturns(nil)
			appFiles.AppFilesInDirReturns(allFiles, nil)
			appBitsRepo.GetApplicationFilesReturns(presentFiles, nil)
		})

		AfterEach(func() {
		})

		Context("when the input is a zipfile", func() {
			BeforeEach(func() {
				zipper.IsZipFileReturns(true)
			})

			It("extracts the zip", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					_, _, err = actor.GatherFiles(appDir, tmpDir)
					Expect(err).NotTo(HaveOccurred())
					Expect(zipper.UnzipCallCount()).To(Equal(1))
				})
			})

			FIt("returns files list with file mode populated", func() {
				var expectedFileMode string

				zipper.UnzipStub = func(source string, dest string) error {
					err := os.Mkdir(filepath.Join(dest, "example-app"), os.ModeDir|os.ModePerm)
					Expect(err).NotTo(HaveOccurred())

					f, err := os.Create(filepath.Join(dest, "example-app/ignore-me"))
					Expect(err).NotTo(HaveOccurred())
					defer f.Close()

					err = ioutil.WriteFile(filepath.Join(dest, "example-app/ignore-me"), []byte("This is a test file"), os.ModePerm)

					info, err := os.Lstat(filepath.Join(dest, "example-app/ignore-me"))
					Expect(err).NotTo(HaveOccurred())

					expectedFileMode = fmt.Sprintf("%#o", info.Mode())

					return nil
				}

				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					actualFiles, _, err := actor.GatherFiles(appDir, tmpDir)
					Expect(err).NotTo(HaveOccurred())

					expectedFiles := []resources.AppFileResource{
						resources.AppFileResource{
							Path: "example-app/ignore-me",
							Mode: expectedFileMode,
						},
					}

					Expect(actualFiles).To(Equal(expectedFiles))
				})
			})

		})

		Context("when the input is a directory full of files", func() {
			BeforeEach(func() {
				zipper.IsZipFileReturns(false)
			})

			It("does not try to unzip the directory", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					files, _, err := actor.GatherFiles(appDir, tmpDir)
					Expect(zipper.UnzipCallCount()).To(Equal(0))
					Expect(err).NotTo(HaveOccurred())
					Expect(files).To(Equal(presentFiles))
				})
			})
		})

		Context("when errors occur", func() {
			It("returns an error if it cannot unzip the files", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					zipper.IsZipFileReturns(true)
					zipper.UnzipReturns(errors.New("error"))
					_, _, err = actor.GatherFiles(appDir, tmpDir)
					Expect(err).To(HaveOccurred())
				})
			})

			It("returns an error if it cannot walk the files", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					appFiles.AppFilesInDirReturns(nil, errors.New("error"))
					_, _, err = actor.GatherFiles(appDir, tmpDir)
					Expect(err).To(HaveOccurred())
				})
			})

			It("returns an error if we cannot reach the cc", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					appBitsRepo.GetApplicationFilesReturns(nil, errors.New("error"))
					_, _, err = actor.GatherFiles(appDir, tmpDir)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("when using .cfignore", func() {
			BeforeEach(func() {
				appBitsRepo.GetApplicationFilesReturns(nil, nil)
				appDir = filepath.Join(fixturesDir, "exclude-a-default-cfignore")
			})

			It("includes the .cfignore file in the upload directory", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					files, _, err := actor.GatherFiles(appDir, tmpDir)
					Expect(err).NotTo(HaveOccurred())

					_, err = os.Stat(filepath.Join(tmpDir, ".cfignore"))
					Expect(os.IsNotExist(err)).To(BeFalse())
					Expect(len(files)).To(Equal(0))
				})
			})
		})
	})

	Describe(".PopulateFileMode()", func() {
		var files []resources.AppFileResource

		BeforeEach(func() {
			files = []resources.AppFileResource{
				resources.AppFileResource{Path: "example-app/.cfignore"},
				resources.AppFileResource{Path: "example-app/app.rb"},
				resources.AppFileResource{Path: "example-app/config.ru"},
			}
		})

		It("returns []resources.AppFileResource with file mode populated", func() {
			actualFiles, err := actor.PopulateFileMode(fixturesDir, files)
			Ω(err).NotTo(HaveOccurred())

			for i, _ := range files {
				fileInfo, err := os.Lstat(filepath.Join(fixturesDir, files[i].Path))
				Ω(err).NotTo(HaveOccurred())
				mode := fileInfo.Mode()
				Ω(actualFiles[i].Mode).To(Equal(fmt.Sprintf("%#o", mode)))
			}
		})
	})

	Describe(".UploadApp", func() {
		It("Simply delegates to the UploadApp function on the app bits repo, which is not worth testing", func() {})
	})
})
