# Ramdisk creator for MacOS

This tool can be used to create a ramdisk in MacOS.
The ramdisk can be found in Finder to store data in memory.

## How to use
 * Help for the parameters:
```
mac-ramdisk-creator -help
```
 * Create a 1 GB ramdisk:
```
mac-ramdisk-creator -size 1024
```
 * Create a 200 MB ramdisk with name '200MB':
```
mac-ramdisk-creator -size 200 -name 200MB
```
 * Eject a ramdisk in this example with name '200MB' (Use the 'Eject' icon in Finder, or use the following command):
```
mac-ramdisk-creator -name 200MB -eject
```
