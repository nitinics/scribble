#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <stddef.h>             /* offsetof */
#include <net/if_arp.h>
#include <net/if.h>
#include <linux/sockios.h>
#include <netinet/in.h>
#if __GLIBC__ >=2 && __GLIBC_MINOR >= 1
#include <netpacket/packet.h>
#include <net/ethernet.h>
#else
#include <asm/types.h>
#include <linux/if_ether.h>
#endif

#define IFNAME "eth0:2"
#define HOST "10.0.0.1"
#define ifreq_offsetof(x)  offsetof(struct ifreq, x)

int main(int argc, char **argv) {

	if (argc != 2)
	{
		printf( "usage: %s src_ip \n", argv[0] );
	}
	else
	{
        struct ifreq ifr;
        struct sockaddr_in sai;
        int sockfd;                     /* socket fd we use to manipulate stuff with */
        int selector;
        unsigned char mask;

        char *p;
	char *hostip = argv[1];


        /* Create a channel to the NET kernel. */
        sockfd = socket(AF_INET, SOCK_DGRAM, 0);

        /* get interface name */
        strncpy(ifr.ifr_name, IFNAME, IFNAMSIZ);

        memset(&sai, 0, sizeof(struct sockaddr));
        sai.sin_family = AF_INET;
        sai.sin_port = 0;
	
	printf ("hostip: %s \n", hostip);

        sai.sin_addr.s_addr = inet_addr(hostip);

        p = (char *) &sai;
        memcpy( (((char *)&ifr + ifreq_offsetof(ifr_addr) )),
                        p, sizeof(struct sockaddr));

        ioctl(sockfd, SIOCSIFADDR, &ifr);
        ioctl(sockfd, SIOCGIFFLAGS, &ifr);

	ifr.ifr_hwaddr.sa_data[0] = 0xDE;
	ifr.ifr_hwaddr.sa_data[1] = 0xAD;
	ifr.ifr_hwaddr.sa_data[2] = 0xBE;
	ifr.ifr_hwaddr.sa_data[3] = 0xEF;
	ifr.ifr_hwaddr.sa_data[4] = 0xCA;
	ifr.ifr_hwaddr.sa_data[5] = 0xFE;
	ifr.ifr_hwaddr.sa_family = ARPHRD_ETHER;

        ifr.ifr_flags |= IFF_UP | IFF_RUNNING;
        // ifr.ifr_flags &= ~selector;  // unset something

        ioctl(sockfd, SIOCSIFFLAGS, &ifr);
	ioctl(sockfd, SIOCSIFHWADDR, &ifr);
        close(sockfd);
        return 0;
	}
}
