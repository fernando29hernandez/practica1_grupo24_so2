obj-m += procs_grupo24.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(shell pwd) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(shell pwd) clean

test:
	sudo dmesg -C
	sudo insmod procs_grupo24.ko
	sudo rmmod procs_grupo24.ko
	dmesg
