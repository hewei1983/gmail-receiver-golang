package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"flag"
	"github.com/Luzifer/go-openssl"
)

func main() {

	/*flag config*/
	credPtr := flag.String("credential", "foo", "format .json:  e.g., client_secret.json. Created by user for a specific gmail address.")
	periodFromPtr := flag.String("from", "foo", "format date query: e.g., \"2017/06/01\"")
	periodToPtr := flag.String("to", "foo", "format date query: e.g., \"2018/03/01\"")
	pwPtr := flag.String("enckey", "foo", "format: an ASII string. Must be a key different from the default \"foo\"")
	encPtr := flag.Bool("enc", false, "format: use -encrypt to encrypt before saving; otherwise miss this flag to save unencrypted email to local file.")
	flag.Parse()

	/*Google gmail access credenital*/
	ctx := context.Background()
	credential, err := ioutil.ReadFile(*credPtr)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	//must provision a unique credential file created to each specific gmail account.
	config, err := google.ConfigFromJSON(credential, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	/*fetching email*/
	user := "me"
	r, err := srv.Users.Messages.List(user).LabelIds("INBOX").Q("after:" + *periodFromPtr + " " + "before:" + *periodToPtr).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages matching labels and appointed epoch. %v", err)
	}

	log.Printf("Number of emails (messages) successfully retrieved is: %v\n", len(r.Messages))

	var msg_plain string
	for _, m := range r.Messages {
		msg, err := srv.Users.Messages.Get(user, m.Id).Do()
		if err != nil {
			log.Fatalf("unable to retrieve message %v: %v", m.Id, err)
		}
		//formatting email content
		msg_plain = msg_plain + fmt.Sprintf("msg status: %v => msg snippet is: %q\n", msg.LabelIds[0], msg.Snippet)
	}

	/*retrieved email*/
	if *pwPtr == "foo" {
		log.Fatalf("Please input your encryption key!")
	}

	f, err := os.Create("./retrieved_email.dat")
	if err != nil {
		log.Fatalf("Error in creating file %s", err)
	}

	if *encPtr == true {

		// instantiate 3rd party lib to encrypt
		o := openssl.New()
		msg_enc, err := o.EncryptString(*pwPtr, msg_plain)
		if err != nil {
			log.Fatalf("Error in encrypting msg %s", err)
		}
		// save encrypted email to .dat file
		_, err = f.WriteString(string(msg_enc))
		if err != nil {
			log.Fatalf("Error in writing to local file %s", err)
		}
	} else {
		// save plain email to .dat file
		_, err = f.WriteString(string(msg_plain))
		if err != nil {
			log.Fatalf("Error in writing to local file %s", err)
		}
	}

	fmt.Printf("Retrieved emails have been saved to retrieved_email.dat.\nEncryption: %v\n\n", *encPtr)
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// request a Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gmail-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
