# reMouseable

> Use your reMarkable tablet as a mouse.

## Overview

I'm a user of the [reMarkable](https://remarkable.com/) tablet. After using it
for a while I started wondering if it could be used as an input for my
computer so I could write and draw on digital whiteboards. It turns out, it can!

There's a great implementation of this feature written in Python at
https://github.com/Evidlo/remarkable_mouse. I'm working on this implementation
so that I can eventually offer pre-built binaries that don't require a
specific language to be installed on the host machine.

## Docs

This README contains how-to information for installing, configuration, and using
the project. To view the code API documentation check out the
[godocs](https://godoc.org/github.com/kevinconway/remouseable).

## Setup And Configuration

The easiest way to get started is to download the pre-compiled binary from the
latest release at https://github.com/kevinconway/remouseable/releases/latest. The
files are named after the OS for which they are built. All builds are currently
64 bit. If you need to build for a different environment see the `Building`
section for instructions.

Most settings default to the correct values. The only value you should need to
set in the common case is the SSH password for the tablet. This password value
is found in the `About` tab of the tablet menu at the bottom of the
`General Information` section. You may either give the password as text with

```bash
remouseable --ssh-password="XYZ123"
```

or you may choose to have a password prompt with:

```bash
remouseable --ssh-password="-"
```

Run one of these commands with your device connected over USB and your stylus
will become a mouse. The stylus is actually active _before_ it touches the
screen which means you can see your mouse move without directly touching the
tablet. Once you touch the tablet surface with the stylus the computer mouse
will click and hold down the left mouse button while you write or draw and then
release the button when you lift the stylus.

### Easier SSH Setup

By default, the tablet only accepts the root password for authentication. It is
possible, though, to install a custom public key on the device so that you can
use either password-less authentication or use a key pair that is encrypted with
the password of your choice rather than the device's default password.

If you'd like to create a key pair especially for accessing the reMarkable
tablet then start with a guide like https://help.github.com/en/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent that walks through
creating a new key pair and registering it with your SSH agent. For advanced
SSH users, such as those using the gpg-agent as the SSH agent, the reMouse
application will talk to any valid SSH agent implementation so long as the
`SSH_AUTH_SOCK` value is set correctly.

Once you have a key pair ready, copy the public key value from `ssh-add -L` for
the key you want to use. Then copy the key over to your tablet with:

```bash
ssh root@10.11.99.1 # This will prompt for password.
mkdir -p ~/.ssh # This directory does not exist by default.
echo 'INSERT YOUR PUBLIC KEY HERE' >> ~/.ssh/authorized_keys
```

Now future connections over SSH will leverage your key pair and you can omit
the usual password flag when running the application.

### Wireless Tablet

The default expectation is that you will have your tablet connected over USB
which makes the default `10.11.99.1` address available. However, it is also
possible to access your device over wifi. If you attempt this method then you
will need to arrange for a static, or at least consistent, IP address for the
tablet. This is something you can usually do through configuring your router to
assign a fixed IP address to the device based on the hardware MAC address.

If you cannot assign the same `10.11.99.1` address in your setup then you may
override the default IP address when running the application:

```bash
remouseable --ssh-ip="192.168.1.110" # or other IP
```

### All Options

```
$ remouseable -h
Usage of remouseable:
      --event-file string     The path on the tablet from which to read evdev events. Probably don't change this. (default "/dev/input/event0")
      --orientation string    Orientation of the tablet. Choices are vertical, right, and left (default "right")
      --screen-height int     The max units per millimeter of the host screen height. Probably don't change this. (default 1080)
      --screen-width int      The max units per millimeter of the host screen width. Probably don't change this. (default 1920)
      --ssh-ip string         The host and port of a tablet. (default "10.11.99.1:22")
      --ssh-password string   An optional password to use when ssh-ing into the tablet. Use - for a prompt rather than entering a value. If not given then public/private keypair authentication is used.
      --ssh-socket string     Path to the SSH auth socket. This must not be empty if using public/private keypair authentication. (default "/run/user/1000/gnupg/S.gpg-agent.ssh")
      --ssh-user string       The ssh username to use when logging into the tablet. (default "root")
      --tablet-height int     The max units per millimeter for the hight of the tablet. Probably don't change this. (default 15725)
      --tablet-width int      The max units per millimeter for the width of the tablet. Probably don't change this. (default 20967)
pflag: help requested
exit status 2
```

## Building

If you want to build this for a specific environment then first make sure you
have all the system dependencies required by https://github.com/go-vgo/robotgo
which is the default driver used to control the mouse. After that, it should be
a matter of:

```bash
go get
make build # alteratively, go build main.go if you want to pass custom options
```

## How It Works

The project is implemented as a set of successive layers that turn the tablet
into a mouse. It follows as:

-   SSH into the device and start streaming `evdev` data back to the host.
-   Convert the raw byte stream into structured `evdev` data containers.
-   Feed all events into a state machine that emits higher level state change
    events like "CLICK" and "MOVE".
-   Use state change events as a trigger for moving or clicking the mouse
    on the host machine.

Each of these layers has an interface defined in the `pkg/domain.go` file.
The mouse interactions on the host are performed by using
https://github.com/go-vgo/robotgo.

## License

    remouseable is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License version 3 as published
    by the Free Software Foundation.

    remouseable is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with remouseable.  If not, see <https://www.gnu.org/licenses/>

## Developing

This project is go1.13+ compatible. A Makefile is included to make some things
easier. Some make targets of note:

-   make generate

    Re-generate any automatically generated code. Note that there is a gomock
    bug making it necessary to manually modify the files after generation
    because it adds a cyclical import.

-   make test

    Run all the unit tests and generate a coverage report in `.coverage/`.

-   make lint

    Run the golangci-lint suite using the included configuration.

-   make fmt

    Apply `goimports` formatting.

-   make build

    Generate a binary from the current project state.

-   make tools

    Generate a `.bin/` directory that contains a built version of each of the
    tools used to build and test the project.

-   make update / make updatetools

    Run `go get -u` for the project or for the project tooling.

-   make clean / make cleantools / make cleancoverage

    Remove files generated by the Makefile. The top-level `clean` should remove
    all artifacts such as `./bin` and `./coverage`. The other are scoped to
    specific artifacts for cases where, for example, you want to remove old
    coverage reports and regenerate them.

## Thanks

I used the https://github.com/gvalkov/golang-evdev project as a reference when
implementing the `evdev` parser. I didn't use it directly because it is very
much oriented towards directly opening and managing a file descriptor for a
device. This project needs to read data from a remote device.

## Future Features

A lot of the low level device interfaces like `evdev` and `hid` are new to me
and I'm still learning them. Getting solid OSX support is top of my list to
complete. Support for multiple monitors would be next. After that, figuring out
how to handle the tilt and pressure events in a way that enable full graphics
tablet functionality on the host would be next.
