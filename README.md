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

## Usage

`sandman <action> <container-name>`

### Build

The build command prepares a local image of the container using `buildah`, to be used with the run command.

### Start or Run

The start or run command spawns the container from the local image

### Test

Validates the connection to the Podman socket

### Sample

Outputs a scaffold of a container TOML file

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

# Running parameters
[Run]

# Allow x11 forwarding
x11 = true

# Allow Wayland forwarding
wayland = false

# Allow GPU acceleration
dri = false
gpu = false

# Allow host IPC
ipc = false

# Allow pulseaudio forwarding
pulseaudio = false

# Allow pipewire forwarding
pipewire = false

# Allow dbus forwarding
dbus = false

# Setup uids, requires /etc/subuid and /etc/subgid to be setup
uidmap = false

# Set to true if you need automatic data persistence. Will mount the /home inside the container to .local/share/sandman/xclock
home = false

# Deprecated option, implies in network = slirp4netns
net = false

# Any networking namespace modes allowed by podman-run
network = "none"

# If you want fonts to be mounted RO
fonts = true

# An optional name, if blank will use the default randomized name
name = "xclock"

# A list of usb device ids (vendor:product) to be mounted in the container
usbDevices = []

volumes = []
devices = []
env = []
ports = []

# If you need to share something from the host OS
# volumes = ['/etc/locale:/etc/locale:ro]

# If you need to access a device
# devices = ['/dev/video0']

# If you need special environment variables
# env = ['ENV=test']

# If you need to expose ports
# ports = ['8080:8080']

[Limits]
# Still being implemented
```

Build it with `sandman build xclock`

Spawn it with `sandman run xclock`

For full reference see `config/config.go`
