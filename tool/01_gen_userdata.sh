#!/bin/bash

# this is just an example, it generate userdata string as output
#
# note:
# 1. size of output should be within 64K according the api doc
# 2. the script seems to be executed on the next boot after the creation

# ex1: just script:
cat - <<EOF
#!/bin/bash

echo do something
EOF

exit 0

# ex2: embed a tar file into userdata
sh_tar="file.tar.xz"

cat - <<EEOF
#!/bin/bash

load_payload() {
	cat - <<EOF
$(cat $sh_tar | base64)
EOF
}

mktemp /root/init.abcd
load_payload | base64 -d | tar -xJ -C /root/init.abcd

echo do something and log the output | tee /var/log/00_init.log

rm -rf /root/init.abcd/
reboot
EEOF
