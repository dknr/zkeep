# zkeep

ZFS snapshot management client and daemon.

## Overview

zkeep provides a simple CLI interface for managing ZFS snapshots through a
background daemon. The daemon runs as root and communicates with the client
via a Unix domain socket.

## Installation

```bash
make install
rc-update add zkeepd default
rc-service zkeepd start
```

**Note**: Only an OpenRC init script is provided. You will need to create a
systemd unit or other init script if using a different init system.

**Prerequisites**:
1. Create the `zkeep` group: `groupadd zkeep`
2. Add your user to the group: `usermod -aG zkeep <username>`
3. The daemon must run as root

## Usage

```bash
# Create a snapshot
zkeep snapshot tank/dataset@daily

# List snapshots
zkeep list tank/dataset

# Destroy a snapshot
zkeep destroy tank/dataset@daily

# Rollback to a snapshot
zkeep rollback tank/dataset@daily

# Check daemon status
zkeep status

# Show version information
zkeep version
```

## Architecture

- **zkeepd** (daemon): Runs as root, manages ZFS operations
- **zkeep** (client): CLI tool, communicates with daemon via Unix socket
- **Socket**: `/var/run/zkeep.sock` (owned by `root:zkeep`)

## License

ISC License - see LICENSE file for details.
