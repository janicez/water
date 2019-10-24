// +build solaris

package water

//#include <errno.h>
//#include <stdio.h>
//#include <sys/ioctl.h>
//#include <unistd.h>
//#include <string.h>
//#include <sys/socket.h>
//#include <sys/types.h>
//#include <stdlib.h>
//#include <stddef.h>
//#include <net/if.h>
//#include <ctype.h>
//#include <sys/stropts.h>
//#include <sys/sockio.h>
//#include <fcntl.h>
//#include <net/route.h>
//int
//snoifsipv6 (int ipFd, int tunFd, int tunFd2, int ppa) {
//  int error = 0;
//	struct lifreq ifr = {
//		.lifr_ppa = ppa,
//		.lifr_flags = IFF_IPV6
//	};
//    if (ioctl(tunFd, I_SRDOPT, RMSGD) < 0) {
//        error = -1;
//
//    } else if (ioctl(tunFd2, I_PUSH, "ip") < 0) {
//        // add the ip module
//        error = -2;
//
//    } else if (ioctl(tunFd2, SIOCSLIFNAME, &ifr) < 0) {
//        // set the name of the interface and specify it as ipv6
//        error = -3;
//
//    } else if (ioctl(ipFd, I_LINK, tunFd2) < 0) {
//        // link the device to the ipv6 router
//        error = -4;
//    }
//  return error;
//};
import "C"

import (
	"errors"
	"os"
	"strconv"
	"golang.org/x/sys/unix"
)

func openDev(config Config) (ifce *Interface, err error) {
	switch config.Name[:8] {
	case "/dev/tap":
		return newTAP(config)
	case "/dev/tun":
		return newTUN(config)
	default:
		return nil, errors.New("unrecognized driver")
	}
}

func newTAP(config Config) (ifce *Interface, err error) {
	panic("This section is currently bypassed because it is not necessary")
	if config.Name[:8] != "/dev/tap" {
		panic("TUN/TAP name must be in format /dev/tunX or /dev/tapX")
	}

	file, err := os.OpenFile("/dev/tap", os.O_RDWR, 0)

	//ppa = sErr := unix.IoctlRetInt(file, ((0x54<<16)|0x2), strconv.Atoi(config.Name()))
//	if sErr != nil {
//		return nil, err
//	}
//	file2, err2 := os.OpenFile("/dev/ip6", os.O_RDWR, 0)
//	file3, err3 := os.OpenFile("/dev/tap", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
//	if err2 != nil {
//		return nil, err2
//	}
//	if err3 != nil {
//		return nil, err3
//	}

	ifce = &Interface{isTAP: true, ReadWriteCloser: file, name: config.Name[5:]}
	return
}

func newTUN(config Config) (ifce *Interface, err error) {
	if config.Name[:8] != "/dev/tun" {
		panic("TUN/TAP name must be in format /dev/tunX or /dev/tapX")
	}

	file, err := os.OpenFile("/dev/tun", os.O_RDWR, 0)

	namenum, discard := strconv.Atoi(config.Name[5:])
	sErr := unix.IoctlSetInt(int(file.Fd()), ((0x54<<16)|0x2), namenum)
	if discard != nil {
		return nil, discard
	}
	if sErr != nil {
		return nil, sErr
	}
	file2, err2 := os.OpenFile("/dev/ip6", os.O_RDWR, 0)
	file3, err3 := os.OpenFile("/dev/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	if err2 != nil {
		return nil, err2
	}
	if err3 != nil {
		return nil, err3
	}

	C.snoifsipv6(C.int(file.Fd()), C.int(file2.Fd()), C.int(file3.Fd()), C.int(namenum));
	ifce = &Interface{isTAP: false, ReadWriteCloser: file, name: config.Name[5:]}
	return
}
