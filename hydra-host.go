// Copied from hydra-host, since it is not possible to import the main package

package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ory-am/hydra/cli/hydra-host/handler"
	//"github.com/ory-am/hydra/cli/hydra-host/templates"
	"errors"
	"fmt"
	"sync"
	"time"
)

func errorWrapper(f func(ctx *cli.Context) error) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		if err := f(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred: %s\n", err)
		}
	}
}

var (
	ctx      = new(handler.DefaultSystemContext)
	cl       = &handler.Client{Ctx: ctx.GetSystemContext()}
	u        = &handler.Account{Ctx: ctx.GetSystemContext()}
	co       = &handler.Core{Ctx: ctx.GetSystemContext()}
	pl       = &handler.Policy{Ctx: ctx.GetSystemContext()}
	Commands = []cli.Command{
		{
			Name:  "client",
			Usage: "Client actions",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  `Create a new client`,
					Action: errorWrapper(cl.Create),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "i, id",
							Usage: "Set client's id",
						},
						cli.StringFlag{
							Name:  "s, secret",
							Usage: "The client's secret",
						},
						cli.StringFlag{
							Name:  "r, redirect-url",
							Usage: `A list of allowed redirect URLs: https://foobar.com/callback|https://bazbar.com/cb|http://localhost:3000/authcb`,
						},
						cli.BoolFlag{
							Name:  "as-superuser",
							Usage: "Grant superuser privileges to the client",
						},
					},
				},
			},
		},
		{
			Name:  "account",
			Usage: "Account actions",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					Usage:     "Create a new account",
					ArgsUsage: "<username>",
					Action:    errorWrapper(u.Create),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "password",
							Usage: "The account's password",
						},
						cli.BoolFlag{
							Name:  "as-superuser",
							Usage: "Grant superuser privileges to the account",
						},
					},
				},
			},
		},
		{
			Name:   "start",
			Usage:  "Start the host service",
			Action: errorWrapper(start),
		},
		{
			Name:  "jwt",
			Usage: "JWT actions",
			Subcommands: []cli.Command{
				{
					Name:   "generate-keypair",
					Usage:  "Create a JWT PEM keypair.\n\n   You can use these files by providing the environment variables JWT_PRIVATE_KEY_PATH and JWT_PUBLIC_KEY_PATH",
					Action: errorWrapper(handler.CreatePublicPrivatePEMFiles),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "s, private-file-path",
							Value: "rs256-private.pem",
							Usage: "Where to save the private key PEM file",
						},
						cli.StringFlag{
							Name:  "p, public-file-path",
							Value: "rs256-public.pem",
							Usage: "Where to save the private key PEM file",
						},
					},
				},
			},
		},
		{
			Name:  "tls",
			Usage: "TLS actions",
			Subcommands: []cli.Command{
				{
					Name:   "generate-dummy-certificate",
					Usage:  "Create a dummy TLS certificate and private key.\n\n   You can use these files (in development!) by providing the environment variables TLS_CERT_PATH and TLS_KEY_PATH",
					Action: errorWrapper(handler.CreateDummyTLSCert),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "c, certificate-file-path",
							Value: "tls-cert.pem",
							Usage: "Where to save the private key PEM file",
						},
						cli.StringFlag{
							Name:  "k, key-file-path",
							Value: "tls-key.pem",
							Usage: "Where to save the private key PEM file",
						},
						cli.StringFlag{
							Name:  "u, host",
							Usage: "Comma-separated hostnames and IPs to generate a certificate for",
						},
						cli.StringFlag{
							Name:  "sd, start-date",
							Usage: "Creation date formatted as Jan 1 15:04:05 2011",
						},
						cli.DurationFlag{
							Name:  "d, duration",
							Value: 365 * 24 * time.Hour,
							Usage: "Duration that certificate is valid for",
						},
						cli.BoolFlag{
							Name:  "ca",
							Usage: "whether this cert should be its own Certificate Authority",
						},
						cli.IntFlag{
							Name:  "rb, rsa-bits",
							Value: 2048,
							Usage: "Size of RSA key to generate. Ignored if --ecdsa-curve is set",
						},
						cli.StringFlag{
							Name:  "ec, ecdsa-curve",
							Usage: "ECDSA curve to use to generate a key. Valid values are P224, P256, P384, P521",
						},
					},
				},
			},
		},
		{
			Name:  "policy",
			Usage: "Policy actions",
			Subcommands: []cli.Command{
				{
					Name:      "import",
					ArgsUsage: "<policies1.json> <policies2.json> <policies3.json>",
					Usage:     `Import a json file which defines an array of policies`,
					Action:    errorWrapper(pl.Import),
					Flags:     []cli.Flag{},
				},
			},
		},
	}
)

