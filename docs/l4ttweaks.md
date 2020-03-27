Tweaks on L4T platform
======================

# Compute Node (L4T) Tweaks

## USB Ethernet

After link USB cable with working L4T system with computer, new network interface created and with fixed IP `192.168.55.100` for the host computer, and `192.168.55.1` point to the target board. Developers can use the ssh command to this target board:
```=shell
$ ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no \
  username@192.168.55.1
```

Since this IP address fixed for all target boards, please note there are multiple same IP addresses for multiple board connected environments. In this case, we should change the target board IP address manually.

## Disable zram Swap

```=shell
$ sudo systemctl disable nvzramconfig.service
```
Alternatively, change kernel config to disable it:
```=config
CONFIG_ZRAM = no
```

## Enable certain kernel flags for k8s installation

```=config
CONFIG_IP_VS
CONFIG_IP_VS_RR
CONFIG_IP_VS_WRR
CONFIG_IP_VS_SH
CONFIG_NF_CONNTRACK
CONFIG_NF_CONNTRACK_IPV4
CONFIG_IP_NF_NAT
CONFIG_IP_NF_RAW
CONFIG_NETFILTER_XT_TARGET_NOTRACK  # For CoreDNS node-local-cache
```

## Enable maximum performance now and after boot

1. Add file `/lib/systemd/system/jetson-clocks.service`:
    ```=config
    [Unit]
    Description=Maximize Jetson Performance

    [Service]
    ExecStart=/usr/bin/jetson_clocks

    [Install]
    WantedBy=multi-user.target
    ```
2. Execute command to enable jetson-clocks service:
    ```=shell
    $ sudo systemctl daemon-reload
    $ sudo systemctl enable jetson-clocks.service
    $ sudo systemctl start jetson-clocks.service
    ```