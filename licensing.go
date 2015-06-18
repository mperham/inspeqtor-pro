package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/mperham/inspeqtor/util"
	"golang.org/x/crypto/nacl/box"
)

type License struct {
	Org       string
	Email     string
	Date      string
	User      string
	Pwd       string
	HostLimit int
	Nonce     int
}

func (lic *License) Signature(hostname string, now int64) string {
	nonce := lic.Nonce
	str := fmt.Sprintf("%s %d %d", hostname, now, nonce)

	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func verifyLicense(path string) (*License, error) {
	result, err := util.FileExists(path + "/license.bin")
	if err != nil {
		return nil, err
	}

	pub := [32]byte{
		0x2c, 0x6b, 0xe8, 0x2e, 0x1f, 0x62, 0x4a, 0x43, 0xe6, 0xd1, 0x69, 0xbc,
		0x11, 0x4c, 0x05, 0x6f, 0xa4, 0xcd, 0x34, 0x8f, 0x9d, 0xdf, 0x07, 0x59,
		0x82, 0x3f, 0x8a, 0x50, 0xcd, 0x20, 0x1e, 0x13,
	}
	prv := [32]byte{
		0x97, 0x95, 0xd6, 0xf1, 0xce, 0xf3, 0xb8, 0xf9, 0x05, 0x47, 0x6b, 0xb0,
		0x69, 0x39, 0xad, 0xa5, 0xb3, 0x43, 0xcf, 0xf1, 0x8d, 0x32, 0xeb, 0x04,
		0xc6, 0x6e, 0x64, 0x73, 0x62, 0xd8, 0x3f, 0x1b,
	}
	nonce := [24]byte{
		0x5e, 0xf8, 0x3d, 0xad, 0x43, 0x54, 0x2a, 0xe7, 0x8e, 0x13, 0x6c, 0xd7,
		0x84, 0xd6, 0xc9, 0x61, 0xd6, 0x69, 0xc3, 0xcd, 0x1f, 0xd7, 0x8e, 0xbd,
	}

	if !result {
		fmt.Println("Unlicensed, non-production use only.")
		return &License{}, nil
	}

	enc, err := ioutil.ReadFile(path + "/license.bin")
	if err != nil {
		return nil, err
	}

	dec, ok := box.Open(nil, enc, &nonce, &pub, &prv)
	if !ok {
		return nil, errors.New("Invalid license file")
	}

	newdoc := &License{}
	err = json.Unmarshal(dec, newdoc)
	if err != nil {
		return nil, errors.New("Corrupt license file?")
	}

	if newdoc.HostLimit > 0 {
		fmt.Printf("Licensed to %s for up to %d hosts\n", newdoc.Org, newdoc.HostLimit)
	} else {
		fmt.Printf("Licensed to %s, unlimited hosts\n", newdoc.Org)
	}
	fmt.Println("")

	return newdoc, nil
}
