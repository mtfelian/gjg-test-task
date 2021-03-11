package storage_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/mtfelian/gjg-test-task/config"
	"github.com/mtfelian/gjg-test-task/service"
	"github.com/mtfelian/gjg-test-task/storage"
	"github.com/mtfelian/gjg-test-task/storage/model"
	"github.com/mtfelian/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestAll(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

// IsInDocker returns true if we are in Docker container
func IsInDocker() bool { return utils.FileExists("/.dockerenv") }

var _ = Describe("testing service API", func() {
	var (
		s                *service.Service
		dbConnectionData = storage.PostgresConnection{
			User:     "postgres",
			Password: "postgres",
			Host:     "127.0.0.1",
			Port:     "5440",
			Name:     "gjg-test-task",
			Schema:   "public",
		}
	)
	BeforeSuite(func() {
		if IsInDocker() {
			dbConnectionData.Host = "postgres"
			dbConnectionData.Password = "postgres"
			dbConnectionData.Port = "5432"
		}

		viper.Set(config.LogLevel, logrus.DebugLevel.String())
		dbConnectionData.SetIntoViper(viper.GetViper())

		Expect(service.NewWithPostgresClient(viper.GetViper())).To(Succeed())
		s = service.Get()
		Expect(s.Storage).NotTo(BeNil())
	})

	Describe("levels storage", func() {
		BeforeEach(func() {
			Expect(s.Storage.ApplyMigrations("/migrations", "up")).To(Succeed())
			Expect(s.Storage.RemoveAll()).To(Succeed())
		})
		AfterEach(func() {})

		It("checks adding level", func() {
			var ids []strfmt.UUID
			newLevels := []model.Level{
				{X: 3, Y: 2, Maze: []byte{
					0, 1, 0,
					0, 0, 4,
				}},
				{X: 3, Y: 4, Maze: []byte{
					0, 1, 0,
					0, 0, 4,
					0, 1, 0,
					0, 1, 0,
				}},
				{X: 4, Y: 4, Maze: []byte{
					0, 1, 0, 0,
					0, 0, 0, 1,
					0, 1, 0, 1,
					0, 1, 0, 1,
				}},
			}
			By("creating newLevels", func() {
				ids = make([]strfmt.UUID, len(newLevels))
				var err error
				for i, level := range newLevels {
					ids[i], err = s.Storage.AddLevel(level)
					Expect(err).NotTo(HaveOccurred())
					Expect(ids[i]).NotTo(BeEmpty())
				}
			})

			By("listing all levels", func() {
				levels, err := s.Storage.GetLevels(model.GetLevelsParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(levels).To(HaveLen(3))
			})
		})
	})
})
