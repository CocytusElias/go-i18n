package main

import (
	"errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/language"
	"os"
	"path"
	"regexp"
	"strings"
)

var genCmd = &cli.Command{
	Name:  "gen",
	Usage: "goI18n gen --bundle=[bundlePath] --sc=[statusCode fileName] --pkgName=[final generated file package name] --outputDir=[final generated directory name]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "bundleDir",
			Value: "bundle",
			Usage: "Path to the toml language file",
		},
		&cli.StringFlag{
			Name:  "scFileName",
			Value: "statusCode",
			Usage: "The file name of the i18n service status and http status code information. This file name should be in the directory corresponding to the bundle parameter. The bundle parameter defaults to the bundle directory under the current directory.",
		},
		&cli.StringFlag{
			Name:  "i18nPkgName",
			Value: "i18n",
			Usage: "The package name of the final generated file.",
		},
		&cli.StringFlag{
			Name:  "outputDir",
			Value: "i18n",
			Usage: "The directory where the last generated file is located.",
		},
		&cli.StringFlag{
			Name:  "defaultLanguage",
			Value: "zh_cn",
			Usage: "The default language. This language needs to have a corresponding language file.",
		},
	},
	Action: func(c *cli.Context) error {

		// bundle parameter
		bundleDir := c.String("bundleDir")
		if _, err := os.Stat(bundleDir); os.IsNotExist(err) {
			return errors.New("bundle directory does not exist")
		}

		// sc parameter
		statusCodeFile := c.String("scFileName")

		var validFilename = regexp.MustCompile(`^[a-zA-Z0-9_\-.]+$`)

		if !validFilename.MatchString(statusCodeFile) {
			return errors.New("invalid filename for sc parameter")
		}

		if !strings.HasSuffix(statusCodeFile, ".toml") {
			statusCodeFile = statusCodeFile + ".toml"
		}

		// pkgName parameter
		i18nPkgName := c.String("i18nPkgName")
		var validPkgName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
		if !validPkgName.MatchString(i18nPkgName) {
			return errors.New("invalid package name")
		}

		// outDir parameter
		outputDir := c.String("outputDir")
		dirStrArray := strings.Split(outputDir, "/")
		if dirStrArray[len(dirStrArray)-1] != i18nPkgName {
			outputDir = path.Join(outputDir, i18nPkgName)
		}

		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return errors.New("failed to create outDir directory")
		}

		// defaultLanguage parameter
		defaultLanguage := c.String("defaultLanguage")
		if tag, err := language.Parse(defaultLanguage); err != nil {
			return err
		} else {
			defaultLanguage = strings.ToLower(tag.String())
		}

		cpp(bundleDir, statusCodeFile, outputDir, i18nPkgName, defaultLanguage)
		return nil
	},
}
