### Testing
Test it: `go install ./...`

* `cobra_tofu -h`
* `cobra_tofu -concise --no-color init -backend`
* `source <(cobra_tofu completion zsh)`
  * test commpletion with tab. Works for flags too
  * good example is to test it for `cobra_tofu workspace <TAB>` where it completes the command and for the `select` one we offer suggestions about the existing workspaces
* `cobra_tofu -chdir commands init` - should succeed
* `cobra_tofu init -chdir commands` - should fail because -chdir is allowed only before the subcommmand

