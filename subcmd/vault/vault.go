package vault

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jobtalk/pnzr/lib"
	"github.com/ieee0824/getenv"
)

var flagSet = &flag.FlagSet{}

var (
	kmsKeyID       *string
	encryptFlag    *bool
	decryptFlag    *bool
	file           *string
	f              *string
)

func init() {
	kmsKeyID = flagSet.String("key_id", getenv.String("KMS_KEY_ID"), "Amazon KMS key ID")
	encryptFlag = flagSet.Bool("encrypt", getenv.Bool("ENCRYPT", false), "encrypt mode")
	decryptFlag = flagSet.Bool("decrypt", getenv.Bool("DECRYPT", false), "decrypt mode")

	file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")
}

func encrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMS()
	_, err = kms.SetKeyID(keyID).Encrypt(bin)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(kms.String()), 0644)
}

func decrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMSFromBinary(bin)
	if kms == nil {
		return errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.SetKeyID(keyID).Decrypt()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainText, 0644)
}

type Vault struct{}

func (c *Vault) Help() string {
	var msg string
	msg += "usage: pnzr vault [options ...]\n"
	msg += "options:\n"
	msg += "    -key_id\n"
	msg += "        set kms key id\n"
	msg += "    -encrypt\n"
	msg += "        use encrypt mode\n"
	msg += "    -decrypt\n"
	msg += "        use decrypt mode\n"
	msg += "    -file\n"
	msg += "        setting target file\n"
	msg += "    -f"
	msg += "        setting target file\n"
	msg += "    -profile\n"
	msg += "        aws credential name\n"
	msg += "    -region\n"
	msg += "        aws region name\n"
	msg += "    -aws-access-key-id\n"
	msg += "        setting aws access key id\n"
	msg += "    -aws-secret-key-id\n"
	msg += "        setting aws secret key id\n"
	msg += "===================================================\n"
	return msg
}

func (c *Vault) Run(args []string) int {
	if err := flagSet.Parse(args); err != nil {
		log.Fatalln(err)
	}

	if *f == "" && *file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		file = &targetName
	}


	if *file == "" {
		file = f
	}
	if *encryptFlag == *decryptFlag {
		log.Fatalln("Choose whether to execute encrypt or decrypt.")
	}
	if *decryptFlag {
		err := decrypt(*kmsKeyID, *file)
		if err != nil {
			log.Fatalln(err)
		}
	} else if *encryptFlag {
		err := encrypt(*kmsKeyID, *file)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return 0
}

func (c *Vault) Synopsis() string {
	return c.Help()
}
