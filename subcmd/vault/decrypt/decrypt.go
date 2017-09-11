package decrypt

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/lib"
	"io/ioutil"
)

type DecryptCommand struct {
	sess           *session.Session
	file           *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
}

func (d *DecryptCommand) parseArgs(args []string) (helpString string) {
	var (
		flagSet = new(flag.FlagSet)
		f       *string
	)

	buffer := new(bytes.Buffer)
	flagSet.SetOutput(buffer)

	d.profile = flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	d.region = flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")
	d.awsAccessKeyID = flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	d.awsSecretKeyID = flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")
	d.file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return buffer.String()
		}
		panic(err)
	}

	if *f == "" && *d.file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		d.file = &targetName
	}

	if *d.file == "" {
		d.file = f
	}

	if *d.awsAccessKeyID != "" && *d.awsSecretKeyID != "" && *d.profile == "" {
		d.sess = session.Must(session.NewSessionWithOptions(session.Options{
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			SharedConfigState:       session.SharedConfigEnable,
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(*d.awsAccessKeyID, *d.awsSecretKeyID, ""),
				Region:      d.region,
			},
		}))
	} else {
		d.sess = session.Must(session.NewSessionWithOptions(session.Options{
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			SharedConfigState:       session.SharedConfigEnable,
			Profile:                 *d.profile,
		}))

		if d.region != nil && *d.region != "" {
			d.sess.Config.Region = d.region
		}
	}

	return
}

func (d *DecryptCommand) decrypt(fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMSFromBinary(bin, d.sess)
	if kms == nil {
		return errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.Decrypt()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainText, 0644)
}

func (d *DecryptCommand) Help() string {
	return d.parseArgs([]string{"-h"})
}

func (d *DecryptCommand) Synopsis() string {
	return "Decryption mode of encrypted file."
}

func (d *DecryptCommand) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println(d.Synopsis())
		return 0
	}
	d.parseArgs(args)
	if err := d.decrypt(*d.file); err != nil {
		panic(err)
	}

	return 0
}
