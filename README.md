# AppleTVRestarter


## Build

You need at least [Go](https://go.dev/) 1.21.3 to build this.

```shell
git clone https://github.com/angelaschule-os/AppleTVRestarter
cd AppleTVRestarter
go build -ldflags "-X main.GitCommit=$(git rev-parse --short HEAD)"

```

## Usage

Create a `.env` file with the following content:

```env
BASE_URL      = "https:///{yourDomain}.jamfcloud.com/api"
NETWORK_ID    = "XXXXXXXX"
KEY           = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
```

```shell
Usage:
  AppleTVRestarter [--refresh] [--restart]
Options:
  --refresh	Send refresh command to devices
  --restart	Send restart command to devices
  --help	Display this help message
```

## Create a release on Github


Creating an annotated tag in Git and share it.
```shell
git tag -a v1.0 -m "Version 1.0"
git push origin v1.0

```

By default, the git push command doesn’t transfer tags to remote servers. You
will have to explicitly push tags to a shared server after you have created
them.


## Links

- [Jamf School API](https://school.jamfcloud.com/api/docs/)
