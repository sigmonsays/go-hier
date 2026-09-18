package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/sigmonsays/go-hier"
	"github.com/spf13/cobra"
)

func init() {
	lookupCmd.Flags().StringP("cfg", "c", "", "deployment configuration")
	lookupCmd.Flags().BoolP("json", "j", false, "json output")
	lookupCmd.Flags().BoolP("merge", "m", false, "merge result")
	// overrides
	lookupCmd.Flags().StringP("fqdn", "f", "", "override fqdn")
	lookupCmd.Flags().StringP("site", "s", "", "override site")
	lookupCmd.Flags().StringP("domain", "D", "", "override domain")
	lookupCmd.Flags().StringP("env", "e", "", "override env")
	rootCmd.AddCommand(lookupCmd)
}

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "lookup variable",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := NewContext(cmd, args)
		defer ctx.Close()
		cfgfile, _ := cmd.Flags().GetString("cfg")
		outputJson, _ := cmd.Flags().GetBool("json")
		mergeResult, _ := cmd.Flags().GetBool("merge")
		hostOverride, _ := cmd.Flags().GetString("host")
		siteOverride, _ := cmd.Flags().GetString("site")
		domainOverride, _ := cmd.Flags().GetString("domain")
		envOverride, _ := cmd.Flags().GetString("env")

		// load spec
		slog.Debug("Load config file", "f", cfgfile)
		spec, err := hier.LoadHierSpec(cfgfile)
		if err != nil {
			return err
		}
		spec.Merge = mergeResult
		if hostOverride != "" {
			spec.Inputs["fqdn"] = hostOverride
		}
		if siteOverride != "" {
			spec.Inputs["site"] = siteOverride
		}
		if domainOverride != "" {
			spec.Inputs["domain"] = domainOverride
		}
		if envOverride != "" {
			spec.Inputs["env"] = envOverride
		}

		wd, _ := os.Getwd()
		spec.FS = os.DirFS(wd)

		for _, arg := range args {
			v, err := hier.Lookup(spec, arg)
			if err != nil {
				if err == hier.LookupNotFound {
					slog.Debug("lookup not found", "arg", arg)
					continue
				}
				slog.Warn("Lookup error", "arg", arg, "err", err)
				continue
			}

			slog.Debug("Value found", "from", v.Path, "value", v.Value)
			if outputJson {
				buf, _ := json.MarshalIndent(v.Value, "", " ")
				fmt.Printf("%s\n", buf)
			} else {
				fmt.Println(v.Value)
			}
		}

		return nil
	},
}
