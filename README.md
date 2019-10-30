vrdwallet
=========

vrdwallet is a daemon handling soter wallet functionality for a
single user.  It acts as both an RPC client to soterd and an RPC server
for wallet clients and legacy RPC applications.

Public and private keys are derived using the hierarchical
deterministic format described by




vrdwallet is not an SPV client and requires connecting to a local or
remote soterd instance for asynchronous blockchain queries and
notifications over websockets.  Full soterd installation instructions
can be found [here](https://github.com/soteria-dag/soterd).  An alternative
SPV mode that is compatible with soterd and Bitcoin Core is planned for
a future release.

Wallet clients can use one of two RPC servers:

  1. A legacy JSON-RPC server mostly compatible with Bitcoin Core

     The JSON-RPC server exists to ease the migration of wallet applications
     from Core, but complete compatibility is not guaranteed.  Some portions of
     the API (and especially accounts) have to work differently due to other
     design decisions (mostly due to BIP0044).  However, if you find a
     compatibility issue and feel that it could be reasonably supported, please
     report an issue.  This server is enabled by default.

  2. An experimental gRPC server

     The gRPC server uses a new API built for vrdwallet, but the API is not
     stabilized and the server is feature gated behind a config option
     (`--experimentalrpclisten`).  If you don't mind applications breaking due
     to API changes, don't want to deal with issues of the legacy API, or need
     notifications for changes to the wallet, this is the RPC server to use.
     The gRPC server is documented [here](./rpc/documentation/README.md).

## Installation and updating

### Windows - MSIs Available

Install the latest MSIs available here:

https://github.com/soteria-dag/soterd/releases

https://github.com/soteria-dag/vrdwallet/releases

### Windows/Linux/BSD/POSIX - Build from source

Building or updating from source requires the following build dependencies:

- **Go 1.5 or 1.6**

  Installation instructions can be found here: http://golang.org/doc/install.
  It is recommended to add `$GOPATH/bin` to your `PATH` at this point.

  **Note:** If you are using Go 1.5, you must manually enable the vendor
    experiment by setting the `GO15VENDOREXPERIMENT` environment variable to
    `1`.  This step is not required for Go 1.6.

- **Glide**

  Glide is used to manage project dependencies and provide reproducible builds.
  To install:

  `go get -u github.com/Masterminds/glide`

Unfortunately, the use of `glide` prevents a handy tool such as `go get` from
automatically downloading, building, and installing the source in a single
command.  Instead, the latest project and dependency sources must be first
obtained manually with `git` and `glide`, and then `go` is used to build and
install the project.

**Getting the source**:

For a first time installation, the project and dependency sources can be
obtained manually with `git` and `glide` (create directories as needed):

```
git clone https://github.com/soteria-dag/vrdwallet $GOPATH/src/github.com/soteria-dag/vrdwallet
cd $GOPATH/src/github.com/soteria-dag/vrdwallet
glide install
```

To update an existing source tree, pull the latest changes and install the
matching dependencies:

```
cd $GOPATH/src/github.com/soteria-dag/vrdwallet
git pull
glide install
```

**Building/Installing**:

The `go` tool is used to build or install (to `GOPATH`) the project.  Some
example build instructions are provided below (all must run from the `vrdwallet`
project directory).

To build and install `vrdwallet` and all helper commands (in the `cmd`
directory) to `$GOPATH/bin/`, as well as installing all compiled packages to
`$GOPATH/pkg/` (**use this if you are unsure which command to run**):

```
go install . ./cmd/...
```

To build a `vrdwallet` executable and install it to `$GOPATH/bin/`:

```
go install
```

To build a `vrdwallet` executable and place it in the current directory:

```
go build
```

## Getting Started

The following instructions detail how to get started with vrdwallet connecting
to a localhost soterd.  Commands should be run in `cmd.exe` or PowerShell on
Windows, or any terminal emulator on *nix.

- Run the following command to start soterd:

```
soterd -u rpcuser -P rpcpass
```

- Run the following command to create a wallet:

```
vrdwallet -u rpcuser -P rpcpass --create
```

- Run the following command to start vrdwallet:

```
vrdwallet -u rpcuser -P rpcpass
```

If everything appears to be working, it is recommended at this point to
copy the sample soterd and vrdwallet configurations and update with your
RPC username and password.

PowerShell (Installed from MSI):
```
PS> cp "$env:ProgramFiles\Soterd Suite\Soterd\sample-soterd.conf" $env:LOCALAPPDATA\Soterd\soterd.conf
PS> cp "$env:ProgramFiles\Soterd Suite\vrdwallet\sample-vrdwallet.conf" $env:LOCALAPPDATA\vrdwallet\vrdwallet.conf
PS> $editor $env:LOCALAPPDATA\Soterd\soterd.conf
PS> $editor $env:LOCALAPPDATA\vrdwallet\vrdwallet.conf
```

PowerShell (Installed from source):
```
PS> cp $env:GOPATH\src\github.com\soteria-dag\soterd\sample-soterd.conf $env:LOCALAPPDATA\Soterd\soterd.conf
PS> cp $env:GOPATH\src\github.com\soteria-dag\vrdwallet\sample-vrdwallet.conf $env:LOCALAPPDATA\vrdwallet\vrdwallet.conf
PS> $editor $env:LOCALAPPDATA\Soterd\soterd.conf
PS> $editor $env:LOCALAPPDATA\vrdwallet\vrdwallet.conf
```

Linux/BSD/POSIX (Installed from source):
```bash
$ cp $GOPATH/src/github.com/soteria-dag/soterd/sample-soterd.conf ~/.soterd/soterd.conf
$ cp $GOPATH/src/github.com/soteria-dag/vrdwallet/sample-vrdwallet.conf ~/.vrdwallet/vrdwallet.conf
$ $EDITOR ~/.soterd/soterd.conf
$ $EDITOR ~/.vrdwallet/vrdwallet.conf
```

## Issue Tracker

The [integrated github issue tracker](https://github.com/soteria-dag/vrdwallet/issues)
is used for this project.

## GPG Verification Key

All official release tags are signed by Conformal so users can ensure the code
has not been tampered with and is coming from the btcsuite developers.  To
verify the signature perform the following:

- Download the public key from the Conformal website at
  https://opensource.conformal.com/GIT-GPG-KEY-conformal.txt

- Import the public key into your GPG keyring:
  ```bash
  gpg --import GIT-GPG-KEY-conformal.txt
  ```

- Verify the release tag with the following command where `TAG_NAME` is a
  placeholder for the specific tag:
  ```bash
  git tag -v TAG_NAME
  ```

## License

vrdwallet is licensed under the liberal ISC License.
