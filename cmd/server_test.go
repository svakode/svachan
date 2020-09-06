package cmd

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/svakode/svachan/dictionary"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/mocks"
)

type ServerTestSuite struct {
	suite.Suite
	ctx        *framework.Context
	server     *framework.Server
	memoryInfo *framework.Memory
	cpuInfo    *framework.CPU
	diskInfo   *framework.Disk

	discord *mocks.DiscordService
	memory  *mocks.MemoryService
	cpu     *mocks.CPUService
	disk    *mocks.DiskService
}

func (suite *ServerTestSuite) SetupTest() {
	suite.memoryInfo = &framework.Memory{
		Used:       1.0,
		Total:      1.0,
		Percentage: 100,
	}
	suite.cpuInfo = &framework.CPU{
		Used: 100,
	}
	suite.diskInfo = &framework.Disk{
		Free:  1.0,
		Total: 1.0,
		Used:  100,
	}

	suite.discord = new(mocks.DiscordService)
	suite.memory = new(mocks.MemoryService)
	suite.cpu = new(mocks.CPUService)
	suite.disk = new(mocks.DiskService)

	suite.server = framework.NewServer(suite.memory, suite.cpu, suite.disk)
	suite.ctx = &framework.Context{
		Discord: suite.discord,
		Server:  suite.server,
	}
}

func (suite *ServerTestSuite) TestServerShouldSuccess() {
	suite.memory.On("Get").Return(suite.memoryInfo, nil)
	suite.cpu.On("Get").Return(suite.cpuInfo, nil)
	suite.disk.On("Get").Return(suite.diskInfo, nil)
	suite.discord.On("GetEmbedLayout").Return(&discordgo.MessageEmbed{})
	suite.discord.On("Embed", mock.Anything).Return(&discordgo.Message{})

	ServerCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Embed", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *ServerTestSuite) TestServerShouldReturnErrorWhenGetMemoryFailed() {
	suite.memory.On("Get").Return(nil, errors.New(dictionary.GeneralError))
	suite.discord.On("Reply", dictionary.GeneralError).Return(&discordgo.Message{})

	ServerCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.GeneralError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *ServerTestSuite) TestServerShouldReturnErrorWhenGetCPUFailed() {
	suite.memory.On("Get").Return(suite.memoryInfo, nil)
	suite.cpu.On("Get").Return(nil, errors.New(dictionary.GeneralError))
	suite.discord.On("Reply", dictionary.GeneralError).Return(&discordgo.Message{})

	ServerCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.GeneralError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *ServerTestSuite) TestServerShouldReturnErrorWhenGetDiskFailed() {
	suite.memory.On("Get").Return(suite.memoryInfo, nil)
	suite.cpu.On("Get").Return(suite.cpuInfo, nil)
	suite.disk.On("Get").Return(nil, errors.New(dictionary.GeneralError))
	suite.discord.On("Reply", dictionary.GeneralError).Return(&discordgo.Message{})

	ServerCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.GeneralError)
	suite.discord.AssertExpectations(suite.T())
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
