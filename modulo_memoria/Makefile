obj-m += mem_grupo24.o


all:
        make -C /lib/modules/$(shell uname -r)/build M=$(shell pwd) modules

clean:
        make -C /lib/modules/$(shell uname -r)/build M=$(shell pwd) clean

test:
        sudo dmesg -C
        sudo insmod meminfo.ko
        sudo rmmod meminfo.ko
        dmesg
