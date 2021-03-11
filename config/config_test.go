package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("checking config", func() {
	It("checks loading from file", func() {
		const (
			port = 3000
		)
		contents := []byte(fmt.Sprintf(`{%s:%d}`, strconv.Quote(Port), port))
		wd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		const testConfig = "test"
		viper.Set(FileName, testConfig)

		fileName := filepath.Join(wd, viper.GetString(FileName))
		fileNameWithExt := fileName + ".json"
		Expect(ioutil.WriteFile(fileNameWithExt, contents, 0644)).To(Succeed())
		defer func() { Expect(os.Remove(fileNameWithExt)).To(Succeed()) }()
		Expect(parseConfigFile(wd)).To(Succeed())

		Expect(viper.GetInt(Port)).To(Equal(port))
	})
})
