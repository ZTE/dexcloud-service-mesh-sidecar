README
===============

The directory is used to store configuration files that configs stream server(s) using tcp/udp protocol.
The config file must be a *.conf file. For example:
#dns.conf
~~~
upstream dns {
	server 192.168.0.1:53535;
}
server {
	listen 53 udp;
	proxy_responses 1;
	proxy_timeout 10s;
	proxy_pass dns;
}
~~~