# Kconf

Kubernetes config files are simple enought if you have only one context, when the number of contexts goes up so it's the difficulty to manage then.

Kconf let you have a bunch of config files, with only one context and then squish then for you, so you don't have to edit a ever more larger kubernetes config

Kconf takes files like this:
```
# $HOME/.helm/configs/context-1
apiVersion: v1
kind: Config
clusters:
- name: cluster-1
  cluster: {}
users:
- name: user-1
  user: {}
contexts:
- name: context-1
  context:
    user: user-1
    cluster: cluster-1
current-context: context-1
```
```
# $HOME/.helm/configs/context-2
apiVersion: v1
kind: Config
clusters:
- name: cluster-2
  cluster: {}
users:
- name: user-2
  user: {}
contexts:
- name: context-2
  context:
    user: user-2
    cluster: cluster-2
current-context: context-2
```
And make then into this:
```
# $HOME/.helm/config
apiVersion: v1
kind: Config
clusters:
- name: cluster-1
  cluster: {}
- name: cluster-2
  cluster: {}
users:
- name: user-1
  user: {}
- name: user-2
  user: {}
contexts:
- name: context-1
  context:
    user: user-1
    cluster: cluster-1
- name: context-2
  context:
    user: user-2
    cluster: cluster-2
current-context: context-1
```

## Usage

- Download a release from `https://github.com/yurifrl/kconf/releases`
  - Unpack the tar.gz of your architeture and and the binary to your path
- This cli reads files from a folder (by default `$HOME/.kconf/configs`) and concatenate then into your `$HOME/.kube/config`
- The files inside `$HOME/.kconf/configs` are a bunch of kubernetes config files (`kubectl config view`)
- If you have two of does files, after running `kconf` you will have one on `$HOME/.kube/config` with both off them concatenated in one

## Build

- `go build -o dist/kconf .`

## TODO

- [ ] Write tests
- [ ] Deploy releases with goreleaser
- [ ] CI/CD
- [ ] Create command `use` that accepts a config and use that config
- [ ] Make the cli backup the ~/.kube/config file
- [ ] Ship on Arch
