package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	secret "inspeqtor-stuff/customers"
	"io"
	"io/ioutil"
	"math/rand"
	"net/smtp"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"text/template"
	"time"

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

type EmailData struct {
	License
	RepoToken string
}

func (l *License) HostDescription() string {
	if l.HostLimit < 0 {
		return "unlimited"
	}
	return fmt.Sprintf("up to %d", l.HostLimit)
}

func main() {
	email := ""
	if len(os.Args) > 1 {
		email = os.Args[1]
	}
	var licinfo *secret.Lic

	for _, lic := range secret.List {
		if lic.Email == email {
			licinfo = &lic
			break
		}
	}

	if licinfo == nil && email != "" {
		panic("No such customer: " + email)
	}

	createMasterPasswd()

	if licinfo != nil {
		repoToken := createMasterToken(licinfo.Email)
		user, pwd := repoToken[0:8], repoToken[9:17]

		doc := License{
			licinfo.Org,
			licinfo.Email,
			time.Now().Format("2006-01-02"),
			user,
			pwd,
			licinfo.HostLimit,
			int(rand.Int31()),
		}

		fmt.Printf("Creating license for %s, Hosts %s\n", doc.Email, doc.HostDescription())
		fmt.Printf("%+v\n", doc)

		licdata := encrypt(doc)

		sendEmail(email, EmailData{doc, repoToken}, licdata)
	}
}

func createMasterToken(email string) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s-%s", secret.Salt, email)))
	return hex.EncodeToString(hash.Sum(nil))
}

func createMasterPasswd() {
	file, err := os.Create("inspeqtor-passwd")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, customer := range secret.List {
		fmt.Println(customer.Email)
		repoToken := createMasterToken(customer.Email)
		user, pwd := repoToken[0:8], repoToken[9:17]
		line := toHttpd(user, pwd)
		file.Write([]byte(line))
		file.Write([]byte("\n"))
	}
}

func toHttpd(user, pwd string) string {
	cmd := exec.Command("htpasswd", "-nb", user, pwd)
	data, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(data))
}

func sendEmail(email string, lic EmailData, licdata []byte) {
	data, err := ioutil.ReadFile("cmd/license-email.txt")
	if err != nil {
		panic(err)
	}
	temp := template.Must(template.New("license-email.txt").Parse(string(data)))
	var doc bytes.Buffer
	err = temp.Execute(&doc, &lic)

	body := doc.Bytes()

	content := makeContent(email, body, licdata)

	auth := smtp.PlainAuth("", "mike@contribsys.com", secret.Pwd, "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "mike@contribsys.com",
		[]string{email}, content)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent email to " + email)
}

func makeContent(to string, body []byte, licdata []byte) []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: Contributed Systems <support@contribsys.com>\n")
	buf.WriteString("To: " + to + "\n")
	buf.WriteString("Subject: Thank you for buying Inspeqtor Pro!\n")
	buf.WriteString("MIME-Version: 1.0\n")

	boundary := "f46d043c813270fc6b04c2d223da"

	buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\n")
	buf.WriteString("--" + boundary + "\n")
	buf.WriteString("Content-Type: text/plain; charset=utf-8\n\n")
	buf.Write(body)

	buf.WriteString("\n\n--" + boundary + "\n")

	buf.WriteString("Content-Type: application/octet-stream\n")
	buf.WriteString("Content-Transfer-Encoding: base64\n")
	buf.WriteString("Content-Disposition: attachment; filename=\"license.bin\"\n\n")

	b := make([]byte, base64.StdEncoding.EncodedLen(len(licdata)))
	base64.StdEncoding.Encode(b, licdata)
	buf.Write(b)
	buf.WriteString("\n--" + boundary)
	buf.WriteString("--")

	return buf.Bytes()
}

func encrypt(doc License) []byte {
	myprv := readKey("src/inspeqtor-stuff/keys/local.prv")
	pub := readKey("src/inspeqtor-stuff/keys/remote.pub")

	fmt.Println("Encrypting this document")
	fmt.Println(doc)
	data, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}

	nonce := loadNonce()

	enc := box.Seal(nil, data, nonce, pub, myprv)

	f, err := os.Create("license.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(enc)

	return enc
}

func loadNonce() *[24]byte {
	usr, _ := user.Current()
	dir := usr.HomeDir
	path := fmt.Sprintf("%s/%s", dir, "src/inspeqtor-stuff/keys/nonce.bin")

	nonce := new([24]byte)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.ReadFull(f, nonce[:])
	return nonce
}

func readKey(name string) *[32]byte {
	usr, _ := user.Current()
	dir := usr.HomeDir
	path := fmt.Sprintf("%s/%s", dir, name)

	pub := new([32]byte)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	_, err = io.ReadFull(f, pub[:])
	if err != nil {
		panic(err)
	}
	return pub
}
