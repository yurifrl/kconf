# Kconf

Kubernetes configurations manager, merges kubernetes configurations from a file into your `~/.kube/config`

## Usage

- Download a release from `https://github.com/yurifrl/kconf/releases`
  - Unpack the tar.gz of your architeture and and the binary to your path
- This cli reads files from a folder (by default `$HOME/.kconf/configs`) and concatenate then into your `$HOME/.kube/config`
- The files inside inside `$HOME/.kconf/configs` are a bunch of kubernetes config files (`kubectl config view`)
```
apiVersion: v1
kind: Config
clusters:
- name: "<name>"
  cluster: {}

users:
- name: "<user-name>"
  user: {}

contexts:
- name: "<context>"
  context: {}

current-context: "<context>"
```
- If you have two of does files, after running `kconf` you will have one inside `$HOME/.kube/config` with both off them concatenated together

## Build

- `go build -o dist/kconf .`

## TODO

- [ ] Write tests
- [ ] Deploy releases with goreleaser
- [ ] CI/CD
- [ ] Create command `use` that accepts a config and use that config
- [ ] Make the cli backup the ~/.kube/config file
- [ ] Ship on Arch
