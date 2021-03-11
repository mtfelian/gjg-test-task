package main_test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	main "github.com/mtfelian/gjg-test-task"
	"github.com/mtfelian/gjg-test-task/api"
	"github.com/mtfelian/gjg-test-task/config"
	"github.com/mtfelian/gjg-test-task/game"
	"github.com/mtfelian/gjg-test-task/gpr"
	"github.com/mtfelian/gjg-test-task/service"
	"github.com/mtfelian/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// IsInDocker returns true if we are in Docker container
func IsInDocker() bool { return utils.FileExists("/.dockerenv") }

var (
	server *httptest.Server
	g      *gpr.GPR
)

func TestAll(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite (v2)")
}

var _ = Describe("Testing manager API", func() {
	BeforeSuite(func() {
		if IsInDocker() {
			viper.Set(config.DBHost, "postgres")
			viper.Set(config.DBPassword, "postgres")
			viper.Set(config.DBPort, "5432")
		}
		viper.Set(config.LogLevel, logrus.DebugLevel.String())
		Expect(service.NewWithPostgresClient(viper.GetViper())).To(Succeed())

		s := service.Get()
		httpServer := s.HTTPServer
		main.RegisterHTTPAPIHandlers(httpServer)
		server = httptest.NewServer(httpServer)

		g = gpr.New(server)
		g.SetLogPerformRequest(true)

		// set server port if needed
		URL, err := url.Parse(server.URL)
		Expect(err).NotTo(HaveOccurred())
		port, err := strconv.Atoi(URL.Port())
		Expect(err).NotTo(HaveOccurred())
		viper.Set(config.Port, uint(port))

		Expect(s.Storage).NotTo(BeNil())
	})

	AfterSuite(func() { server.Close() })

	Context("api.SubmitLevel request", func() {
		It("checks that creating levels works OK in valid cases", func() {
			tcs := [][][]byte{
				{
					{0, 1, 0, 0, 2},
					{0, 0, 4, 1, 0},
					{0, 1, 0, 1, 1},
					{0, 1, 0, 1, 2},
				},
				{
					{0, 1, 0},
					{0, 0, 4},
					{0, 1, 0},
					{0, 1, 0},
				},
			}
			ids := make([]strfmt.UUID, len(tcs))
			for i, tc := range tcs {
				p := api.SubmitLevelParams{Maze: tc}
				var r api.SubmitLevelResponse
				g.PerformSubmitLevelRequest(utils.MushMarshalJSON(p), http.StatusCreated, &r)
				Expect(r.LevelID).NotTo(BeEmpty(), "case %d", i)
				ids[i] = r.LevelID
			}
		})

		It("checks that creating level fails, invalid case: just wrong data", func() {
			var r game.Error
			g.PerformSubmitLevelRequest([]byte(`{"maze":"q"}`), http.StatusUnprocessableEntity, &r)
			Expect(r.Code).To(Equal(service.ErrValidationRequest))
		})

		It("checks that creating level fails, invalid case: related to validation (not rectangular maze)", func() {
			p := api.SubmitLevelParams{Maze: [][]byte{
				{0, 1, 0, 0, 2},
				{0, 0, 4, 1, 0},
				{0, 1, 0, 1},
				{0, 1, 0, 1, 2},
			}}
			var r game.Error
			g.PerformSubmitLevelRequest(utils.MushMarshalJSON(p), http.StatusBadRequest, &r)
			Expect(r.Code).To(Equal(service.ErrValidationFieldIsNotRectangular))
		})
	})
})
