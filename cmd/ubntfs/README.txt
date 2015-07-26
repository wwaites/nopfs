9P Filesystem for Ubiquiti Radios
=================================

To install, make sure your Go toolchain is built for
ARM without floating point support

   export GOROOT=/where/go/is/installed
   export GOARCH=arm
   export GOARM=5
   (cd ${GOROOT}/src; ./make.bash)

Then build this package,

   go build .

The resulting binary, ubntfs, will run on Ubiquiti radios,
but it is too big to be stored in the only user-writeable
partition on their flash storage. Only 256KB is available
and this binary is around 3MB.

The workaround is to either scp ubntfs to the radio each
time it starts, or to have an rc.poststart script that 
fetches it. The latter can be done like this on the radio.

    cat > /etc/persistent/rc.poststart <<EOF
    #!/bin/sh

    mkdir /var/bin
    wget -O /var/bin/ubntfs http://192.0.2.1/dist/ubntfs
    chmod 755 /var/bin/ubntfs
    /var/bin/ubntfs &
    EOF

    chmod 755 /etc/persistent/rc.poststart
    cfgmtd -w -p /etc/

This will obviously fail if the host distributing the binary,
192.0.2.1, is unreachable.
