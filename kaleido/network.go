package kaleido

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	kld "github.com/kaleido-io/kaleido-sdk-go/kaleido"
)

type KaleidoNetwork struct {
	Consortium   kld.Consortium
	Environment  kld.Environment
	Memberships  []kld.Membership
	MyMembership kld.Membership
	MyCA         kld.Service
	Orderers     []kld.Node
	Peers        []kld.Node

	client kld.KaleidoClient
}

func NewNetwork() *KaleidoNetwork {
	kn := KaleidoNetwork{}
	return &kn
}

func (kn *KaleidoNetwork) Initialize() {
	fmt.Println("Initializing Kaleido Network...")

	apikey := os.Getenv("APIKEY")
	if apikey == "" {
		apikey = "u0j2w38ebe-JQBUeDqJlJmSNW/4q+ebNoth24Y6a1kMMMt64ya8qyY="
	}

	url := getApiUrl()
	fmt.Printf("URL: %v\n", url)

	kn.client = kld.NewClient(url, apikey)

	kn.selectConsortium()
	kn.selectEnvironment()
	kn.selectMembership()
	kn.getMyCA()
	kn.getOrderersAndPeers()
}

func (kn *KaleidoNetwork) selectConsortium() {
	var consortiums []kld.Consortium
	var targetCon kld.Consortium

	data, err := kn.client.ListConsortium(&consortiums)
	log.Println(data)
	if err != nil {
		fmt.Printf("Failed to get list of consortiums. %v\n", err)
		os.Exit(1)
	}
	liveConsortiums := []kld.Consortium{}
	for _, con := range consortiums {
		if con.State != "deleted" {
			liveConsortiums = append(liveConsortiums, con)
		}
	}

	if len(liveConsortiums) > 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Select the target consortium:")
		for i, con := range liveConsortiums {
			fmt.Printf("\t[%v] %v (%v)\n", i, con.Name, con.State)
		}

		for {
			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)

			i, _ := strconv.Atoi(text)
			targetCon = liveConsortiums[i]
			break
		}
	} else {
		log.Println(liveConsortiums)
		targetCon = liveConsortiums[0]
	}

	fmt.Printf("Target consortium: %v (id=%v)\n", targetCon.Name, targetCon.ID)
	kn.Consortium = targetCon
}

func (kn *KaleidoNetwork) selectEnvironment() {
	var envs []kld.Environment
	var targetEnv kld.Environment

	_, err := kn.client.ListEnvironments(kn.Consortium.ID, &envs)
	if err != nil {
		fmt.Printf("Failed to get list of environments. %v\n", err)
		os.Exit(1)
	}
	liveEnvs := []kld.Environment{}
	for _, env := range envs {
		if env.State != "deleted" {
			liveEnvs = append(liveEnvs, env)
		}
	}

	if len(liveEnvs) > 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Select the target environment:")
		for i, env := range liveEnvs {
			fmt.Printf("\t[%v] %v (%v)\n", i, env.Name, env.State)
		}

		for {
			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)

			i, _ := strconv.Atoi(text)
			targetEnv = liveEnvs[i]
			break
		}
	} else {
		targetEnv = liveEnvs[0]
	}

	fmt.Printf("Target environment: %v (id=%v)\n", targetEnv.Name, targetEnv.ID)
	kn.Environment = targetEnv
}

func (kn *KaleidoNetwork) selectMembership() {
	var memberships []kld.Membership
	var targetMembership kld.Membership

	_, err := kn.client.ListMemberships(kn.Consortium.ID, &memberships)
	if err != nil {
		fmt.Printf("Failed to get list of memberships. %v\n", err)
		os.Exit(1)
	}

	desiredMembership := os.Getenv("SUBMITTER")
	if len(memberships) > 1 && desiredMembership == "" {
		//reader := bufio.NewReader(os.Stdin)
		fmt.Println("Select the membership to submit transactions from:")
		for i, mem := range memberships {
			fmt.Printf("\t[%v] %v (%v)\n", i, mem.OrgName, mem.ID)
		}

		for {
			//fmt.Print("-> ")
			//text, _ := reader.ReadString('\n')
			//// convert CRLF to LF
			//text = strings.Replace(text, "\n", "", -1)
			//
			//i, _ := strconv.Atoi(text)
			//targetMembership = memberships[i]
			fmt.Print("-> 0 \n")
			targetMembership = memberships[0]
			break
		}
	} else {
		if desiredMembership != "" {
			for _, mem := range memberships {
				if mem.ID == desiredMembership {
					targetMembership = mem
					break
				}
			}
		} else {
			targetMembership = memberships[0]
		}
	}

	fmt.Printf("Target membership: %v (id=%v)\n", targetMembership.OrgName, targetMembership.ID)
	kn.MyMembership = targetMembership
	kn.Memberships = memberships
}

func (kn *KaleidoNetwork) getOrderersAndPeers() {
	var nodes []kld.Node

	peers := []kld.Node{}
	orderers := []kld.Node{}

	_, err := kn.client.ListNodes(kn.Consortium.ID, kn.Environment.ID, &nodes)
	if err != nil {
		fmt.Printf("Failed to get list of nodes. %v\n", err)
		os.Exit(1)
	}
	if len(nodes) == 0 {
		fmt.Println("The environment does not have any orderers or peers.")
		os.Exit(1)
	}

	for i := range nodes {
		if nodes[i].Role == "orderer" {
			fmt.Printf("Found orderer %s (membership=%s)\n", nodes[i].ID, nodes[i].MembershipID)
			orderers = append(orderers, nodes[i])
		} else if nodes[i].Role == "peer" {
			fmt.Printf("Found peer %s (membership=%s)\n", nodes[i].ID, nodes[i].MembershipID)
			peers = append(peers, nodes[i])
		}
	}

	if len(orderers) == 0 {
		fmt.Println("The environment does not have any orderers")
		os.Exit(1)
	}
	if len(peers) == 0 {
		fmt.Println("The environment does not have any peers")
		os.Exit(1)
	}

	kn.Orderers = orderers
	kn.Peers = peers
}

func (kn *KaleidoNetwork) getMyCA() {
	var services []kld.Service

	_, err := kn.client.ListServices(kn.Consortium.ID, kn.Environment.ID, &services)
	if err != nil {
		fmt.Printf("Failed to get list of services. %v\n", err)
		os.Exit(1)
	}
	if len(services) == 0 {
		fmt.Println("The environment does not have any services.")
		os.Exit(1)
	}

	for i := range services {
		if services[i].Service == "fabric-ca" && services[i].MembershipID == kn.MyMembership.ID {
			fmt.Printf("Found fabric-ca service %s for my membership %s\n", services[i].ID, kn.MyMembership.ID)
			kn.MyCA = services[i]
			return
		}
	}
	fmt.Printf("Failed to locate fabric-ca service for my membership %s\n", kn.MyMembership.ID)
	os.Exit(1)
}

func getApiUrl() string {
	url := os.Getenv("KALEIDO_URL")
	if url == "" {
		url = "https://console-ap.kaleido.io/api/v1"
	}
	return url
}
