<a id="markdown-settings---typed-configuration-toolkit-for-go" name="settings---typed-configuration-toolkit-for-go"></a>
# Settings - Typed Configuration Toolkit For Go
[![GoDoc](https://godoc.org/github.com/asecurityteam/settings?status.svg)](https://godoc.org/github.com/asecurityteam/settings)
[![Build Status](https://travis-ci.com/asecurityteam/settings.png?branch=master)](https://travis-ci.com/asecurityteam/settings)
[![codecov.io](https://codecov.io/github/asecurityteam/settings/coverage.svg?branch=master)](https://codecov.io/github/asecurityteam/settings?branch=master)

*Status: Incubation*

<!-- TOC -->

- [Settings - Typed Configuration Toolkit For Go](#settings---typed-configuration-toolkit-for-go)
    - [Overview](#overview)
    - [Data Sources](#data-sources)
    - [Component API](#component-api)
    - [Hierarchy API](#hierarchy-api)
    - [Adapter API](#adapter-api)
    - [Contributing](#contributing)
        - [License](#license)
        - [Contributing Agreement](#contributing-agreement)

<!-- /TOC -->

<a id="markdown-overview" name="overview"></a>
## Overview

There aren't very many well tested and maintained libraries in the ecosystem for
managing configuration values in a go application. The two that appear with the most
prevalence when searching are [viper](https://github.com/spf13/viper) and
[go-config](https://github.com/micro/go-config). Both of these projects provide both
dynamic configuration reloading and a broad interface for fetching and converting
configuration values. Both of those libraries are fairly stable and excellent choices
for configuration if they fit your needs.

This project grew out of a desire for a configuration system that:

-   Allows for the minimum amount of coupling between components needing
    configuration and the configuration system.

-   Focuses on enabling plugin, or at least swappable component, based systems.

-   Enables developers to define configuration as standard go types rather than as
    strings in a global, opaque registrar.

-   Provides the ability to remix the basic configuration components to create new
    complex systems that may differ from our initial assumptions about how we need
    configuration to work in a given system.

-   Generates useful help output and example configurations from a given set of
    options.

Honestly, [viper] and [go-config] both come close to satisfying most of these wants
but fall short of what we needed in small, but consequential, ways. This project
attempts to overcome those small deficits by offering:

-   An extremely minimal interface for defining sources of configuration data so that
    new or proprietary data sources can be more easily added.

-   A high level API for working with loosely coupled components that describe their
    configurations with standard go types.

-   A mid level API for defining and managing configuration hierarchies.

-   A low level API for interfacing between statically typed options and weakly typed
    configuration sources.

<a id="markdown-data-sources" name="data-sources"></a>
## Data Sources

All forms of our API interact in some with a source of configuration data. The
interface for a source is:

```golang
type Source interface {
    Get(ctx context.Context, path []string) (interface{}, error)
}
```

It is intentionally simple and leaves every implementation detail to the the team
creating a new source. Packaged with this project are Source implementations for
JSON, YAML, and ENV. We also provide a minimal set of tools for arranging and
composing data sources:

```golang
jsonSource, _ := settings.NewFileSource("config.json")
yamlSource, _ := settings.NewFileSOurce("config.yaml")
envSource := settings.NewEnvSource(os.Environ())
// Apply a static prefix to all env lookups. For example,
// here we add a APP_ to all names.
prefixedEnv := &settings.PrefixSource{
    Source: envSource,
    Prefix: []string{"APP"},
}
// Create an ordered lookup where ENV is first with a fallback to
// YAML which can fallback to JSON.
finalSource := []settings.MultiSource{prefixedEnv, yamlSource, jsonSource}

v, found := finalSource.Get(context.Background(), "setting")
```

Sources may be used as-is by passing them around to components that need to fetch
values. However, the values returned from `Get()` are opaque and highly dependent on
the implementation. For example, the ENV source will always return a string
representation of the value because that is what is available in the environment.
Alternatively, the JSON and YAML sources may return other data types as they
typically unmarshal into native go types. Each component fetching values from a
source is responsible for safely converting the result into a useful value.

We recommend using one of the API layers we provide to do this for you.

<a id="markdown-component-api" name="component-api"></a>
## Component API

With one of our goals being the support of plugin based systems, we've built
configurable components into the higher level interface of the project.
`NewComponent` is the entry point for the high-level, Component API. This method
manages much of the complexity of adding configuration to a system.

```golang
func NewComponent(ctx context.Context, s settings.Source, value interface{}, destination interface{}) error
```

The given context and source are used for all lookups of configuration data. The
given value must be an implementation of the component contract and the destination
is a pointer created with `new(T)` where `T` is the output type (or equivalent
convertible) to the element produced by the component contract implementation.

The component contract is an interface that all input values must conform to and is
roughly equivalent to the Factory or Constructor concepts. Each instance of the
component contract must define two methods: `Setting() C` and `New(context.Context,
C) (T, error)`. Due to the lack of generics in go, there's no way to describe this
contract as an actual go interface that would benefit from static typing support. As
a result, `NewComponent` uses reflection to enforce the contract in order to allow for
`C` to be any type that is convertible to configuration via the `settings.Convert()`
method and for `T` to be any type that your use case requires.

For example, the most minimal implementation of the contract would look like:

```golang
// All configuration is driven by a struct that uses standard go types.
type Config struct {}
// The resulting element is virtually anything that you need it to be
// for the purposes of the system you are building.
type Result struct {}
// The component contract interfaces between the thing you want to make,
// the result, and the settings project. It is responsible for producing
// instances of the configuration struct with any default values populated.
// It is also used to construct new instances of the result using a
// populated configuration. No references to the settings project are
// required to create any part of this setup.
type Component struct {}
func (*Component) Settings() *Config { return &Config{} }
func (*Component) New(_ context.Context, c *Config) (*Result, error) {
    return &Result{}, nil
}
```

From here, any number of settings and sub-trees may be added to `Config`, any methods
or attributes may be added to `Result`, and any complexity in the creation of
`Result` may be be added to the `Component.New` method. To then use this basic
example as a component you would:

```golang
component := &Component{}
r := new(Result)
err := settings.NewComponent(context.Background(), source, component, r)
```

If the resulting error is `nil` then the destination value, `r` in this example, now
points to the output of the `Component.New` method. The method returns an error any
time the given component does not satisfy the contract, any time the configuration
loading fails, or the `Component.New` returns an error.

The benefits of using this API are that it is highly flexible with respect to types
and it prevents plugins or components from needing to import and using elements from
this project. This makes it a bit easier to write tests by removing the need to
orchestrate an entire configuration system.

A potential downside to this API is that the resulting configuration hierarchy is not
easily modified. The structure is enforced is such that each component receives a top
level key and all nested structs result in sub-trees. The name of every setting is
generated from the field name and this is not changeable. The description of each
field can be set using struct tags. The name and description of each tree may be
defined by implementing a `Name()` and `Description()` method but the overall
arrangement is fixed.

```golang
type InnerConfig struct {
    Value2 string `description:"a string"`
}
func (c *InnerConfig) Name() string {
    return "subtree"
}
func (c *InnerConfig) Description() string {
    return "a nesting configuration tree"
}

type OuterConfig struct {
    Value1 int `description:"the first value"`
    InnerConfig *InnerConfig
}
func (c *OuterConfig) Name() string {
    return "toptree"
}
func (c *OuterConfig) Description() string {
    return "the top configuration tree"
}
```

The above will equate to a configuration like:

```yaml
toptree:
    value1: 0
    subtree:
        value2: ""
```

The descriptions are used to annotate example configurations and help output.

<a id="markdown-hierarchy-api" name="hierarchy-api"></a>
## Hierarchy API

If the component API is too restrictive for your use case then the Hierarchy API
might be of more use. This layer of the API is based on the `settings.Setting` and
`settings.Group` types which represent individual configuration options and
sub-trees, respectively. We include a `settings.SettingGroup` implementation of the
`settings.Group` which allows you to construct any arbitrary hierarchy of
configuration as needed. It also, however, requires coupling the code to this project
and using a much more verbose style of defining options:

```golang
top := &settings.SettingGroup{
    NameValue: "root",
    GroupValues: []settings.Group{}, // Add any sub-tree here.
    SettingValues: []settings.Setting{ // Add any settings here.
        settings.NewIntSetting("Value1", "an int value", 2),
    },
}

err := settings.LoadGroups(finalSource, []Group{top})
```

After loading is complete, each `Setting` value will contain either the given default
for a the value found in the `Source`. This is the same API we used to create the
Component API.

<a id="markdown-adapter-api" name="adapter-api"></a>
## Adapter API

If none of the higher API layers provide what you need then we also offer a lower
level tool set for creating new high level APIs. At the most basic level, this
project contains a large set of strongly types adapters for content pulled from a
source. Each adapter is responsible for converting from the empty interface into a
native type that the setting exposes. These are named with the pattern of
`<Type>Setting`. We make use of the [`cast`](https://github.com/spf13/cast) project
to handle converting from arbitrary types to target types.

This is the layer to target when adding new supported configuration types or when
replacing the type converters with something else. These elements, in possible
conjunction with elements from the Hierarchy API, are flexible enough to build
anything you need.

## Special Type Parsing and Casting

We use the `cast` library for casting values read in from configurations into their go types. The `cast` library 
falls back to JSON for complex types expressed as string values. Here are some examples of how we parse different types:

**[]string**

For a given configuration
```go
type Config struct {
    TheSlice []string
}
```
The values in the following examples will all be parsed as a string slice.

*yaml*
```yaml
config:
  theslice:
    - "a"
    - "b"
    - "c"
```

*JSON*
```json
{"config": {"theslice":  ["a", "b", "c"]}}
```

You can also set an environment variable and reference them in a YAML or JSON file like below. Note that this
environment variable value will be parsed as a slice where each letter will be a value since it gets split by any
space in the string.

*Environment Variable*
```shell
CONFIG_THESLICE="a b c"`
```

*yaml*
```yaml
config:
  theslice: "${CONFIG_THESLICE}"
```

*JSON*
```json
{"config": {"theslice":  "${CONFIG_THESLICE}"}}
```

**map[string][]string**

For a given configuration
```go
type Config struct {
    allowedStrings map[string][]string
}
```

The values in the following examples will all be parsed as a string map string slices where the key `letters` and
`symbols` gets included as the string map key and their values are a string slice.

*yaml*
```yaml
config:
  allowedStrings:
    letters:
      - "a"
      - "b"
      - "c"
    symbols:
      - "@"
      - "!"
```

*JSON*
```json
{
	"config": {
		"allowedStrings": {
			"letters": ["a", "b", "c"],
			"symbols": ["@", "!"]
		}
	}
}
```

**time.Time**

For a given configuration
```go
type Config struct {
    TheTime time.Time
}
```

The following examples will be parsed using the RFC3339 format by `time.Parse(time.RFC3339, value)`

*yaml*
```yaml
"config":
  "thetime": "2012-11-01T22:08:41+00:00"
```

*JSON*
```json
{"config": {"thetime": "2012-11-01T22:08:41+00:00"}}
```

*Environment Variable*
```shell
CONFIG_THETIME="2012-11-01T22:08:41+00:00"`
```

**time.Duration**

For a given configuration
```go
type Config struct {
	TimeLength time.Duration
}
```

The following examples will be parsed using `time.Duration`
*yaml*
```yaml
"config":
  "timeLength": "4h"
```

*JSON*
```json
{"config": {"timeLength": "4h"}}
```

*Environment Variable*
```shell
CONFIG_TIMELENGTH="4h"
```

<a id="markdown-contributing" name="contributing"></a>
## Contributing

<a id="markdown-license" name="license"></a>
### License

This project is licensed under Apache 2.0. See LICENSE.txt for details.

<a id="markdown-contributing-agreement" name="contributing-agreement"></a>
### Contributing Agreement

Atlassian requires signing a contributor's agreement before we can accept a patch. If
you are an individual you can fill out the [individual
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d).
If you are contributing on behalf of your company then please fill out the [corporate
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b).
