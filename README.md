# nsoctl
Small cli tool to help with day to day Cisco NSO operations

## Config
The default config file is located at `~/.nsoctl.toml`.<br>
The path can be changed using the the command: `nsoctl --config /new_path/.nsoctl.toml`

The config file should contain the following:
````toml
[nso]
insecureSkipVerify = false   # true/false depending on where or not the client trusts the ssl certificate of server
serverFqdn = ""              # server fqdn
username = ""                # username
password = ""                # password
````
