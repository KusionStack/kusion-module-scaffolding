[package]
name = "example"

[dependencies]
kawesome = { oci = "oci://ghcr.io/kusionstack/kawesome", tag = "0.2.0" }
service = {oci = "oci://ghcr.io/kusionstack/service", tag = "0.1.0" }
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.2.0" }

[profile]
entries = ["main.k"]
