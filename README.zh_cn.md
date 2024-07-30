# Go I18n

## Overview

goI18n 是一个命令行工具，旨在从基于 TOML 的国际化 (i18n) 消息包生成 Go 文件。该项目通过将 TOML 文件转换为可以处理不同语言和格式的结构化 Go 代码，促进了 Go 应用程序中多语言支持的创建。

## Features

- **基于 TOML 的配置**：使用 TOML 文件定义您的 i18n 消息和状态代码。
- **多语言支持**：使用结构化消息定义轻松支持多种语言。
- **代码生成**：自动生成包含您的 i18n 消息和状态代码的 Go 文件。
- **模板处理**：支持带有动态内容占位符的模板化消息。
- **命令行**：简单直观的 CLI 来生成所需的 Go 文件。

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

我们拿项目中 `example` 目录下的代码作为示例。

1. 安装库。

```shell
go get github.com/eliassama/go-i18n
```

2. 安装生成工具。

```shell
go install github.com/eliassama/go-i18n/goI18n@latest
```

3. 在项目根目录下，使用生成命令直接生成。

```shell
 goI18n gen --bundleDir=example/bundle --scFileName=statusCode --i18nPkgName=i18n --outputDir=example/i18n --defaultLanguage=zh_cn
```

4. `example` 目录下生成了一个 `i18n` 目录和一个 `i18n.go` 的代码文件。

直接使用此文件即可，生成的代码会自动引入 `go-i18n` 的核心结构体。

## TOML Configuration

使用这个项目，最关键的就在于编写 TOML 配置文件。

TOML 配置文件分为两种，一种是 `Status And Code Toml`，另一种是 `Language Toml`。

编写的 TOML 基础格式如下:

```toml
[MessageKey] # MessageKey 就是对应的消息名称。
ItemKey1 = ItemVal2 # 这就是一个 KV 键值对
ItemKey2 = ItemVal2 # 这是第二个 KV 键值对
```

上面的格式，我们可以把他看做是一个 `map[string]string` 的 `go struct`。

为了便于我们理解这个配置文件的意思，我们把他转换为 `go` 代码。

```go
MessageKey := map[string]string{
"ItemKey1": "ItemVal2",
"ItemKey2": "ItemVal2",
}
```

`MessageKey` 其实就是这个消息的名称，也可以叫做标识。下面的内容就是消息的状态码或者中英文消息文本了。

当然，实际 `toml` 可支持的格式要比这个更复杂，更多样。但是在这个项目里，`TOML` 这样写就够了。

### Status And Code Toml

`Status And Code Toml` 主要存放的是每一个消息的状态码，分为 `Status` 和 `Code` 两个。

```toml
[TestMsgKey]
status = 555
code = 100555

[TestMsgKey2]
status = 888
code = 100888
```

> 之所以设计两个，是因为可能会涉及到 `HttpStatus` 和 `ServiceCode` 两种。
>
> 这个文件不是必须的，但我仍建议设置。因为服务并不是独立存在的，而看状态码比看异常信息更容易定位问题发生在哪里。
>
> 这个文件中你可以只写 `Status` 或者只写 `Code`，不过如果应用在一个生产级别的服务上，我建议你两者都写。

另外需要注意：

1. `Status` 或 `Code` 只能设置为大于 `0` 的值。可设置的最大值为 `int`，最小值为 `1`
2. 如果 `Status` 或 `Code` 某一个没有设置的话，则生成的代码只会返回一个状态码，如果两个都没有设置的话，则只返回消息。
3. `Status` 或 `Code` 的设置必须统一，不允许某个 `MessageKey` 设置了 `Status` 和 `Code`，而另外一个 `MessageKey` 只设置一个或者两个都不设置。
4. 这个配置文件如果存在，那 `MessageKey` 必须和语言文件的 `MessageKey` 相统一，不允许多一个或少一个。
5. 这个文件的文件名可以通过命令的 `scFileName` 参数设置，但是该文件必须和语言文件在同一个目录下。

### Language Toml

`Language Toml` 主要存放的是每一个消息的不同语言的文本。

从消息内容上，分为纯字符串和模板字符串。

从消息类型上，分为单数消息(`msg`)、复数消息(`singular`)、默认消息(`plural`)。

```toml
[NetworkAuthenticationRequiredMsg]
msg = "Network authentication required"

[TestMsg]
msg = "This is a test message. Name: {{.Name}}, Age: {{.Age}}, Home Address: {{.Phone}}"
singular = "This is a singular test message, suitable for the case where the data quantity is 1. Name: {{.Name}}, Age: {{.Age}}, Home Address: {{.Phone}}"
plural = "This is a plural test message, suitable for cases where the data quantity is greater than 1. Name: {{.Name}}, Age: {{.Age}}, Home Address: {{.Phone}}"
```

> `Language Toml` 的文件名称格式为 `[模块或其他分类命名-可选].[语言代码-必选].toml`
>
> 有效的文件名称示例：`en_us.toml`、`zh_cn.toml`、`user.zh_cn.toml`、`user.en_us.toml`。

