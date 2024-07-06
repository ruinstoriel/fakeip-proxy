#!/bin/sh

# 启用 IPv4 转发
echo "启用 IPv4 转发..."
sysctl -w net.ipv4.ip_forward=1

# 启用 IPv6 转发
echo "启用 IPv6 转发..."
sysctl -w net.ipv6.conf.all.forwarding=1

# 确保在引导时启用 IPv4 和 IPv6 转发
echo "保存设置以在引导时启用..."
echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
echo "net.ipv6.conf.all.forwarding=1" >> /etc/sysctl.conf

echo "路由规则..."
ip rule add fwmark 1 lookup 100
ip route add local 0.0.0.0/0 dev lo table 100
echo "完成！"