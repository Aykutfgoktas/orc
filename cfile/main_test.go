package cfile

import (
	"math"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cfile Suite")
}

type Config struct {
	APIKey              string   `json:"key"`
	DefaultOrganization string   `json:"org"`
	Organizations       []string `json:"orgs"`
}

var _ = Describe("Config file service", func() {
	var fileName = "test.json"

	configService := New(fileName)

	Describe("NewConfigFile", func() {
		It("should get the config service", func() {
			Expect(configService).To(Not(BeNil()))
		})
	})

	Describe("ConfigFile", func() {
		It("should return file name with test.json", func() {
			Expect(configService.ConfigFile()).To(Equal(fileName))
		})
	})

	Describe("CheckConfigFile", func() {
		It("should return false when it checks the file", func() {
			Expect(configService.CheckConfigFile()).To(Equal(false))
		})

		It("should return true when it checks the file", func() {
			_ = os.WriteFile(fileName, []byte{}, 0600)

			Expect(configService.CheckConfigFile()).To(Equal(true))

			os.Remove(fileName)
		})
	})

	Describe("ConfigWriter", func() {
		It("should return error while writing to file", func() {

			_, err := configService.Writer(math.Inf(1))

			Expect(err).To(Not(BeNil()))

			os.Remove(fileName)

		})

		It("should return path", func() {

			result, err := configService.Writer("test")

			Expect(err).To(BeNil())

			cur, _ := os.Getwd()

			Expect(result).To(Equal(cur))

			os.Remove(fileName)

		})
	})

	Describe("ConfigReader", func() {
		It("should return error while reading the file", func() {

			_, err := configService.Reader()

			Expect(err).To(Not(BeNil()))

			os.Remove(fileName)

		})

		It("should return error while reading the file", func() {

			conf := Config{
				APIKey:              gofakeit.Word(),
				DefaultOrganization: gofakeit.Word(),
				Organizations:       []string{gofakeit.Word()},
			}

			_, err := configService.Writer(conf)

			Expect(err).To(BeNil())

			result, err := configService.Reader()

			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())

			var newConf Config

			err = result.Decode(&newConf)

			Expect(err).To(BeNil())

			Expect(newConf).To(Equal(conf))

			os.Remove(fileName)
		})

		It("should return error while decoding file", func() {

			conf := Config{
				APIKey:              gofakeit.Word(),
				DefaultOrganization: gofakeit.Word(),
				Organizations:       []string{gofakeit.Word()},
			}

			_, err := configService.Writer(conf)

			Expect(err).To(BeNil())

			result, err := configService.Reader()

			Expect(result).To(Not(BeNil()))
			Expect(err).To(BeNil())

			var newConf Config

			err = result.Decode(newConf)

			Expect(err).To(Not(BeNil()))

			os.Remove(fileName)
		})
	})
})
