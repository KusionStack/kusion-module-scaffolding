[package]
name = "example"

[dependencies]
kam = { git = "https://github.com/KusionStack/kam.git", tag = "v0.1.0-beta" }
kawesome = { path = "../../v1" }
[profile]
entries = ["main.k"]

