package main

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-smtp"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// The Backend implements SMTP server methods.
type Backend struct{}

// Login handles a login command with username and password.
func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {

	return &Session{}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

// A Session is returned after successful login.
type Session struct{}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Incoming mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("To:", to)
	//do your thing here!

	parts := strings.Split(to, "@")
	url := viper.GetString("action_url")
	furl := fmt.Sprintf(url, parts[0])
	fetchurl(furl)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
    	if(viper.GetBool("debug")){
		log.Println("Data:", string(b))
        }
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func main() {

	config()
	
	viper.SetDefault("port", "25")
	viper.SetDefault("usetls", false)
	viper.SetDefault("debug", false)

	be := &Backend{}

	usetls := viper.GetBool("usetls")

	s := smtp.NewServer(be)

	s.Addr = ":" + viper.GetString("port")
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	if viper.GetBool("debug") {
		file, err := os.Create("./tls.log")
		if(err!=nil){
    		
    		log.Fatal(err)
		}
		lg := log.New(file, "", log.LstdFlags|log.Lshortfile)
		s.ErrorLog = lg
		s.Debug = file
	}

	if usetls {

		cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
		if err != nil {
			log.Println(err)
			return
		}

		s.TLSConfig = &tls.Config{
			InsecureSkipVerify:       true,
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS10,
			MaxVersion:               tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_FALLBACK_SCSV,
			},
			Certificates: []tls.Certificate{cer},
		}

	}

	log.Println("Starting server at", s.Addr)
	
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func fetchurl(url string) {

	log.Println("Making request:",url)
	response, err := http.Get(url)
	if err != nil {
        if(viper.GetBool("debug")){
		 log.Printf("%s", err)
        }
		
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
           if(viper.GetBool("debug")){
			log.Printf("%s", err)
           }
			
		}
        if(viper.GetBool("debug")){
		log.Printf("%s\n", string(contents))
		}
		
	}
}

func config() {

	viper.SetConfigName("config")          // name of config file (without extension)
	viper.SetConfigType("json")            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/ipcamemapi/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.ipcamemapi") // call multiple times to add many search paths
	viper.AddConfigPath(".")
	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			log.Println(fmt.Errorf("Fatal error config file: %s \n", err))
		}

	})
	// optionally look for config in the working directory
	err2 := viper.ReadInConfig() // Find and read the config file
	if err2 != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err2))
	}

}