func main() {
	NewApp().Run(os.Args)
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "hydra-host"
	app.Usage = `Dragons guard your resources`
	app.Commands = Commands
	return app
}

/**
* Start and check id superclient is available
*
**/
func start(ctx *cli.Context) error {

	var wg sync.WaitGroup
	wg.Add(1)
	log.Infoln("In Hydra start")

	go func() {
		wg.Wait()
		// CLIENT CREATION that will be used by login App and other apps
		//get environment variables for context
		id := os.Getenv("CLIENT_ID")
		if id == "" {
			log.Errorln(errors.New("Environment variable CLIENT_ID should be set"))
		}
		secret := os.Getenv("CLIENT_SECRET")
		if secret == "" {
			log.Errorln(errors.New("Environment variable CLIENT_SECRET should be set"))
		}
		client, err := cl.Ctx.GetOsins().GetClient(id)
		if err != nil {
			//attempt to create the client
			log.Infoln("Attempt to create the super client because ", err)
			flagSet := flag.NewFlagSet("flagset", *new(flag.ErrorHandling))
			flagSet.String("id", id, "Set client's id")
			flagSet.String("secret", secret, "The client's secret")
			flagSet.String("redirect-url", "http://localhost:8080/authenticate/callback", "The redirect url")
			flagSet.Bool("as-superuser", false, "Grant superuser privileges to the client") // DO NOT GRANT SUPER PRIVILEGES

			ctx = cli.NewContext(nil, flagSet, ctx)
			errCreate := cl.Create(ctx)
			if errCreate != nil {
				log.Errorln("Error on create client : ", errCreate)
			}
		} else {
			log.Infoln("Client found with account id ", client.GetId())
			if client.GetSecret() != secret {
				log.Errorln(errors.New("the secret does not match"))
			}
		}

		// SUPERUSER CREATION
		//get environment variables for context
		superAccountUsername := os.Getenv("SUPERACCOUNT_USERNAME")
		if superAccountUsername == "" {
			log.Errorln(errors.New("Environment variable SUPERACCOUNT_USERNAME should be set"))
		}
		superAccountSecret := os.Getenv("SUPERACCOUNT_SECRET")
		if superAccountSecret == "" {
			log.Errorln(errors.New("Environment variable SUPERACCOUNT_SECRET should be set"))
		}
		account, errAccount := u.Ctx.GetAccounts().GetByUsername(superAccountUsername)
		if errAccount != nil {
			//attempt to create the account
			log.Infoln("Attempt to create the super account because ", errAccount)
			flagSet := flag.NewFlagSet("flagset", *new(flag.ErrorHandling))
			flagSet.String("password", superAccountSecret, "The user's password")

			flagSet.Bool("as-superuser", true, "Grant superuser privileges to the account")

			flagSet.Parse([]string{superAccountUsername})

			ctx = cli.NewContext(nil, flagSet, ctx)
			errCreateAccount := u.Create(ctx)
			if errCreateAccount != nil {
				log.Errorln("Error on create super account : ", errCreateAccount)
			}
		} else {
			if account != nil {
				log.Infof("Super Account found with account %v, id %v", account.GetUsername(), account.GetID())
			} else {
				log.Infoln("No error found when retreiving super account, but account is nil")
			}
		}
	}()
	return co.StartWithInform(ctx, &wg)
}
