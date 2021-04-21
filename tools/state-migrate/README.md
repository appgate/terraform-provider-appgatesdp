# about this tool


Prior to version 0.5.0, the project was named `terraform-provider-appgate-sdp`, and the provider was not yet published to registry.terraform.io.
When we published it, we noticed problems with using kebab-case in the name, which forced us to re-name the project to `terraform-provider-appgatesdp`


This tool will target a terraform plan directory and transform all appgate names found in .tf and .tfstate files to the new appgatesdp provider name. It creates a backup of the target directory <plan-directory>.backup as a sibling folder.


```sh
$ go run main.go migrate -dir /path/to/terraform-resources

```

or use the built binary
```sh
$ ./state-migrate migrate -dir /path/to/terraform-resources

```
