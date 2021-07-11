# Sandman

Sandboxes with Podman

This is a work in progress.

## Usage

`sandman <action> <container-name>` where action is `build` or `run`

## Example configuration file

```toml
[build]
image_name = 'xclock'
instructions = '''
FROM archlinux
RUN pacman -Syu xorg-xclock --noconfirm
CMD "/usr/bin/xclock"
'''

[run]
x11 = true
dri = true
ipc = true
pulseaudio = true
dbus = true
net = true
volumes = ['/home/julio/Containers/xclock:/root']
devices = ['/dev/tun']

[env]
TEST = "value"
```