另外需要注意：

1. 允许单独只设置 `msg`，允许同时设置 `msg`、`singular`、`plural`，也允许只设置 `singular`、`plural`。但是不允许单独设置 `singular` 或 `plural`。因为 `singular` 和 `plural` 是单复数的形式，不能只出现一个。
2. `msg`、`singular`、`plural` 只能设置为非空字符串。
3. 如果 `msg`、`singular`、`plural` 都设置的话，则只会根据 `count` 决定返回 `singular` 还是 `plural`，不会返回 `msg`。
4. `msg`、`singular`、`plural` 的设置必须统一，不允许某个 `MessageKey` 设置了 `msg`，而另外一个 `MessageKey` 没有设置 `msg`。`singular` 和 `plural` 同理。
5. 不同语言文件的 `MessageKey` 必须相统一，不允许一个语言有，而另一个语言没有。
6. 消息支持的模板语言，请参考 `text/template` 包。
7. 语言文件的目录可以通过命令的 `bundleDir` 参数设置。

> TOML 文件示例可以参考 [example](example) 目录。

## CLI Usage

命令行工具比较简单，因为咱们是极简的嘛。（笑）

安装命令如下：

```shell
go get github.com/eliassama/go-i18n
```

示例完整命令：

```shell
goI18n gen --bundleDir=example/bundle --scFileName=statusCode --i18nPkgName=i18n --outputDir=example/i18n --defaultLanguage=zh_cn
```

### CLI Command And Parameter

命令行的命令目前只有一个 gen 是用来生成的。

咱们主要来介绍下命令行的四个参数：

> **defaultLanguage** 参数 **『可选』**
>
> 简单介绍：用来指定默认语言包，防止生成的方法中，传入非有效语言代码，导致无法返回内容。这个默认语言包，必须是存在相关的语言包 toml。比如只有 zh_cn 和 en_us 的语言包 toml，但是指定默认语言包为 zh_hk，就不行。
>
> 默认值：zh_cn
>
> 示例：`goI18n gen --defaultLanguage=zh_cn`，默认语言为 zh_cn


> **bundleDir** 参数 **『可选』**
>
> 简单介绍：用来指定语言包的 Toml 文件路径，这个路径下还应该存放 statusCode.toml（状态码）
>
> 默认值：bundle，命令行执行目录下的 `./bundle/` 目录
>
> 示例：`goI18n gen --bundle=example/bundle`，命令行执行目录下的 `./bundle/` 目录

> **scFileName** 参数 **『可选』**
>
> 简单介绍：用来指定状态码的文件名称，这个文件必须存放在 bundleDir 目录下。如果没有指定名称的文件，则最后生成的代码，不会返回任何状态码。这个文件名，如果不以 .toml 结尾的话，会自动加上 .toml。文件名需要符合正则：`^[a-zA-Z0-9_\-.]+$`
>
> 默认值：statusCode，对应的就是 statusCode.toml
>
> 示例：goI18n gen --scFileName=statusCode 或 goI18n gen --scFileName=statusCode.toml ，两者是等价的。

> **i18nPkgName** 参数 **『可选』**
>
> 简单介绍：生成的 i18n 语言 go 代码文件的 package 名称、文件名、所在目录名。需要符合正则：`^[a-zA-Z_][a-zA-Z0-9_]*$`
>
> 默认值：i18n，会在下面的 outputDir 指定的目录下创建一个 i18n 目录，在 i18n 目录中创建 i18n.go 并生成代码填充，然后 i18n.go 的 package name 就是 i18n。如果 outputDir 也是 i18n 的话，就直接只创建 i18n。
>
> 示例：goI18n gen --i18nPkgName=i18n

> **outputDir** 参数 **『可选』**
>
> 简单介绍：生成的 i18n 语言 go 代码文件输出的前置路径。如果和 i18nPkgName 同名的话，就忽略。如果最后的路径和 i18nPkgName 同名的话，则最后的路径也忽略。
>
> 默认值：i18n。直接创建一个 i18n 目录，生成的代码输出到这个目录里。
>
> 示例：goI18n gen --i18nPkgName=example

`i18nPkgName` 和 `outputDir` 的相互作用关系，我这里举几个例子来说明下。

> - **outputDir 为 example，i18nPkgName 为 i18n**：
    >   - 在当前目录下，创建 `example` 目录，`example` 目录下创建 `i18n` 目录。在 `i18n` 目录中创建 `i18n.go` 并生成代码填充。

> - **outputDir 为 i18n，i18nPkgName 为 i18n**：
    >   - 在当前目录下，创建 `i18n` 目录。在 `i18n` 目录中创建 `i18n.go` 并生成代码填充。

> - **outputDir 为 example/i18n，i18nPkgName 为 i18n**：
    >   - 在当前目录下，创建 `example` 目录，`example` 目录下创建 `i18n` 目录。在 `i18n` 目录中创建 `i18n.go` 并生成代码填充。

## maintains

[@eliassama](https://github.com/eliassama)
