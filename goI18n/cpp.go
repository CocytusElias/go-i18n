package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/eliassama/go-i18n/core"
	"golang.org/x/text/language"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type I18nMessage struct {
	Msg      string `toml:"msg,omitempty" json:"msg,omitempty"`
	Plural   string `toml:"plural,omitempty" json:"plural,omitempty"`
	Singular string `toml:"singular,omitempty" json:"singular,omitempty"`
}

type I18nMessages map[string]I18nMessage

type I18nStatusCode struct {
	Status int `toml:"status,omitempty"  json:"status,omitempty"`
	Code   int `toml:"code,omitempty"    json:"code,omitempty"`
}

type I18nStatusCodes map[string]I18nStatusCode

func cpp(bundleDir, statusCodeFile, outputDir, i18nPkgName, defaultLanguage string) {

	i18nMessages, langCodes := readI18nMessages(bundleDir, defaultLanguage)
	i18nStatusCodes := readI18nStatusCodes(bundleDir, statusCodeFile)
	i18nBundles, statusExist, codeExist := mergeVerifyI18nBundle(i18nMessages, langCodes, i18nStatusCodes)

	generateGoFile(i18nBundles, outputDir, i18nPkgName, defaultLanguage, statusExist, codeExist)
}

func readI18nMessages(bundleDir, defaultLanguage string) (map[string]I18nMessages, map[string]bool) {
	messages := make(map[string]I18nMessages)
	langCodes := make(map[string]bool)
	defaultLanguageExist := false

	files, err := os.ReadDir(bundleDir)
	if err != nil {
		panic(fmt.Sprintf("Failed to read I18n Bundel directory: %v", err))
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			filePath := filepath.Join(bundleDir, fileName)
			langCode := strings.ToLower(getLangCodeFromFileName(fileName))

			if langCode != "" {
				if langCode == defaultLanguage {
					defaultLanguageExist = true
				}

				langCodes[langCode] = true

				bundleMessages := make(I18nMessages)
				if _, err = toml.DecodeFile(filePath, &bundleMessages); err != nil {
					panic(fmt.Sprintf("Failed to decode I18n Bundel File 『%s』: %v", fileName, err))
				}

				if len(bundleMessages) == 0 {
					continue
				}

				for messageKey, messageVal := range bundleMessages {
					if messageKey == "" {
						panic(fmt.Sprintf("In the 『%s』 file under the 『%s』 path, there is an empty identifier", fileName, filePath))
					}

					messageVal.Msg = strings.TrimSpace(messageVal.Msg)
					messageVal.Singular = strings.TrimSpace(messageVal.Singular)
					messageVal.Plural = strings.TrimSpace(messageVal.Plural)

					if messageVal.Msg == "" && messageVal.Singular == "" && messageVal.Plural == "" {
						panic(fmt.Sprintf("In the 『%s』 file under the 『%s』 path, the language information marked by 『%s』 is empty", fileName, filePath, messageKey))
					}

					if _, exist := messages[messageKey]; !exist {
						messages[messageKey] = make(map[string]I18nMessage)
					}
					messages[messageKey][langCode] = messageVal
				}
			}
		}
	}

	if !defaultLanguageExist {
		panic("The default language is missing, please check your language bundle file.")
	}

	return messages, langCodes
}

func readI18nStatusCodes(bundleDir string, statusCodeFile string) I18nStatusCodes {
	statusCodes := make(I18nStatusCodes)

	filePath := path.Join(bundleDir, statusCodeFile)
	if _, err := toml.DecodeFile(filePath, &statusCodes); err != nil {
		return nil
	}

	return statusCodes
}

func getLangCodeFromFileName(filename string) string {
	parts := strings.Split(filename, ".")
	n := len(parts)

	if n >= 2 {
		filename = strings.Join(parts[n-2:], ".")
	} else {
		return ""
	}

	if languageTag, err := language.Parse(parts[n-2]); err == nil {
		return languageTag.String()
	}

	return ""
}

