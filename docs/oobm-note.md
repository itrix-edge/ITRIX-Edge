OOBM Hardware Management
========================

## OOBM Part

In the below example, we use Raspberry Pi 4 2GB as the sample OOBM controller.

It includes Ethernet, USB, and GPIO expansion pin that suitable for both (and only) **Power control** and **Serial Console** access to compute nodes. 

Pi 3 rev.B, as another HW option, may also suitable for this example.


### Enable Hardware Watchdog on Pi

1. Add modules in `/etc/modules`:
    ```=config
    # /etc/modules: kernel modules to load at boot time.
    #
    # This file contains the names of kernel modules that should be loaded
    # at boot time, one per line. Lines beginning with "#" can be ignored.

    bcm2835_wdt  # For Pi 4
    bcm2708_wdog # For Pi 3
    ```
2. Enable kernel module: reboot device or use modprobe:
    ```=shell
    $ modprobe bcm2835_wdt        // For Pi 4
    $ modprobe bcm2708_wdog       // For Pi 3
    ```
3. Install watchdog daemon:
    ```=shell
    sudo apt-get install watchdog
    ```
4. Modify config `/etc/watchdog`:
    ```=config
    watchdog-device = /dev/watchdog
    #...
    
    Uncomment to enable test. Setting one of these values to '0' disables it.
    These values will hopefully never reboot your machine during normal use
    (if your machine is really hung, the loadavg will go much higher than 25)
    max-load-1             = 24
    max-load-5             = 18
    max-load-15            = 12
    #...
    
    # Defaults compiled into the binary
    temperature-device     = /sys/class/thermal/thermal_zone0/temp
    max-temperature        = 75000
    ```
5. Enable service:
    ```=shell
    $ update-rc.d watchdog defaults
    update-rc.d: using dependency based boot sequencing
    
    $ /etc/init.d/watchdog start
    ```
### Power Control

For board power control, onboard GPIO pins can be used for relay control.

Please note that GPIO pins on Pi have only 3.3V as `OUTPUT_HIGH`; this is not enough signal for the relay board(5V). In case to directly trigger the relay board without any additional parts, the program should manually trigger, GPIO pins as `OUTPUT-INPUT` direction change to make sure the relay operation as expected.  



### OpenBMC for Compute Node Management

There are possible exist solution that uses OpenBMC to control all compute nodes; in this case, GPIOs uses only for power control. Some tweaks shall apply for this case, and compute node physical devices control is still facts big challenge.

#### Build SD image

1. Select any released version. The development branch may not work correctly.
2. Following Poky & OpenEmbedded instructions, with prepared [RPi configuration](https://gist.github.com/stevennick/ddcfd9ac551982f6469b2d50d103dc87).
3. Once build is complete, use dd command or Disk Tools copy file `raspberrypi3-64/obmc-phosphor-image-raspberrypi3-64.rpi-sdimg` to the SD card. This file locate under directory `build/tmp/deploy/images/`.

#### SD card size offset for OpenBMC
The following flash size offset is required for successful image build.
```=config
# Config for RaspberryPi 3 or armelf build.
FLASH_ROFS_OFFSET = "16384"
FLASH_RWFS_OFFSET = "65536"
FLASH_SIZE = "65536"
```
```=config
# Config for RaspberryPi 4 or aarch64 build.
FLASH_ROFS_OFFSET = "32768"
FLASH_RWFS_OFFSET = "131072"
FLASH_SIZE = "131072"
```
> **Note:** This setting may not bootable on mtd devices. E.g., RaspberryPi compute module.

#### Web-UI & bmcweb development

1. Use bmcweb as backend core
2. Disable XSS prevention by change `BMCWEB_INSECURE_DISABLE_XSS_PREVENTION` and `BMCWEB_INSECURE_DISABLE_CSRF_PREVENTION` to `ON` in file `CMakeLists.txt` (See below)

#### Disable bmcweb XSS/CSRF prevention
1. Import the bmcweb in https://github/openbmc/bmcweb to newer github repos.
2. clone it into local folder.
3. In file CMakeLists.txt, change`BMCWEB_INSECURE_DISABLE_XSS_PREVENTION` option:
    ```=config
     option (BMCWEB_INSECURE_DISABLE_XSS_PREVENTION "Disable XSS preventions" ON)
     option (BMCWEB_INSECURE_DISABLE_CSRF_PREVENTION "Disable CSRF preventions" ON)
    ```
4. Commit and push.
5. Changes in openbmc project:
    a. Find the location of the recipes for bmcweb:
    ```=shell
     ~/openbmc$ find . -name bmcweb*
    ```
    We find it in `./meta-phosphor/recipes-phosphor/interfaces/bmcweb_git.bb`
    b. Change bmcweb_git.bb. Redirect to new target github repo and commit version
    ```=config
     SRC_URI = "git://github.com/[yourname]/bmcweb.git"
     SRCREV = "7d044f74293a959a7fc8b63f7c4052a70b56dde9"
    ```
    > **Note:** The option `SRCREV` can be set as `${AUTOREV}`, make bitbake fetch latest commit version in automatic.
    
    c. Rebuild Image
    ```=shell
     . openbmc-env
     bitbake obmc-phosphor-image
    ```

References:
[1] [bmcweb Update for Web UI development](https://michaeltien8901.github.io/2019/03/06/bmcweb-Update-For-WEB-UI-development.html)
[2] [yocto recipe : (3) 使用 External Source 來編譯軟體](http://yi-jyun.blogspot.com/2018/05/yocto-4-external-source.html)
[3] [OpenBMC Web User Interface Development](https://github.com/openbmc/docs/blob/master/development/web-ui.md)


#### Let u-boot use VFAT environment file `uboot.env`

Because OpenBMC uses`fw_setenv` command for environment store during booting and system operation, on Raspberry Pi or any SDCard booted board, we can enable store environment on `/boot/uboot.env` file.

Below are all the necessary steps to make this happened:

1. Store default(builtin) u-boot environment into env file. 
   To save this initial file, ether access u-boot command line via attach serial console or download BSP-aware environment file from here:
   ```=shell
   U-Boot> saveenv
   Saving Environment to FAT... OK
   ```
   After successfully save environment file, use command `boot` continue boot sequence:
   ```=shell
   U-Boot> boot
   ```
2. Update `/etc/fstab`, let system mount `/boot` during boot:
   ```=config
   /dev/mmcblk0p1       /boot          auto       defaults,sync  0  0
   ```
   The device name `/dev/mmcblk0p1` may not same as your BSP, use `fdisk -l` check available devices.
   After modify `/etc/fstab`, execute `mount` command to apply the change:
   ```=shell
   $ mount -a
   ```
3. Update `/etc/fw_env.config`, uncomment `/boot/uboot.env` and comment out any other lines, let `fw_setenv` and `fw_printenv` read/write environment variables from `/boot/uboot.env`:
   ```=config
   # VFAT example
   /boot/uboot.env 0x0000          0x4000
   ```

#### Future Works

1. Use external bmcweb for UI change
2. Update BSP for GPIO power button & power status control
3. External GPIO via I2C
4. obmc agent for k8s, either k8s-externalService-OBMC or OBMC-externalService-k8s
5. **(Optional)** Compute node OS reset stack

# Reference

- [為你的 Raspberry Pi 裝上 Watchdog](https://coldnew.github.io/68fbe311/)
- [樹莓派硬體看門狗（Watchdog）：當機時自動重新開機](https://blog.gtwang.org/iot/raspberry-pi-hardware-watchdog/)