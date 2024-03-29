[package]
name = "example"

[dependencies]
kawesome = { oci = "oci://ghcr.io/kusionstack/kawesome", tag = "0.1.0" }
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.1.0" }

[profile]
entries = ["main.k"]
