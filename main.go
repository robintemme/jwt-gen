package main

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type tokenStruct struct {
	Subject    string
	Expiration time.Time
	Token      string
}

func generateToken(sub string, exp time.Time, privateKey *rsa.PrivateKey) (string, error) {
	// create the token
	token := jwt.New(jwt.SigningMethodRS256)

	// set some claims
	token.Claims["exp"] = exp.Unix()
	token.Claims["sub"] = sub

	// sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func getPrivateKey(file string) (*rsa.PrivateKey, error) {
	var key *rsa.PrivateKey
	privateKeyFile, err := os.Open(file)
	if err != nil {
		return key, err
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	key, err = x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		return key, err
	}

	return key, err
}

func main() {
	// flag parsing
	sub := flag.String("sub", "nobody", "The subject of the token (owner)")
	exp := flag.String("exp", "01.01.1970", "Expiration date (German format)")
	key := flag.String("key", "key.pem", "RSA private key file")
	flag.Parse()

	privateKey, err := getPrivateKey(*key)
	if err != nil {
		panic(err)
	}

	// parse params into struct
	var token tokenStruct
	token.Subject = *sub

	// FIXME hardcoded time zone and format
	loc, _ := time.LoadLocation("Europe/Berlin")
	const germanDateOnlyFormat = "02.01.2006"
	token.Expiration, _ = time.ParseInLocation(germanDateOnlyFormat, *exp, loc)

	t, err := generateToken(token.Subject, token.Expiration, privateKey)
	if err != nil {
		panic(err)
	}
	token.Token = t

	// print out the token
	fmt.Println(token.Token)
}
