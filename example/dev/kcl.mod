[package]
name = "example"

[dependencies]
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.2.0-beta" }
service = { oci = "oci://ghcr.io/kusionstack/service", tag = "0.1.0-beta" }
monitoring = { oci = "oci://ghcr.io/kusionstack/monitoring", tag = "0.1.0-beta-test" }
kawesome = { oci = "oci://ghcr.io/kusionstack/kawesome", tag = "0.1.0-beta" }
[profile]
entries = ["main.k"]

