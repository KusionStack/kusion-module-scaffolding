[package]
name = "kawesome-example"

[dependencies]
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.1.0" }
kawesome = { path = "../../kawesome/kawesome-schema" }

[profile]
entries = ["main.k"]
