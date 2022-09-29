# Sandman

Sandboxes with Podman

**Notice:** This is a work in progress.

The aim of Sandman is to provide a simple command to build and run Linux applications within a containerized rootless sandbox using Buildah and Podman. 

Each application has its own TOML configuration file with the instructions to build its base image and running options, where you can specify Dockerfile equivalent build instructions and what do you want enabled inside the sandbox. Bear in mind that the more toggles you enable, the less constrained is the sandbox.

Due to having a base image already built, the applications run in ephemeral containers. That means that each time you launch a sandbox, it's brand new. If you need data persistence, you can export a volume inside the container where you need the data to persist.

This also means that you have full control of the sandboxing, you choose what do you want exposed into the container, differently from other solutions where the app declares its own sandbox boundaries.

## Requirements

It is required to have `podman` and `buildah` installed and configured to work with rootless containers.

Sandman uses the podman socket to send binded commands, you need to ensure it is active (e.g. `systemctl --user enable --now podman.socket`)

In order to use `uidmap` you need to setup subuid and subgid.

More information can be found at: https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md

## Files

The TOML files with the container configuration are stored in `.config/sandman` inside your home.

The optional local storage is stored in `.local/share/sandman` inside your home.

## Installing

A Makefile is provided with basic commands. You can use `make all` to download dependencies, test everything, compile and install.

## Usage

`sandman <action> <container-name>`

### Build

The build command prepares a local image of the container using `buildah`, to be used with the run command.

### Start or Run

The start or run command spawns the container from the local image. Starting spawns a dettached container, Run will auto-attach.

### Test

Validates the connection to the Podman socket

### Sample

Outputs a simplified scaffold of a container TOML file. See full example below

### Help

Provides help for any subcommand

## Example configuration file

**~/.config/sandman/xclock.toml**

```toml
[Build]

# Directory where files will be pulled from
ContextDirectory = ""

# Build instructions, conforms to Dockerfile syntax
Instructions = '''
FROM archlinux
RUN pacman -Syu xorg-xclock --noconfirm
CMD "/usr/bin/xclock"
'''

# Any additional image names
AdditionalImageNames = ["sandman/xclock:2.0.alpha"]

[Build.Limits]
# Change limits for build
Ulimit = ["nofile=4096"]

# Running parameters
[Run]

# Allow x11 forwarding
X11 = true

# Allow Wayland forwarding
Wayland = false

# Allow GPU acceleration
Dri = false
Gpu = false

# Allow host IPC
Ipc = false

# Allow pulseaudio forwarding
Pulseaudio = false

# Allow pipewire forwarding
Pipewire = false

# Allow dbus forwarding
Dbus = false

# Setup uids, requires /etc/subuid and /etc/subgid to be setup
Uidmap = false

# Set to true if you need automatic data persistence. Will mount the /home inside the container to .local/share/sandman/xclock
Home = false

# Deprecated option, implies in network = slirp4netns
Net = false

# Any networking namespace modes allowed by podman-run
Network = "none"

# If you want fonts to be mounted RO
Fonts = true

# An optional name, if blank will use the default randomized name
Name = "xclock"

# A list of usb device ids (vendor:product) to be mounted in the container
UsbDevices = []

Volumes = []
Devices = []
Env = []
Ports = []

# If you need to share something from the host OS
# Volumes = ['/etc/locale:/etc/locale:ro']

# If you need to access a device
# Devices = ['/dev/video0']

# If you need special environment variables
# Env = ['ENV=test']

# If you need to expose ports
# Ports = ['8080:8080']

[Run.Limits]
# Still being implemented. See full list in config/config.go
```

Build it with `sandman build xclock`

Spawn it with `sandman start xclock`

For full reference see `config/config.go`
