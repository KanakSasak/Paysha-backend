package kaleido

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type Channel struct {
	ChannelID string
	client    *channel.Client
	sdk       *fabsdk.FabricSDK
}

func NewChannel(channelId string, sdk *fabsdk.FabricSDK) *Channel {
	return &Channel{
		ChannelID: channelId,
		sdk:       sdk,
	}
}

func (c *Channel) Connect(signer *msp.IdentityIdentifier) error {
	channelContext := c.sdk.ChannelContext(c.ChannelID, fabsdk.WithUser(signer.ID), fabsdk.WithOrg(signer.MSPID))
	channelClient, err := channel.New(channelContext)
	if err != nil {
		return fmt.Errorf("Failed to create channel client. %s", err)
	}
	c.client = channelClient
	return nil
}

func (c *Channel) InitChaincode(chaincodeId string) error {
	return nil
}

func (c *Channel) ExecChaincode(chaincodeId, FcName string, idwallettujuan string, amount string) (string, error) {
	return c.invokeChaincode(chaincodeId, FcName, idwallettujuan, amount)
}

func (c *Channel) invokeChaincode(chaincodeId string, FcName string, idwallettujuan string, amount string) (string, error) {

	var response string

	if FcName == "account" {
		log.Println("Invoke ClientAccountID")
		resp, err := c.client.Execute(
			channel.Request{ChaincodeID: chaincodeId, Fcn: "ClientAccountID"},
			channel.WithRetry(retry.DefaultChannelOpts),
		)

		fmt.Printf("wallet: %s", resp.Payload)

		if err != nil {
			return "", fmt.Errorf("Failed to send transaction to invoke the chaincode. %s", err)
		} else {
			response = string(resp.Payload)
		}
	}

	if FcName == "mint" {
		log.Println("Invoke Mint")
		resp, err := c.client.Execute(
			channel.Request{ChaincodeID: chaincodeId, Fcn: "Mint", Args: [][]byte{[]byte(amount)}},
			channel.WithRetry(retry.DefaultChannelOpts),
		)
		log.Println(resp)
		if err != nil {
			return "", fmt.Errorf("Failed to send transaction to invoke the chaincode. %s", err)
		} else {
			response = string(resp.Payload)
		}
	}

	if FcName == "balanceof" {
		log.Println("Invoke BalanceOf")
		resp, err := c.client.Execute(
			channel.Request{ChaincodeID: chaincodeId, Fcn: "BalanceOf", Args: [][]byte{[]byte("")}},
			channel.WithRetry(retry.DefaultChannelOpts),
		)
		//for i, _ := range resp.Responses {
		log.Println(string(resp.Payload))

		//}
		if err != nil {
			return "", fmt.Errorf("Failed to send transaction to invoke the chaincode. %s", err)
		} else {
			response = string(resp.Payload)
		}
	}

	if FcName == "clientaccountbalance" {
		log.Println("Invoke ClientAccountBalance")
		resp, err := c.client.Execute(
			channel.Request{ChaincodeID: chaincodeId, Fcn: "ClientAccountBalance"},
			channel.WithRetry(retry.DefaultChannelOpts),
		)
		//for i, _ := range resp.Responses {
		log.Println(string(resp.Payload))

		//}
		if err != nil {
			return "", fmt.Errorf("Failed to send transaction to invoke the chaincode. %s", err)
		} else {
			response = string(resp.Payload)
		}
	}

	if FcName == "total" {
		log.Println("Invoke TotalSupply")
		resp, err := c.client.Execute(
			channel.Request{ChaincodeID: chaincodeId, Fcn: "TotalSupply"},
			channel.WithRetry(retry.DefaultChannelOpts),
		)
		//for i, _ := range resp.Responses {
		log.Println(string(resp.Payload))

		//}
		if err != nil {
			return "", fmt.Errorf("Failed to send transaction to invoke the chaincode. %s", err)
		} else {
			response = string(resp.Payload)
		}
	}

	if FcName == "transfer" {
		log.Println("Invoke Transfer")
		resp, err := c.client.Execute(
			channel.Request{ChaincodeID: chaincodeId, Fcn: "Transfer", Args: [][]byte{[]byte(idwallettujuan), []byte(amount)}},
			channel.WithRetry(retry.DefaultChannelOpts),
		)
		//for i, _ := range resp.Responses {
		log.Println(string(resp.Payload))

		//}
		if err != nil {
			return "", fmt.Errorf("Failed to send transaction to invoke the chaincode. %s", err)
		} else {
			response = string(resp.Payload)
		}
	}

	return response, nil
}
