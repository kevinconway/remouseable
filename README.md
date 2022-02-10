# reMouseable

> Use your reMarkable tablet as a mouse.

- [reMouseable](#remouseable)
  - [Project Status](#project-status)
    - [Alternative Projects](#alternative-projects)
  - [Overview](#overview)
  - [Code And Developer Documentation](#code-and-developer-documentation)
  - [Installation](#installation)
    - [Windows](#windows)
    - [OSX](#osx)
    - [Linux](#linux)
  - [Usage](#usage)
    - [reMarkable 2 Tablets](#remarkable-2-tablets)
    - [Wireless Tablet](#wireless-tablet)
    - [Advanced SSH Setup](#advanced-ssh-setup)
    - [All Options](#all-options)
  - [Common Issues And Solutions](#common-issues-and-solutions)
    - [OSX Privacy Settings](#osx-privacy-settings)
    - [Getting "panic: dial unix: missing address" On Windows](#getting-panic-dial-unix-missing-address-on-windows)
  - [Building](#building)
    - [Linux](#linux-1)
    - [OSX](#osx-1)
    - [Windows](#windows-1)
      - [Windows On Linux](#windows-on-linux)
  - [How It Works](#how-it-works)
  - [License](#license)
  - [Developing](#developing)
  - [Thanks](#thanks)

## Project Status

The project is stable and works for Remarkable 1 and Remarkable 2 tablets. I
will continue to have the latest release [available for
download](https://github.com/kevinconway/remouseable/releases).

Due to time constraints, I will no longer provide support through GitHub issues
or email. I will also no longer review pull requests that modify the project
code. If you want to add new features then please fork the project.

### Alternative Projects

If you are concerned that I'm no longer developing this project and want a more
active alternative then your best choice is
<https://github.com/Evidlo/remarkable_mouse>. It is actively developed and
supports more features than this project such as multi-monitor support.

If you are maintaining a fork of this project with new features and would like
to be mentioned here then make a PR to update this section of the README. Note
that my response to the PR will likely be slow.

## Overview

I'm a user of the [reMarkable](https://remarkable.com/) tablet. After using it
for a while I started wondering if it could be used as an input for my
computer so I could write and draw on digital whiteboards. It turns out, it can!

There's a great implementation of this feature written in Python at
<https://github.com/Evidlo/remarkable_mouse>. I'm working on this
implementation so that I can offer pre-built binaries that don't require a
specific language to be installed on the host machine.

## Code And Developer Documentation

This README contains how-to information for installing, configuration, and using
the project. To view the code API documentation check out the
[godocs](https://godoc.org/github.com/kevinconway/remouseable).

If you would like to modify the project or add a feature then see the technical
documentation in the `technical-documentation` directory.

## Installation

### Windows

Go to <https://github.com/kevinconway/remouseable/releases/latest> and download
the file named `windows.exe`. Then rename the file to `remouseable.exe`. You
can now open the Windows command prompt and start the program with:

```shell
cd Downloads
remouseable.exe
```

If a new version of the program comes out then you can overwrite your
`remouseable.exe` with a new version using exactly the same steps.

### OSX

Go to <https://github.com/kevinconway/remouseable/releases/latest> and download
the file named `osx`. Then rename the file to `remouseable`. Next, make the
program runnable with by opening a command line prompt and:

```shell
cd ~/Downloads
chmod +x remouseable
```

You can now run the program by opening a command line prompt and:

```shell
cd ~/Downloads
./remouseable
```

Note that the first time you run the application your system will prompt you
with a security notice. The remouseable application works by controlling your
mouse and OSX does not allow this by default. To enable the application you
must grant your command line prompt accessibility settings which allow it to
move the mouse. To do this, navigate to
`System Preferences -> Security & Privacy -> Privacy -> Accessibility`. You will
see your terminal or shell in the list of applications that have requested
accessibility permissions.

If you'd like to be able to launch the application through spotlight instead of
only the terminal then check out <https://github.com/isaacwisdom/reMouseableApp>
where another developer has created an Applscript wrapper that makes remouseable
act more like a typical OSX application.

### Linux

Go to <https://github.com/kevinconway/remouseable/releases/latest> and download
the file named `linux`. Then rename the file to `remouseable`. Next, make the
program runnable with by opening a command line prompt and:

```shell
cd ~/Downloads
chmod +x remouseable
```

You can now run the program by opening a command line prompt and:

```shell
cd ~/Downloads
./remouseable
```

## Usage

Most settings default to the correct values. The only value you should need to
set in the common case is the SSH password for the tablet. This password value
is found in the settings menu under `Help` and then `Copyrights and licenses`.
Your password will be near the bottom of the page. If you have an older tablet
that has not been updated to the latest software then your password may be
found in the `About` tab of the tablet menu at the bottom of the `General
Information` section. You may either give the password as text with

```bash
remouseable --ssh-password="XYZ123"
```

or you may choose to have a password prompt with:

```bash
remouseable --ssh-password="-"
```

Run one of these commands with your device connected over USB and your stylus
will become a mouse. The stylus is actually active _before_ it touches the
screen. This means you can see your mouse move by hovering the stylus just above
the writing surface but without directly touching the tablet. Once you touch the
tablet surface with the stylus the computer mouse will click and hold down the
left mouse button while you write or draw and then release the button when you
lift the stylus.

### reMarkable 2 Tablets

The application should work with both reMarkable and reMarkable 2 tablets.
However, the reMarkable 2 requires that you add
`--event-file /dev/input/event1` when executing because of a slight change in
where the stylus events are written in the new tablets. The full command should
look like
`remousable --ssh-password="MYPASSWORD" --event-file="/dev/input/event1"`.

### Wireless Tablet

The default expectation is that you will have your tablet connected over USB
which makes the default `10.11.99.1` address available. However, it is also
possible to access your device over wifi. If you attempt this method then you
will need to arrange for a static, or at least consistent, IP address for the
tablet. This is something you can usually do through configuring your router to
assign a fixed IP address to the device based on the hardware MAC address.

If you cannot assign the same `10.11.99.1` address in your setup then you may
override the default IP address when running the application:

### Advanced SSH Setup

By default, the tablet only accepts the root password for authentication. It is
possible, though, to install a custom public key on the device so that you can
use either password-less authentication or use a key pair that is encrypted with
the password of your choice rather than the device's default password.

If you'd like to create a key pair especially for accessing the reMarkable
tablet then start with a guide like
<https://help.github.com/en/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent>
that walks through creating a new key pair and registering it with your SSH
agent. For even more advanced SSH users, such as those using the gpg-agent as
the SSH agent, the remouseable application will talk to any valid SSH agent
implementation so long as the `SSH_AUTH_SOCK` value is set correctly.

Once you have a key pair ready, copy the public key value from `ssh-add -L` for
the key you want to use. Then copy the key over to your tablet with:

```bash
ssh root@10.11.99.1 # This will prompt for password.
mkdir -p ~/.ssh # This directory does not exist by default.
echo 'INSERT YOUR PUBLIC KEY HERE' >> ~/.ssh/authorized_keys
```

Now future connections over SSH will leverage your key pair and you can omit
the usual password flag when running the application.

Note that windows builds cannot use this option due to incompatibilities with
the current version of the windows ssh-agent.

Note that if you encounter the `Invalid MIT-MAGIC-COOKIE-1 key` error it means
that most likely the ssh fingerprint of the device might have changed to an
update of the tablet OS. Follow the ssh suggestion of removing the outdated
fingerprint then if you are satisfied that your device is indeed the right one
try connecting again.

```bash
remouseable --ssh-ip="192.168.1.110:22" # or other IP
```

### All Options

```
$ remouseable -h
Usage of remouseable:
      --debug-events             Stream hardware events from the tablet instead of acting as a mouse. This is for debugging.
      --disable-drag-event       Disable use of the custom OSX drag event. Only use this drawing on an Apple device is not working as expected.
      --event-file string        The path on the tablet from which to read evdev events. Probably don't change this. (default "/dev/input/event0")
      --orientation string       Orientation of the tablet. Choices are vertical, right, and left (default "right")
      --pressure-threshold int   Change the click detection sensitivity. 1000 is when the pen makes contact with the tablet. Set higher to require more pen pressure for a click. (default 1000)
      --screen-width int         Width of area confining the tablet pointer (defaults to full desktop width)
      --screen-height int        Height of area confining the tablet pointer (defaults to full desktop height)
      --screen-offset-x int      X offset of area confining the tablet pointer (defaults to 0)
      --screen-offset-y int      T offset of area confining the tablet pointer (defaults to 0)
      --ssh-ip string            The host and port of a tablet. (default "10.11.99.1:22")
      --ssh-password string      An optional password to use when ssh-ing into the tablet. Use - for a prompt rather than entering a value. If not given then public/private keypair authentication is used.
      --ssh-socket string        Path to the SSH auth socket. This must not be empty if using public/private keypair authentication.
      --ssh-user string          The ssh username to use when logging into the tablet. (default "root")
      --tablet-height int        The max units for the hight of the tablet. Probably don't change this. (default 15725)
      --tablet-width int         The max units for the width of the tablet. Probably don't change this. (default 20967)
pflag: help requested
exit status 2
```

## Common Issues And Solutions

### OSX Privacy Settings

If you are using this on an Apple or OSX device then you will need to give the
terminal or shell you are using permissions to control your mouse. Mouse
permissions are treated as an accessibility feature. If you are not prompted by
the operating system to update your permissions the first time you run the
application then you can navigate to
`System Preferences -> Security & Privacy -> Privacy -> Accessibility`. You will
see your terminal or shell in the list of applications that have requested
accessibility permissions.

### Getting "panic: dial unix: missing address" On Windows

This error message happens most often when the `--ssh-password` flag is missing
when running the application. On Windows, you must run the application with
either `remouseable.exe --ssh-password="MYPASSWORD"` or
`remouseable.exe --ssh-password="-"`.

## Building

There are pre-built binaries attached to each release that should work for all
64bit versions of linux, osx, and windows. However, if you prefer to generate
your own build then the following sections detail building a binary on
different platforms.

### Linux

Linux builds are dependent on:

- gcc
- x11 dev headers
- xtst dev headers
- xorg dev headers

These package will vary by name depending on your chosen linux distro. Debian
and Ubuntu users can install these with:

```shell
apt-get install -y gcc libc6-dev libx11-dev xorg-dev libxtst-dev
```

From there you run `make build`.

### OSX

OSX builds will require xcode and the xcode command line tools. These must be
installed through the Apple store.

Beyond xcode the build also requires installing support for gnu make if you want
to use the Makefile for generating a build. Homebrew users can install this
with:

```shell
brew install make coreutils findutils gnu-tar gnu-sed gawk gnutls gnu-indent gnu-getopt grep
export PATH="$(brew --prefix)/opt/make/libexec/gnubin:${PATH}"
```

From there you run `make build`.

### Windows

Windows builds require a GCC implementation. I recommend
<https://jmeubank.github.io/tdm-gcc/>. During installation you will be given the
option to add the GCC install to your path. If you choose not to then you will
need to temporarily add it to your path in PowerShell with:

```shell
$env:Path += ";C:\TDM-GCC-64\bin\"
```

The included Makefile contains too many bash specific commands to work in
PowerShell but you can still generate a binary by running:

```shell
go build main.go
```

#### Windows On Linux

If you want to generate a windows build from a linux machine then you will need
to install a MinGW implementation. Debian and Ubuntu users can do this with:

```shell
apt-get install -y gcc-multilib gcc-mingw-w64
```

The included Makefile does not have a build option for this but you can generate
the binary with:

```shell
CC=x86_64-w64-mingw32-gcc GOOS=windows go build main.go
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

The mouse interactions on the host are performed by using a modified version of
<https://github.com/go-vgo/robotgo>. The `pkg/internal/robotgo` directory
contains a stripped down version of `robotgo` that contains only the portions
required to detect the screen dimensions and send mouse events. The actual
`robotgo` project contains support for a much larger set of features such as
taking screen shots and controlling windows on the screen. However, each of
those additional features comes with additional system dependencies that make
creating a portable binary build difficult.

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

This project is go1.16+ compatible. A Makefile is included to make some things
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

I used the <https://github.com/gvalkov/golang-evdev> project as a reference when
implementing the `evdev` parser. I didn't use it directly because it is very
much oriented towards directly opening and managing a file descriptor for a
device. This project needs to read data from a remote device.

I used the <https://github.com/go-vgo/robotgo> project as the basis for
interacting with the operating system. I embedded portions of it here instead
of importing the Go package in order to limit the number of dependencies
required to build the project.