func mergeVerifyI18nBundle(i18nMessages map[string]I18nMessages, langCodes map[string]bool, i18nStatusCodes I18nStatusCodes) (map[string]*core.Bundle, bool, bool) {
	var i18nBundle = make(map[string]*core.Bundle)

	var i18nStatusMap = make(map[string]int)
	var i18nCodeMap = make(map[string]int)
	var i18nMessageKeys = make([]string, 0)

	for key, message := range i18nMessages {
		i18nMessageKeys = append(i18nMessageKeys, key)

		i18nBundle[key] = &core.Bundle{}

		if i18nStatusCodes != nil || len(i18nStatusCodes) != 0 {
			if _, ok := i18nStatusCodes[key]; !ok {
				panic(fmt.Sprintf("Identifier 『%s』 declares the message content, but does not declare status and code", key))
			}

			if i18nStatusCodes[key].Status > 0 {
				i18nStatusMap[key] = i18nStatusCodes[key].Status
				i18nBundle[key].Status = i18nStatusCodes[key].Status
			}

			if i18nStatusCodes[key].Code > 0 {
				i18nCodeMap[key] = i18nStatusCodes[key].Code
				i18nBundle[key].Code = i18nStatusCodes[key].Code
			}
		}

		verifyMsgMap := make(map[string]bool)
		verifyPluralMap := make(map[string]bool)
		verifySingularMap := make(map[string]bool)

		for langCode := range langCodes {

			if i18nMessage, ok := message[langCode]; !ok {
				panic(fmt.Sprintf("There is no message in language 『%s』 under identifier 『%s』, please check", langCode, key))
			} else {
				if len(i18nBundle[key].Msg) == 0 {
					i18nBundle[key].Msg = make(map[string]string)
				}

				if len(i18nBundle[key].Singular) == 0 {
					i18nBundle[key].Singular = make(map[string]string)
				}

				if len(i18nBundle[key].Plural) == 0 {
					i18nBundle[key].Plural = make(map[string]string)
				}

				if i18nMessage.Msg != "" {
					verifyMsgMap[langCode] = true
					i18nBundle[key].Msg[langCode] = i18nMessage.Msg
				}

				if i18nMessage.Singular != "" {
					verifySingularMap[langCode] = true
					i18nBundle[key].Singular[langCode] = i18nMessage.Singular
				}

				if i18nMessage.Plural != "" {
					verifyPluralMap[langCode] = true
					i18nBundle[key].Plural[langCode] = i18nMessage.Plural
				}

				if (i18nMessage.Singular == "" && i18nMessage.Plural != "") || (i18nMessage.Singular != "" && i18nMessage.Plural == "") {
					panic(fmt.Sprintf("Plural and Singular message contents must be set at the same time. Because Plural and Singular are messages in singular and plural form.\nIf you only want to set one message, regardless of singular or plural, please set Msg message.\nPlease adjust the message content under the %s identifier.", key))
				}

			}

		}

		checkVerifyMapErrMsg := fmt.Sprintf("type message settings for the 『%s』 identification language are not uniform:", key)

		if err := checkVerifyLangCodeMap(langCodes, verifyMsgMap); err != nil {
			panic(fmt.Sprintf("The 『Msg』 %s %s", checkVerifyMapErrMsg, err.Error()))
		}

		if err := checkVerifyLangCodeMap(langCodes, verifyPluralMap); err != nil {
			panic(fmt.Sprintf("The 『Plural』 %s %s", checkVerifyMapErrMsg, err.Error()))
		}

		if err := checkVerifyLangCodeMap(langCodes, verifySingularMap); err != nil {
			panic(fmt.Sprintf("The 『Singular』 %s %s", checkVerifyMapErrMsg, err.Error()))
		}

	}

	if i18nStatusCodes != nil || len(i18nStatusCodes) == 0 {
		uselessKeysByStatusCode := make([]string, 0)
		for key := range i18nStatusCodes {
			if _, exist := i18nBundle[key]; !exist {
				uselessKeysByStatusCode = append(uselessKeysByStatusCode, key)
			}
		}

		if len(uselessKeysByStatusCode) > 0 {
			panic(fmt.Sprintf("The identifier only exists in the statusCode toml configuration file, but not in the language configuration file. Please delete these keys: 『%s』", strings.Join(uselessKeysByStatusCode, ",")))
		}
	}

	var codeExist = false
	var statusExist = false

	if len(i18nStatusMap) != 0 {
		statusExist = true
		if err := checkVerifyStatusCodeMap(i18nMessageKeys, i18nStatusMap); err != nil {
			panic(fmt.Sprintf("The Toml configuration for storing Status Code has a validation error. There is a problem with 『status』: %s", err.Error()))
		}
	}

	if len(i18nCodeMap) != 0 {
		codeExist = true
		if err := checkVerifyStatusCodeMap(i18nMessageKeys, i18nCodeMap); err != nil {
			panic(fmt.Sprintf("The Toml configuration for storing Status Code has a validation error. There is a problem with 『code』: %s", err.Error()))
		}
	}

	return i18nBundle, statusExist, codeExist

}

func checkVerifyLangCodeMap(langCodes, verifyMap map[string]bool) error {
	if len(verifyMap) == 0 {
		return nil
	}

	present := make([]string, 0)
	missing := make([]string, 0)
	for code := range langCodes {
		if _, exists := verifyMap[code]; exists {
			present = append(present, code)
		} else {
			missing = append(missing, code)
		}
	}

	if len(missing) > 0 {
		return errors.New(fmt.Sprintf("The corresponding message content is set in the 『%s』 language, but not in the 『%s』 language.", strings.Join(present, ","), strings.Join(missing, ",")))
	}

	return nil
}

func checkVerifyStatusCodeMap(i18nMessageKeys []string, i18nStatusCodeMap map[string]int) error {
	var i18nMessageKeyNotExistStatusCode = make([]string, 0)

	for _, i18nMessageKey := range i18nMessageKeys {
		if _, ok := i18nStatusCodeMap[i18nMessageKey]; !ok {
			i18nMessageKeyNotExistStatusCode = append(i18nMessageKeyNotExistStatusCode, i18nMessageKey)
		}
	}

	if len(i18nMessageKeyNotExistStatusCode) > 0 {
		return errors.New(fmt.Sprintf("These messages Key 『%s』 do not exist", strings.Join(i18nMessageKeyNotExistStatusCode, ",")))
	}
	return nil
}
