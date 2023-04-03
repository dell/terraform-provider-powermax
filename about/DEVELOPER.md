## Developing the provider plugin

### Prerequisites
- Go 1.18+ installed and configured.
- Terraform v1.0.3+ installed locally.


### Prepare env for local provider
1. Check and set GOBIN where Go installs your binaries
    ```
    go env GOBIN
    ```
    if empty, set GOBIN to default, for example
    ```
    export GOBIN=/usr/lib/go/bin
    ```

2. Create a new file called `.terraformrc` in your home directory (~), then add the `dev_overrides` block using the provider address and GOBIN path
    ```
    provider_installation {
        dev_overrides {
            "registry.terraform.io/dell/powermax" = "/usr/lib/go/bin"
        }
    
        # For all other providers, install them directly from their origin provider
        # registries as normal. If you omit this, Terraform will _only_ use
        # the dev_overrides block, and so no other providers will be available.
        direct {}
    }
    ```


### Build the provider
1. Clone the repo and switch to the develop branch
    ```
    git clone git@eos2git.cec.lab.emc.com:Jason-jin1/terraform-provider-powermax-framework.git
    git checkout -b <branch-name>
    ```

2. Enter the repo
    ```
    cd terraform-provider-powermax-framework
    ```

3. Build the provider
    ```
    go install
    ```

### Run and test
1. Navigate to the path where tf file locates, for example
    ```
    cd examples
    ```
2. Change the tf configuration and run a Terraform plan. It will report a warning that overrides are in effect, and you are good to go!
    ```
    terraform plan
    ```


## Debugging with GoLand
1. Navigate to `main.go` file under project directory, set `debug` value to true
    ```
    flag.BoolVar(&debug, "debug", true, "set to true to run the provider with support for debuggers like delve")
    ```

2. Execute the debug command

3. Copy&Paste the output command to set `TF_REATTACH_PROVIDERS` in your terminal(Note: the command will change everytime you rerun)

4. Add your breakpoint and navigate to the tffile and run with `terraform plan`
   
