package database

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var (
	clientDBfirebase   *db.Client
	clientAuthfirebase *auth.Client
)

func FirebaseConnect() {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://transpay.firebaseio.com/",
	}
	// Fetch the service account key JSON file contents
	pathX, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	dir := path.Dir(pathX)

	os := runtime.GOOS
	switch os {
	case "darwin":
		fmt.Println("MAC operating system")
		dir = "/Users/laluraynaldi/TestingApps/hyperledger-fabric/PaySha/kaleido-network/golang"
	case "windows":

	default:
		fmt.Printf("%s.\n", os)
	}

	exPath := filepath.FromSlash(dir + "/service_account.json")
	abs, err := filepath.Abs(exPath)
	if err == nil {
		fmt.Println("Absolute:", abs)
	}
	opt := option.WithCredentialsFile(abs)
	//opt := option.WithCredentialsFile("/home/laluraynaldi/gotrans/gotrans-project-firebase-adminsdk.json")

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		panic(err)
	}

	clientDB, err := app.Database(ctx)
	if err != nil {
		log.Println("Error initializing database client:", err)
		panic(err)
	} else {
		log.Println("firebase initializing successfull")
	}

	clientAuth, err := app.Auth(ctx)
	if err != nil {
		log.Println("Error initializing database client:", err)
		panic(err)
	} else {
		log.Println("firebase initializing successfull")
	}

	SetUpfirebasedb(clientDB)
	SetUpfirebaseAuth(clientAuth)

}

func SetUpfirebasedb(DB *db.Client) {
	clientDBfirebase = DB
}

func SetUpfirebaseAuth(auth *auth.Client) {
	clientAuthfirebase = auth
}

func Getfirebasedb() *db.Client {
	return clientDBfirebase
}

func GetfirebaseAuth() *auth.Client {
	return clientAuthfirebase
}
