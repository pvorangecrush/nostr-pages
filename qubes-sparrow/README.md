# Setting up Sparrow on Qubes OS 4.1
In this tutorial, I will be walking you through setting up Sparrow Wallet in its own AppVM. We will be enhancing our security with a private electrum server connected over the TOR network.
You will need:
+ QubesOS 4.1 Installation (amd64)
+ (*optional*) Private Electrum server running on [Umbrel](https://umbrel.com), [RoninDojo](https://ronindojo.io)
  NOTE: You can use a public electrum server with potential loss of transction privacy

## Create and configure a new AppVM Qube
#### Open a dom0 terminal: ```[user@dom0 ~]$```
Create a new **sparrow** AppVM Qube
```bash
qvm-create sparrow -t debian-11 -l orange
```

Set the net-vm qube for our Sparrow AppVM to sys-whonix
```bash
qvm-prefs sparrow netvm sys-whonix
```

#### Create an allowance for our sparrow AppVM qube to bind ports on sys-whonix
Edit the network policy on dom0:
```bash
sudo nano /etc/qubes/policy.d/30-user-networking.policy
```

Add the following line to 30-user-networking.policy:
```nano
qubes.ConnectTCP * sparrow @default allow target=sys-whonix
```

## Launch a terminal in your sparrow AppVM qube:
Launch a sparrow AppVm Terminal from Qubes menu: ```[Qubes Launcher] > [Qube: sparrow] > [sparrow: Terminal]```

#### In terminal: ```[user@sparrow ~]$```
Setup qubes-bind-dirs:
```bash
sudo mkdir -p /rw/config/qubes-bind-dirs.d
sudo mkdir -p /rw/bind-dirs/opt/sparrow
sudo mkdir -p /rw/bind-dirs/usr/share/desktop-directories
```

Setup qubes-bind-dirs.d to bind directories on launch:
```bash
sudo nano /rw/config/qubes-bind-dirs.d/50_user.conf
```

Add the following lines 50_user.conf:
```nano
binds+=( '/opt/sparrow' )
binds+=( '/usr/share/desktop-directories' )
```

Setup /rw/config/rc.local to bind port 9050 on sys-whonix:
```bash
sudo nano /rw/config/rc.local
```

Add the following lines rc.local:
```nano
qvm-connect-tcp 9050:@default:9050
```
Shutdown your Sparrow AppVM qube using Qube Manager

## Download Sparrow Wallet on a disposible VM
Launch fedora-37-dvm (dvm): Firefox

Open [Sparrow Wallet download page](https://sparrowwallet.com/download)

Download 3 files:
+ sparrow_1.8.1-1_arm64.deb
+ sparrow-1.8.1-manifest.txt
+ sparrow-1.8.1-manifest.txt.asc

Copy files to our sparrow AppVm:
+ Open the downloads page on Firefox ```[Ctrl+Shift+Y]```
+ Click the folder icon
+ Select the three files, then right click and choose "Copy to other AppVm"
+ Choose sparrow as the target AppVm

## Verify and install Sparrow Wallet on sparrow AppVm
Launch a sparrow AppVm Terminal from Qubes menu: ```[Qubes Launcher] > [Qube: sparrow] > [sparrow: Terminal]```

#### In terminal: ```[user@sparrow ~]$```
Change directory to ~/QubesIncoming/disp*:
```bash
cd ~/QubesIncoming/disp*
```

Grab the key of the developer (Craig Raw):
```bash
curl https://keybase.io/craigraw/pgp_keys.asc | gpg --importsha256sum --check sparrow-1.8.1-manifest.txt --ignore-missing
```

Don't trust; verify the manifest:
```bash
gpg --verify sparrow-1.8.1-manifest.txt.asc
```

Now, verify the installation package:
```bash
sha256sum --check sparrow-1.8.1-manifest.txt --ignore-missing
```

Install the package
```bash
sudo apt install ./sparrow_1.8.1-1_amd64.deb
```

You can ignore this warning:
```N: Download is performed unsandboxed as root as file '/home/user/Downloads/sparrow_1.8.1-1_amd64.deb' couldn't be accessed by user '_apt'. - pkgAcquire::Run (13: Permission denied)```

## Final setup and configuration

#### Setup a Sparrow launch key-binding
+ [Qubes Launcher] > [System Tools] > [Keyboard]
+ Choose the [Application Shortcuts] Tab, then click the [+Add] button
+ Command: ```qvm-run -q -a sparrow /opt/sparrow/bin/Sparrow```
+ Enter Keyboard Shortcut: ```[Ctrl+Alt+S]```

#### Test your Sparrow Qube installation
**NOTE:** If you have a previous Sparrow Wallet installation, copy your existing .sparrow folder into sparrow AppVm /home/user/ directory
To Launch Sparrow Wallet: Press ```[Ctrl+Alt+S]```

#### Configure Sparrow to use sys-whonix torsocks proxy
+ Choose the Server tab (bottom):
+ Choose [Type:] Private Electrum
+ Enter your Umbrel or Dojo or other Electrum [URL:] ````**************.onion```
+ Open Sparrow Wallet preferences: Press ```[Ctrl+P]```
+ Toggle on [Use Proxy:]
+ Set [Proxy URL:] localhost   9050

## Final Thoughts:
The goal of this set up is to maximize security and privacy of our Sparrow Wallet setup. It is also possible to connect USB harwarde signers to the instance. Raise an issue [here](https://github.com/pvorangecrush/nostr-pages/) if you want me to add udev settings for usb-connected hardware wallets (I would rather you stay safe and airgap)
