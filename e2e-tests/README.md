## Prerequisites

* Vagrant
* VirtualBox (or another provider like VMware or libvirt, but the Vagrantfile will need modifications)



## Possible Errors You May See

When running on macOS with VirtualBox, you may occasionally see an error saying:

```
Received unexpected error:
exit status 1
```

Open up VirtualBox and attempt to start the Vagrant VM by hand. If you see an error about a kernel mode driver not being installed, it's likely because you've updated the OS since the time you installed VirtualBox. To fix it, just reinstall VirtualBox. Be sure to reboot after installation.