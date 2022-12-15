#!/bin/bash

# It generate userdata string (script content) as output, digitalocean will
# initialize droplet on the first boot with it.
#
# this is just an example
#
# note:
# 1. size of output should be within 64K according to the api doc
# 2. the script seems to be executed on the next boot after the creation

# ex1: just a script:
cat - <<EOF
#!/bin/bash

echo do something
EOF

exit 0

# ex2: a template script and embed a tar file
sh_tar="file.tar.xz"

cat - <<EEOF
#!/bin/bash

load_payload() {
	cat - <<EOF
$(base64 "$sh_tar")
EOF
}

mktemp /root/init.abcd
load_payload | base64 -d | tar -xJ -C /root/init.abcd

echo do something and log the output | tee /var/log/00_init.log

rm -rf /root/init.abcd/
reboot
EEOF
