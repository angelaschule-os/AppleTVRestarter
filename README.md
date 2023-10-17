# AppleTVRestarter


## Build

You need at least [Go](https://go.dev/) 1.21.3 to build this.

```shell
git clone https://github.com/angelaschule-os/AppleTVRestarter
cd AppleTVRestarter
go build
```

## Usage

Create a `.env` file with the following content:

```env
BASE_URL      = "https:///{yourDomain}.jamfcloud.com/api"
NETWORKT_ID   = "XXXXXXXX"
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

## Links

- [Jamf School API](https://school.jamfcloud.com/api/docs/)
