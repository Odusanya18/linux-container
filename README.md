# linux-container
Code here, builds a container from scratch up, direct syscalls, and quite sufficient isolation as can be gotten form docker, lxc etc. Demonstrates that containers are processes!

Provides equivalent of ```docker run -it debian:buster bash```

# Run
```shell
root@<xx,xx,xx>: $ ./start-container.sh
```

# Dependencies
1. Netsetgo https://github.com/teddyking/netsetgo (if you need internet connectivity, run the NAT commands in the repo)

