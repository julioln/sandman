# Sandman

Sandboxes with Podman

**Notice:** This is a work in progress. This is the first time I'm developing in Rust, so the code will not be in the best shape it could be.

The aim of Sandman is to provide a simple command to build and run Linux applications within a containerized rootless sandbox using Buildah and Podman. 

Each application has its own TOML configuration file with the instructions to build its base image and running options, where you can specify Dockerfile equivalent build instructions and what do you want enabled inside the sandbox. Bear in mind that the more toggles you enable, the less constrained is the sandbox.

Due to having a base image already built, the applications run in ephemeral containers. That means that each time you launch a sandbox, it's brand new. If you need data persistence, you can export a volume inside the container where you need the data to persist.

This also means that you have full control of the sandboxing, you choose what do you want exposed into the container, differently from other solutions where the app declares its own sandbox boundaries.

## Requirements

It is required to have `podman` and `buildah` installed and configured to work with rootless containers.

In order to use `uidmap` you need to setup subuid and subgid.

More information can be found at: https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md

## Files

The TOML files with the container configuration are stored in `.config/sandman` inside your home.

The optional local storage is stored in `.local/share/sandman` inside your home.

## Usage

`sandman <action> <container-name>` where action is `build`, `run` or `args`.

### Build

The build command prepares a local image of the container using `buildah`, to be used with the run command.

### Run

The run command spawns the container from the local image, compiling a series of parameters to `podman` based on the TOML configuration.

### Args

The args command shows what parameters are being passed to the `podman` command without executing it.

## Example configuration file (aka Hello World)

**~/.config/sandman/xclock.toml**

```toml
# Build instructions, conforms to Dockerfile syntax
[build]
instructions = '''
FROM archlinux
RUN pacman -Syu xorg-xclock --noconfirm
CMD "/usr/bin/xclock"
'''

# Running parameters
[run]

# Allow x11 forwarding
x11 = true

# Allow Wayland forwarding
wayland = false

# Allow GPU acceleration
dri = false

# Allow host IPC
ipc = false

# Allow pulseaudio forwarding
pulseaudio = false

# Allow dbus forwarding
dbus = false

# Setup uids, requires /etc/subuid and /etc/subgid to be setup
uidmap = false

# Set to true if you need automatic data persistence. Will mount the /home inside the container to .local/share/sandman/xclock
home = false

# Deprecated option, implies in network = slirp4netns
net = false

# Any networking modes allowed by podman-run
network = "none"

volumes = []
devices = []
env = []
ports = []

# If you need to share something from the host OS
# volumes = ['/usr/share/fonts:/usr/share/fonts:ro]

# If you need to access a device
# devices = ['/dev/sdd', '/dev/video0']

# If you need special environment variables
# env = ['ENV=test']

# If you need to expose ports
# ports = ['8080:8080']
```

Build it with `sandman build xclock`

Run it with `sandman run xclock`
