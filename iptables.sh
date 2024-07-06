iptables -t mangle -N HYSTERIA

iptables -t mangle -A HYSTERIA -p tcp -m socket --transparent -j MARK --set-mark 0x1
iptables -t mangle -A HYSTERIA -p udp -m socket --transparent -j MARK --set-mark 0x1
iptables -t mangle -A HYSTERIA -m socket -j RETURN

iptables -t mangle -A HYSTERIA ! -d 198.19.0.0/16 -j RETURN

iptables -t mangle -A HYSTERIA -p tcp -j TPROXY --on-port 9898 --on-ip 127.0.0.1 --tproxy-mark 0x1
iptables -t mangle -A HYSTERIA -p udp -j TPROXY --on-port 9898 --on-ip 127.0.0.1 --tproxy-mark 0x1

iptables -t mangle -A PREROUTING -j HYSTERIA

