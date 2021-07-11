# Sandman

Sandboxes with Podman

**Notice:** This is a work in progress. This is the first time I'm developing in Rust, so the code will not be in the best shape it could be.

The aim of Sandman is to provide a simple command to build and run Linux applications within a containerized rootless sandbox using Buildah and Podman. 

Each application has its own TOML configuration file with the instructions to build its base image and running options, where you can specify Dockerfile equivalent build instructions and what do you want enabled inside the sandbox. Bear in mind that the more toggles you enable, the less constrained is the sandbox.

Due to having a base image already built, the applications can run in ephemeral containers. That means that each time you launch a sandbox, it's brand new. If you need data persistence, you can export a volume inside the container where you need the data to persist.

## Usage

`sandman <action> <container-name>` where action is `build` or `run`

## Example configuration file

```toml
[build]
image_name = 'sandman/xclock'
instructions = '''
FROM archlinux
RUN pacman -Syu xorg-xclock --noconfirm
CMD "/usr/bin/xclock"
'''

[run]
x11 = true
dri = false
ipc = false
pulseaudio = false
dbus = false
net = false
volumes = []
devices = []

# If you need data persistence
# volumes = ['/home/julio/Containers/xclock:/root']

# If you need to access a device
# devices = ['/dev/sdd']

[env]
TEST = "value"
```

Build it with `sandman build xclock`
Run it with `sandman run xclock`
