package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var wgCommand = &cobra.Command{
	Use:          "wg",
	Short:        "WireGuard utility commands",
	SilenceUsage: true,
}

var genPrivateKeyCommand = &cobra.Command{
	Use:          "genkey",
	Short:        "Generates a new private key and writes it to stdout",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		privateKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return errors.Wrap(err, "failed to generate private key")
		}

		fmt.Println(base64.StdEncoding.EncodeToString(privateKey[:]))
		return nil
	},
}

var genPublicKeyCommand = &cobra.Command{
	Use:          "pubkey",
	Short:        "Reads a private key from stdin and writes a public key to stdout",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		scanner := bufio.NewScanner(cmd.InOrStdin())
		if !scanner.Scan() {
			return errors.New("failed to read private key from stdin")
		}

		privateKeyStr := strings.TrimSpace(scanner.Text())
		privateKey, err := wgtypes.ParseKey(privateKeyStr)
		if err != nil {
			return errors.Wrap(err, "failed to parse private key")
		}

		publicKey := privateKey.PublicKey()

		fmt.Println(base64.StdEncoding.EncodeToString(publicKey[:]))
		return nil
	},
}

func init() {
	wgCommand.AddCommand(genPrivateKeyCommand)
	wgCommand.AddCommand(genPublicKeyCommand)
	rootCmd.AddCommand(wgCommand)
}
