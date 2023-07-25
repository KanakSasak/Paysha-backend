package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func GenerateKeyFiles(id string) error {
	// Create the keys
	priv, pub := GenerateRsaKeyPair()

	// Export the keys to pem string
	priv_pem := ExportRsaPrivateKeyAsPemStr(priv)
	pub_pem, _ := ExportRsaPublicKeyAsPemStr(pub)

	// Import the keys from pem string
	priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(priv_pem)
	pub_parsed, _ := ParseRsaPublicKeyFromPemStr(pub_pem)

	// Export the newly imported keys
	priv_parsed_pem := ExportRsaPrivateKeyAsPemStr(priv_parsed)
	pub_parsed_pem, _ := ExportRsaPublicKeyAsPemStr(pub_parsed)

	fmt.Println(priv_parsed_pem)
	fmt.Println(pub_parsed_pem)

	pathX, err := os.Executable()
	if err != nil {
		log.Println(err)
		//lognew.Error(err)
		return err

	}
	dir := path.Dir(pathX)

	var osflag = ""
	opSys := runtime.GOOS
	switch opSys {
	case "darwin":
		fmt.Println("MAC operating system")
		osflag = "1"
	case "windows":
		fmt.Println("windows operating system")
		osflag = "2"
	default:
		fmt.Printf("%s.\n", opSys)
	}

	if osflag != "1" && osflag != "2" {
		folderInfo, err := os.Stat("/usr/src")
		if os.IsNotExist(err) {
			err := os.Mkdir("/usr/src", 7777)
			if err != nil {
				log.Fatal(err)
			}

		}
		dir = "/usr/src"
		log.Println(folderInfo)
	} else {
		if osflag == "1" {
			dir = "/Users/laluraynaldi/TestingApps/jwtrsakey"
		} else if osflag == "2" {
			dir = "C:\\Users\\lalu.putra\\Documents\\JWT keys app"
		} else {
			dir = "notfound"
		}

	}
	exPathpriv := filepath.FromSlash(dir + "/" + id + "_jwk_key.key")
	abspriv, err := filepath.Abs(exPathpriv)
	if err != nil {

		return err
	}

	exPathpub := filepath.FromSlash(dir + "/" + id + "_jwk_key.pub")
	abspub, err := filepath.Abs(exPathpub)
	if err != nil {
		//fmt.Println("Absolute:", abspub)
		return err

	}

	log.Println(exPathpriv)

	err1 := writeNewFile(abspub, []byte(pub_parsed_pem), 0444)
	if err1 != nil {
		//fmt.Println("Absolute:", abspub)
		return err

	}
	err2 := writeNewFile(abspriv, []byte(priv_parsed_pem), 0400)
	if err2 != nil {
		//fmt.Println("Absolute:", abspub)
		return err

	}

	// Check that the exported/imported keys match the original keys
	if priv_pem != priv_parsed_pem && pub_pem != pub_parsed_pem {
		fmt.Println("Failure: Export and Import did not result in same Keys")
		return err
	} else {
		fmt.Println("Success")
		return nil
	}
}

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return string(pubkey_pem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

// writeNewFile is shameless copy-paste from ioutil.WriteFile with a bit
// different flags for OpenFile.
func writeNewFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
