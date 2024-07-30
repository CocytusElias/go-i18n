# Go I18n

## Overview

goI18n is a command-line tool designed to generate Go files from TOML-based internationalization (i18n) message packages. This project facilitates the creation of multi-language support in Go applications by converting TOML files into structured Go code that can handle different languages and formats.
## Features

- **TOML-based configuration**: Define your i18n messages and status codes using TOML files.
- **Multi-language support**: Easily support multiple languages using structured message definitions.
- **Code generation**: Automatically generate Go files containing your i18n messages and status codes.
- **Template handling**: Supports templated messages with dynamic content placeholders.
- **Command line**: Simple and intuitive CLI to generate the required Go files.
&nbsp;

---

## TOC

- [Quick Start](#quick-start)
- [TOML Configuration](#toml-configuration)
  - [Status And Code Toml](#status-and-code-toml)
  - [Language Toml](#language-toml)
- [CLI Usage](#cli-usage)
  - [CLI Command And Parameter](#cli-command-and-parameter)
- [maintains](#maintains)

## Quick Start

Let's take the code in the `example` directory of the project as an example.

1. Install the library.

```shell
go get github.com/eliassama/go-i18n
```

2. Install the build tool.

```shell
go install github.com/eliassama/go-i18n/goI18n@latest
```

3. In the project root directory, use the build command to build directly.

```shell
goI18n gen --bundleDir=example/bundle --scFileName=statusCode --i18nPkgName=i18n --outputDir=example/i18n --defaultLanguage=en_us
```

4. An `i18n` directory and a `i18n.go` code file are generated in the `example` directory.

You can use this file directly, and the generated code will automatically introduce the core structure of `go-i18n`.

## TOML Configuration

The key to using this project is to write a TOML configuration file.

There are two types of TOML configuration files, one is `Status And Code Toml`, and the other is `Language Toml`.

The basic format of the written TOML is as follows:

```toml
[MessageKey] # MessageKey is the corresponding message name.
ItemKey1 = "ItemVal2" # This is a KV key-value pair
ItemKey2 = "ItemVal2" # This is the second KV key-value pair
```

In the above format, we can regard it as a `go struct` of `map[string]string`.

In order to make it easier for us to understand the meaning of this configuration file, we convert it to `go` code.

```
MessageKey := map[string]string{
    "ItemKey1": "ItemVal2",
    "ItemKey2": "ItemVal2",
}
```

`MessageKey` is actually the name of the message, which can also be called an identifier. 

The following content is the status code of the message or the Chinese and English message text.

Of course, the actual `toml` supported formats are more complex and diverse than this. But in this project, `TOML` is enough.

### Status And Code Toml

`Status And Code Toml` mainly stores the status code of each message, which is divided into `Status` and `Code`.

```toml
[TestMsgKey]
status = 555
code = 100555

[TestMsgKey2]
status = 888
code = 100888
```

> The reason for designing two is that `HttpStatus` and `ServiceCode` may be involved.
>
> This file is not required, but I still recommend setting it up. Because services do not exist independently, it is easier to locate the problem by looking at the status code than by looking at the exception information.
>
> In this file, you can only write `Status` or only write `Code`, but if it is applied to a production-level service, I recommend you write both.

Also note:
1. `Status` or `Code` can only be set to a value greater than `0`. The maximum value that can be set is `int` and the minimum value is `1`

2. If one of `Status` or `Code` is not set, the generated code will only return a status code. If neither is set, only the message will be returned.

3. The settings of `Status` or `Code` must be unified. It is not allowed for a `MessageKey` to set `Status` and `Code`, while the other `MessageKey` only sets one or neither.

4. If this configuration file exists, the `MessageKey` must be unified with the `MessageKey` of the language file. One more or one less is not allowed.

5. The file name of this file can be set by the `scFileName` parameter of the command, but the file must be in the same directory as the language file.

### Status And Code Toml

`Status And Code Toml` mainly stores the status code of each message, which is divided into `Status` and `Code`.

```toml
[TestMsgKey]
status = 555
code = 100555

[TestMsgKey2]
status = 888
code = 100888
```

> The reason for designing two is that `HttpStatus` and `ServiceCode` may be involved.
>
> This file is not necessary, but I still recommend setting it. Because services do not exist independently, it is easier to locate the problem by looking at the status code than by looking at the exception information.
>
> In this file, you can only write `Status` or only write `Code`, but if it is applied to a production-level service, I recommend you write both.

Also note:

1. `Status` or `Code` can only be set to a value greater than `0`. The maximum value that can be set is `int`, and the minimum value is `1`

2. If one of `Status` or `Code` is not set, the generated code will only return a status code. If neither is set, only the message will be returned.

3. The settings of `Status` or `Code` must be consistent. It is not allowed for a `MessageKey` to set `Status` and `Code`, while the other `MessageKey` only sets one or neither.

4. If this configuration file exists, the `MessageKey` must be consistent with the `MessageKey` of the language file. One more or one less is not allowed.

5. The file name of this file can be set by the `scFileName` parameter of the command, but the file must be in the same directory as the language file.

### Language Toml
`Language Toml` mainly stores the text of each message in different languages.

In terms of message content, it is divided into pure string and template string.

In terms of message type, it is divided into singular message (`msg`), plural message (`singular`), and default message (`plural`).

```toml
[NetworkAuthenticationRequiredMsg]
msg = "Network authentication required"

[TestMsg]
msg = "This is a test message. Name: {{.Name}}, Age: {{.Age}}, Home Address: {{.Phone}}"
singular = "This is a singular test message, suitable for the case where the data quantity is 1. Name: {{.Name}}, Age: {{.Age}}, Home Address: {{.Phone}}"
plural = "This is a plural test message, suitable for cases where the data quantity is greater than 1. Name: {{.Name}}, Age: {{.Age}}, Home Address: {{.Phone}}"
```

> The file name format of `Language Toml` is `[Module or other category name - optional].[Language code - required].toml`
>
> Examples of valid file names: `en_us.toml`, `zh_cn.toml`, `user.zh_cn.toml`, `user.en_us.toml`.

Other notes:

1. It is allowed to set only `msg`, it is allowed to set `msg`, `singular`, `plural` at the same time, and it is also allowed to set only `singular`, `plural`. However, it is not allowed to set only `singular` or `plural`. Because `singular` and `plural` are singular and plural forms, only one cannot appear.
2. `msg`, `singular`, `plural` can only be set to non-empty strings.
3. If `msg`, `singular`, `plural` are all set, only `singular` or `plural` will be returned based on `count`, and `msg` will not be returned.
4. The settings of `msg`, `singular`, and `plural` must be consistent. It is not allowed for a `MessageKey` to have `msg` set while another `MessageKey` does not have `msg` set. The same applies to `singular` and `plural`.
5. The `MessageKey` of different language files must be consistent. It is not allowed for one language to have it while another language does not.
6. For template languages supported by messages, please refer to the `text/template` package.
7. The directory of the language file can be set by the `bundleDir` parameter of the command.

> For examples of TOML files, please refer to the [example](example) directory.

## CLI Usage

The command line tool is relatively simple, because we are minimalist. (Laugh)

The installation command is as follows:

```shell
go get github.com/eliassama/go-i18n
```

Example full command:

```shell
goI18n gen --bundleDir=example/bundle --scFileName=statusCode --i18nPkgName=i18n --outputDir=example/i18n --defaultLanguage=zh_cn
```

### CLI Command And Parameter

Currently, there is only one command in the command line, gen, which is used for generation. Let's mainly introduce the four parameters of the command line:

> **defaultLanguage** parameter **『Optional』**
>
> Brief introduction: used to specify the default language package to prevent the generation method from passing in invalid language codes, resulting in failure to return content. This default language package must have a related language package toml. For example, if there are only zh_cn and en_us language package tomls, but the default language package is specified as zh_hk, it will not work.
>
> Default value: zh_cn
>
> Example: `goI18n gen --defaultLanguage=zh_cn`, the default language is zh_cn

> **bundleDir** parameter **『Optional』**
>
> Brief introduction: used to specify the Toml file path of the language pack, and statusCode.toml (status code) should also be stored in this path
>
> Default value: bundle, the `./bundle/` directory under the command line execution directory
>
> Example: `goI18n gen --bundle=example/bundle`, the `./bundle/` directory under the command line execution directory

> **scFileName** parameter **『Optional』**
>
> Brief introduction: used to specify the file name of the status code, which must be stored in the bundleDir directory. If there is no file with the specified name, the generated code will not return any status code. If this file name does not end with .toml, .toml will be automatically added. The file name needs to conform to the regular expression: `^[a-zA-Z0-9_\-.]+$`
>
> Default value: statusCode, corresponding to statusCode.toml
>
> Example: goI18n gen --scFileName=statusCode or goI18n gen --scFileName=statusCode.toml , the two are equivalent.

> **i18nPkgName** Parameter **『Optional』**
>
> Brief introduction: The package name, file name, and directory name of the generated i18n language go code file. Need to conform to the regular expression: `^[a-zA-Z_][a-zA-Z0-9_]*$`
>
> Default value: i18n, an i18n directory will be created in the directory specified by outputDir below, i18n.go will be created in the i18n directory and code filling will be generated, and the package name of i18n.go will be i18n. If outputDir is also i18n, just create i18n directly.
>
> Example: goI18n gen --i18nPkgName=i18n

> **outputDir** Parameter **『Optional』**
>
> Brief introduction: The leading path of the generated i18n language go code file output. If it has the same name as i18nPkgName, it is ignored. If the last path has the same name as i18nPkgName, the last path is also ignored.
>
> Default value: i18n. Create an i18n directory directly and output the generated code to this directory.
>
> Example: goI18n gen --i18nPkgName=example

The interaction between `i18nPkgName` and `outputDir`, I will give a few examples to illustrate.

> - **outputDir is example, i18nPkgName is i18n**:
> - Create an `example` directory in the current directory, and an `i18n` directory in the `example` directory. Create `i18n.go` in the `i18n` directory and generate code polyfills.

> - **outputDir is i18n, i18nPkgName is i18n**:
> - Create an `i18n` directory in the current directory. Create `i18n.go` in the `i18n` directory and generate code polyfills.

> - **outputDir is example/i18n, i18nPkgName is i18n**:
> - Create an `example` directory in the current directory, and an `i18n` directory in the `example` directory. Create `i18n.go` in the `i18n` directory and generate code polyfills.

## maintains

[@eliassama](https://github.com/eliassama)
