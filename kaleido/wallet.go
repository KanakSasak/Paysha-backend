package kaleido

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	coremsp "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type Wallet struct {
	UserName string
	Signer   coremsp.SigningIdentity
	network  KaleidoNetwork
	sdk      *fabsdk.FabricSDK
}

func NewWallet(user string, network KaleidoNetwork, sdk *fabsdk.FabricSDK) *Wallet {
	return &Wallet{
		UserName: user,
		network:  network,
		sdk:      sdk,
	}
}

func (w *Wallet) InitIdentity() error {
	ctxProvider := w.sdk.Context()
	mspClient, err := msp.New(ctxProvider, msp.WithCAInstance(w.network.MyCA.MembershipID))
	if err != nil {
		fmt.Printf("Failed to create Fabric CA client. %v\n", err)
		return err
	}
	si, err := mspClient.GetSigningIdentity(w.UserName)
	if err != nil {
		fmt.Printf("Failed to get signing identity. %v. Assuming this is a new user and will attempt to register and enroll.\n", err)
		baseUrl := getApiUrl()
		svcId := w.network.MyCA.ID
		// register for the signing identity
		reg1 := make(map[string]string)
		reg1["enrollmentID"] = w.UserName
		reg1["role"] = "client"
		payload := make(map[string][]interface{})
		payload["registrations"] = []interface{}{reg1}

		result := make(map[string][]map[string]interface{})

		fullUrl := fmt.Sprintf("%s/fabric-ca/%s/register", baseUrl, svcId)
		res, err := w.network.client.Client.R().SetBody(payload).SetResult(&result).Post(fullUrl)
		if err != nil {
			fmt.Printf("Failed to register user with Fabric CA. %v\n", err)
			return err
		}

		fmt.Printf("Register result: %v. Response: %v\n", result, res)
		registrations := result["registrations"]
		si, err = w.enroll(registrations[0], mspClient)
		if err != nil {
			fmt.Printf("Failed to enroll user %s", w.UserName)
			return err
		}
	}
	w.Signer = si
	fmt.Printf("Signing identity: %v\n", si.Identifier())

	return nil
}

func (w *Wallet) enroll(userRegistration map[string]interface{}, mspClient *msp.Client) (coremsp.SigningIdentity, error) {
	secret := userRegistration["enrollmentSecret"].(string)
	err := mspClient.Enroll(w.UserName, msp.WithSecret(secret))
	if err != nil {
		fmt.Printf("Failed to enroll. %v\n", err)
		return nil, err
	}

	si, err := mspClient.GetSigningIdentity(w.UserName)
	if err != nil {
		fmt.Printf("Failed to get signing identity. %v\n", err)
		return nil, err
	}

	return si, nil
}
