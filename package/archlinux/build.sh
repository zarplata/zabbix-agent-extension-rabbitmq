#!/bin/bash

pkgname="zabbix-agent-extension-rabbitmq"
rm -rf src/ pkg/ $pkgname
makepkg -Cod; PKGVER=$(cd $(pwd)/src/$pkgname/ && make ver) makepkg -esd
