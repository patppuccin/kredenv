# Prompt Frameworks

kredenv exposes environment variables that prompt frameworks can use to display when secrets are loaded in the current shell session.

## Available Variables

| Variable               | Description                                 |
| ---------------------- | ------------------------------------------- |
| `KREDENV_LOADED_COUNT` | Number of secrets currently loaded          |
| `KREDENV_LOADED_VARS`  | Comma-separated list of loaded secret names |

Any prompt framework that supports environment variable modules can integrate with kredenv using these variables.

## Starship

[Starship](https://starship.rs) is a cross-shell prompt with an `env_var` module that reads environment variables and displays them in the prompt.

Add the following to your `~/.config/starship.toml`:

```toml
format = """
$username\
at $hostname \
in $directory \
(on $git_branch )\
(with ${env_var.kredenv})\
$line_break\
$character\
"""

[env_var.kredenv]
variable = "KREDENV_LOADED_COUNT"
format = "[$symbol $env_value secrets]($style) "
symbol = "🔑"
style = "bold yellow"
disabled = false
```

This displays a key emoji and the number of loaded secrets whenever a `kredsfile.yaml` is in scope with `autoload: true`.

**Example prompt output:**

```
patppuccin at athena in ~/projects/myapp on main with 🔑 3 secrets
❯
```

### Showing loaded variable names

To display the names of loaded secrets instead of just the count:

```toml
[env_var.kredenv]
variable = "KREDENV_LOADED_VARS"
format = "[$symbol $env_value]($style) "
symbol = "🔑"
style = "bold yellow"
disabled = false
```

This shows the comma-separated list of secret names: `🔑 AWS_ACCESS_KEY_ID,DATABASE_PASSWORD`.
