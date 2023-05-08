package config

import (
	"encoding/json"
	"errors"

	"github.com/Aykutfgoktas/orc/cfile/mocks"

	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config service", func() {
	var (
		configFileService *mocks.MockIConfigFile
		readerMock        *mocks.MockIReader
		configService     Service
		ctrl              *gomock.Controller
		conf              Config
		org               string
		errMsg            error
	)

	BeforeEach(func() {

		ctrl = gomock.NewController(GinkgoT())
		configFileService = mocks.NewMockIConfigFile(ctrl)
		readerMock = mocks.NewMockIReader(ctrl)
		configService = New(configFileService)
		errMsg = errors.New(gofakeit.Error().Error())
		org = gofakeit.Company()
		organization := gofakeit.Company()
		conf = Config{
			APIKey:              gofakeit.Word(),
			DefaultOrganization: organization,
			Organizations:       []string{organization},
		}

	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("NewConfig", func() {
		It("should get the config service", func() {

			Expect(configService).To(Not(BeNil()))
		})
	})

	Describe("ConfigFile", func() {
		It("should return the config file name", func() {

			configFileService.EXPECT().ConfigFile().Times(1).Return("test")

			result := configService.ConfigFile()

			Expect(result).To(Equal("test"))
		})
	})

	Describe("CheckConfigFile", func() {
		It("should return the true for config file", func() {

			configFileService.EXPECT().CheckConfigFile().Times(1).Return(true)

			result := configService.CheckConfigFile()

			Expect(result).To(Equal(true))
		})
	})

	Describe("Create", func() {

		It("should return the error", func() {

			configFileService.EXPECT().Writer(conf).Times(1).Return("", errors.New("a"))

			result, err := configService.Create(conf.APIKey, conf.DefaultOrganization)

			Expect(result).To(Equal(""))
			Expect(err).To(Not(BeNil()))
		})

		It("should return the path", func() {

			configFileService.EXPECT().Writer(conf).Times(1).Return("path", nil)

			result, err := configService.Create(conf.APIKey, conf.DefaultOrganization)

			Expect(result).To(Equal("path"))
			Expect(err).To(BeNil())
		})
	})

	Describe("Read", func() {

		It("should return the error", func() {

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, errors.New("a"))

			result, err := configService.Read()

			Expect(result).To(BeNil())
			Expect(err).To(Not(BeNil()))
		})

		It("should return the config", func() {
			conff := Config{}
			b, _ := json.Marshal(conf)

			readerMock.EXPECT().Decode(&conff).Times(1).Do(func(d interface{}) error {
				err := json.Unmarshal(b, d)

				if err != nil {
					return err
				}

				return nil
			})

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, nil)

			result, err := configService.Read()

			Expect(result.APIKey).To(Equal(conf.APIKey))
			Expect(result.DefaultOrganization).To(Equal(conf.DefaultOrganization))
			Expect(result.DefaultOrganization).To(Equal(conf.DefaultOrganization))
			Expect(err).To(BeNil())
		})

		It("should return the decode error", func() {
			conff := Config{}

			decodeError := decodeError(errMsg)

			readerMock.EXPECT().Decode(&conff).Times(1).Return(errMsg)

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, nil)

			_, err := configService.Read()

			Expect(err).To(Not(BeNil()))
			Expect(err).To(Equal(decodeError))
		})
	})

	Describe("UpdateDefaultOrganization", func() {

		It("should return the reader error", func() {
			readerError := readerError(errMsg)

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, errMsg)

			err := configService.UpdateDefaultOrganization(org)

			Expect(err).To(Equal(readerError))
		})

		It("should return the decode error", func() {
			conff := Config{}

			decodeError := decodeError(errMsg)

			readerMock.EXPECT().Decode(&conff).Times(1).Return(errMsg)

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, nil)

			err := configService.UpdateDefaultOrganization(org)

			Expect(err).To(Equal(decodeError))
		})

		It("should return the writer error", func() {
			conff := Config{}

			b, _ := json.Marshal(conf)

			writerError := writerError(errMsg)

			readerMock.EXPECT().Decode(&conff).Times(1).Do(func(d interface{}) error {
				err := json.Unmarshal(b, d)

				if err != nil {
					return err
				}

				return nil
			})

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, nil)

			conf.DefaultOrganization = org

			configFileService.EXPECT().Writer(conf).Times(1).Return("", errMsg)

			err := configService.UpdateDefaultOrganization(org)

			Expect(err).To(Equal(writerError))
		})

		It("should return the success", func() {
			conff := Config{}

			b, _ := json.Marshal(conf)

			readerMock.EXPECT().Decode(&conff).Times(1).Do(func(d interface{}) error {
				err := json.Unmarshal(b, d)

				if err != nil {
					return err
				}

				return nil
			})

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, nil)

			conf.DefaultOrganization = org

			configFileService.EXPECT().Writer(conf).Times(1).Return("", nil)

			err := configService.UpdateDefaultOrganization(org)

			Expect(err).To(BeNil())
		})

	})
	Describe("AddOrganization", func() {

		It("should return the reader error", func() {
			readerError := readerError(errMsg)

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, errMsg)

			_, err := configService.AddOrganization(org)

			Expect(err).To(Equal(readerError))
		})

		It("should return the decode error", func() {
			conff := Config{}

			decodeError := decodeError(errMsg)

			readerMock.EXPECT().Decode(&conff).Times(1).Return(errMsg)

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, nil)

			_, err := configService.AddOrganization(org)

			Expect(err).To(Equal(decodeError))
		})

		It("should return the writer error", func() {
			conff := Config{}

			b, _ := json.Marshal(conf)

			writerError := writerError(errMsg)

			readerMock.EXPECT().Decode(&conff).Times(1).Do(func(d interface{}) error {
				err := json.Unmarshal(b, d)

				if err != nil {
					return err
				}

				return nil
			})

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, nil)

			conf.Organizations = append(conf.Organizations, org)

			configFileService.EXPECT().Writer(conf).Times(1).Return("", errMsg)

			_, err := configService.AddOrganization(org)

			Expect(err).To(Equal(writerError))
		})

		It("should return the success with true", func() {
			conff := Config{}

			b, _ := json.Marshal(conf)

			readerMock.EXPECT().Decode(&conff).Times(1).Do(func(d interface{}) error {
				err := json.Unmarshal(b, d)

				if err != nil {
					return err
				}

				return nil
			})

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, nil)

			conf.Organizations = append(conf.Organizations, org)

			configFileService.EXPECT().Writer(conf).Times(1).Return("", nil)

			result, err := configService.AddOrganization(org)

			Expect(err).To(BeNil())
			Expect(result).To(Equal(false))
		})

		It("should return the success with true scenerio existing organization", func() {
			conff := Config{}

			b, _ := json.Marshal(conf)

			readerMock.EXPECT().Decode(&conff).Times(1).Do(func(d interface{}) error {
				err := json.Unmarshal(b, d)

				if err != nil {
					return err
				}

				return nil
			})

			configFileService.EXPECT().Reader().Times(1).Return(readerMock, nil)

			result, err := configService.AddOrganization(conf.DefaultOrganization)

			Expect(err).To(BeNil())
			Expect(result).To(Equal(true))
		})

	})
})
