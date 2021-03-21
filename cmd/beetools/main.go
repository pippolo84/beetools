package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	encodeCmd := &cobra.Command{
		Use:   "encode",
		Short: "Encode data from JSON",
		Long:  "Encode data from JSON to bencode format.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			if err := encode(os.Stdout, os.Stdin); err != nil {
				fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
			}
		},
	}

	decodeCmd := &cobra.Command{
		Use:   "decode",
		Short: "Decode data from bencode",
		Long:  "Decode data from bencode format to JSON.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			if err := decode(os.Stdout, os.Stdin); err != nil {
				fmt.Fprintf(os.Stderr, "decode error: %v\n", err)
			}
		},
	}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "show data from bencode-encoded data",
		Long:  "Show relevant info from bencoded-encoded data.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			if err := show(os.Stdout, os.Stdin); err != nil {
				fmt.Fprintf(os.Stderr, "show error: %v\n", err)
			}
		},
	}

	rootCmd := &cobra.Command{
		Use:   "beetools",
		Short: "beetools is a set of tools to manage bencode format",
		Long: `
beetols is a sample CLI application, able to encode data to and decode data
from bencode format`,
	}
	rootCmd.AddCommand(encodeCmd)
	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(showCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
