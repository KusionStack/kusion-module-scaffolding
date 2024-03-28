# Kusion Module Scaffolding

*This is a quick start template repository for developing a Kusion module, which is built on the [Kusion Module Framework](https://github.com/KusionStack/kusion-module-framework).*

This template repository is intended as a starting point for creating a Kusion Module, containing: 

- An example Kusion Module **kawesome** (`kawesome.k` and `src/`) with both schema definition and generator implementation. It will generate a Kubernetes `Service` and a Terraform `random_password` resource. 

- An example application (`example/`) with both developer configuration codes and platform workspace configurations using the kawesome module.

## Pre-requisite

- Go 1.22 or later

## Directory Structure

:::tip

For more details on what Kusion Module is, please refer to the [concept doc](https://www.kusionstack.io/docs/kusion/concepts/kusion-module).
:::

A Kusion Module contains the following:
- A `kcl.mod` file describing the module metadata, such as `name`, `version`, etc
- A `*.k` file with the KCL schema definition. This KCL schema will be what the targeted module users see. In the context of platform engineering, we recommended to only expose concepts that are well-known to them.
- A `src` directory containing the generator implementation written in `Go`. This should be a complete, build-able `Go` project with unit tests. The `Makefile` provides a way to build and test the module locally without publishing them.
- An optional `example` directory containing a complete Kusion project that serves as a sample when using this Kusion Module.

```
├── example
│   ├── dev
│   │   ├── example_workspace.yaml
│   │   ├── kcl.mod
│   │   ├── main.k
│   │   └── stack.yaml
│   └── project.yaml
├── kawesome.k
├── kcl.mod
└── src
    ├── Makefile
    ├── go.mod
    ├── go.sum
    ├── kawesome_generator.go
    └── kawesome_generator_test.go
```

## Build the Kusion Module

1. Clone this scaffold repository:
```
git clone https://github.com/KusionStack/kusion-module-scaffolding.git
```

2. Switch into the `src/` directory:
```
cd kusion-module-scaffolding/src
```

3. Build the Kusion Module using the `make` command:
```
make install
```

# Use the Kusion Module

After running `make install`, the newly-built module binary is now copied into a directory (`${KUSION_HOME}/modules/` by default) where `kusion` will look for during execution.

To declare a module as a dependency to your application configuration, add its reference to the `kcl.mod`:
```
[package]
name = "example"

[dependencies]
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.1.0" }
kawesome = { oci = "oci://ghcr.io/kusionstack/kawesome", tag = "0.1.0" }

[profile]
entries = ["main.k"]
```

:::tip

To understand more about what `kcl.mod` is, please refer to [this section](https://www.kusionstack.io/docs/kusion/configuration-walkthrough/overview#understanding-kclmod).
:::

Then in the KCL configuration code:
```
kawesome: ac.AppConfiguration {
    ...
    # Declare the kawesome module configurations. 
    accessories: {
        "kawesome": ks.Kawesome {
            service: ks.Service{
                port: 5678
            }
            randomPassword: ks.RandomPassword {
                length: 20
            }
        }
    }
}
```

For the complete configuration including the workspace configuration, please refer to the `example` directory.

# Develop a Kusion Module

To develop your own module, start with this scaffold and:
- Modify the `kcl.mod` to include the proper module name
- Modify the `kawesome.k` to include the proper KCL schema definition
- Modify the generator implementation including unit tests in the `src/` directory
- Modify the example in the `example` directory and `README.md` accordingly

For a more through walkthrough including how to publish a module, please refer to the [Kusion Module Developer Guide](https://www.kusionstack.io/docs/kusion/concepts/kusion-module/develop-guide).