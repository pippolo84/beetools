package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	encodeCmd := &cobra.Command{
		Use:   "encode",
		Short: "Encode data from JSON",
		Long:  "Encode data from JSON to bencode format.",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				r io.Reader
				w io.Writer
			)

			r = os.Stdin
			if len(args) > 0 {
				in, err := os.Open(args[0])
				if err != nil {
					return err
				}
				defer in.Close()

				r = in
			}

			w = os.Stdout
			if len(args) > 1 {
				out, err := os.Create(args[1])
				if err != nil {
					return err
				}
				defer out.Close()

				w = out
			}

			if err := encode(w, r); err != nil {
				fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
			}
			return nil
		},
	}

	decodeCmd := &cobra.Command{
		Use:   "decode",
		Short: "Decode data from bencode",
		Long:  "Decode data from bencode format to JSON.",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				r io.Reader
				w io.Writer
			)

			r = os.Stdin
			if len(args) > 0 {
				in, err := os.Open(args[0])
				if err != nil {
					return err
				}
				defer in.Close()

				r = in
			}

			w = os.Stdout
			if len(args) > 1 {
				out, err := os.Create(args[1])
				if err != nil {
					return err
				}
				defer out.Close()

				w = out
			}

			if err := decode(w, r); err != nil {
				fmt.Fprintf(os.Stderr, "decode error: %v\n", err)
			}
			return nil
		},
	}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "show data from bencode-encoded data",
		Long:  "Show relevant info from bencoded-encoded data.",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var r io.Reader

			r = os.Stdin
			if len(args) > 0 {
				in, err := os.Open(args[0])
				if err != nil {
					return err
				}
				defer in.Close()

				r = in
			}

			if err := show(os.Stdout, r); err != nil {
				fmt.Fprintf(os.Stderr, "show error: %v\n", err)
			}
			return nil
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
