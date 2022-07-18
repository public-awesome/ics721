package e2e_test

import (
	"encoding/json"
	"testing"

	"fmt"

	wasmibctesting "github.com/CosmWasm/wasmd/x/wasm/ibctesting"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransferTestSuite struct {
	suite.Suite

	coordinator *wasmibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *wasmibctesting.TestChain
	chainB *wasmibctesting.TestChain

	chainAICSAddress sdk.AccAddress
	chainBICSAddress sdk.AccAddress
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = wasmibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(wasmibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(wasmibctesting.GetChainID(1))
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

	chainAStoreResp := suite.chainA.StoreCodeFile("contracts/ics721.wasm")
	require.Equal(suite.T(), uint64(1), chainAStoreResp.CodeID)
	chainBStoreResp := suite.chainB.StoreCodeFile("contracts/ics721.wasm")
	require.Equal(suite.T(), uint64(1), chainBStoreResp.CodeID)

	instantiateICS721 := InstantiateICS721{
		100,
		2,
		"cw721",
	}
	instantiateICS721Raw, err := json.Marshal(instantiateICS721)
	require.NoError(suite.T(), err)
	suite.chainAICSAddress = suite.chainA.InstantiateContract(1, instantiateICS721Raw)
	suite.chainBICSAddress = suite.chainB.InstantiateContract(1, instantiateICS721Raw)
	info := suite.chainA.ContractInfo(suite.chainAICSAddress)
	fmt.Println(suite.chainAICSAddress.String(), suite.chainBICSAddress.String())
	fmt.Println(info.IBCPortID)
}

func (suite *TransferTestSuite) TestICSConnection() {
	var (
		sourcePortID      = suite.chainA.ContractInfo(suite.chainAICSAddress).IBCPortID
		counterpartPortID = suite.chainB.ContractInfo(suite.chainBICSAddress).IBCPortID
	)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
	suite.coordinator.UpdateTime()

	require.Equal(suite.T(), suite.chainA.CurrentHeader.Time, suite.chainB.CurrentHeader.Time)
	path := wasmibctesting.NewPath(suite.chainA, suite.chainB)
	path.EndpointA.ChannelConfig = &ibctesting.ChannelConfig{
		PortID:  sourcePortID,
		Version: "ics721-1",
		Order:   channeltypes.UNORDERED,
	}
	path.EndpointB.ChannelConfig = &ibctesting.ChannelConfig{
		PortID:  counterpartPortID,
		Version: "ics721-1",
		Order:   channeltypes.UNORDERED,
	}

	suite.coordinator.SetupConnections(path)
	suite.coordinator.CreateChannels(path)
}

func TestTransferTest(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}