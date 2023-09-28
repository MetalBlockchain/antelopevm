package cmd

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createKeyCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}

var createKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Create a new keypair and print the public and private keys",
	Long:  `Create a new keypair and print the public and private keys`,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := ecc.NewRandomPrivateKey()
		if err != nil {
			fmt.Printf("failed to generate keypair: %v\n", err)
		}

		fmt.Printf("Private key: %s\n", key.String())
		fmt.Printf("Public key: %s\n", key.PublicKey())
	},
}
