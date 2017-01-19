#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stddef.h>             /* offsetof */
#include <net/if.h>
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
#define HOST "192.168.1.204"
#define ifreq_offsetof(x)  offsetof(struct ifreq, x)

int main(int argc, char **argv) {

        struct ifreq ifr;
        struct sockaddr_in sai;
        int sockfd;                     /* socket fd we use to manipulate stuff with */
        int selector;
        unsigned char mask;

        char *p;


        /* Create a channel to the NET kernel. */
        sockfd = socket(AF_INET, SOCK_DGRAM, 0);

        /* get interface name */
        strncpy(ifr.ifr_name, IFNAME, IFNAMSIZ);

        memset(&sai, 0, sizeof(struct sockaddr));
        sai.sin_family = AF_INET;
        sai.sin_port = 0;

        sai.sin_addr.s_addr = inet_addr(HOST);

        p = (char *) &sai;
        memcpy( (((char *)&ifr + ifreq_offsetof(ifr_addr) )),
                        p, sizeof(struct sockaddr));

        ioctl(sockfd, SIOCSIFADDR, &ifr);
        ioctl(sockfd, SIOCGIFFLAGS, &ifr);

        ifr.ifr_flags |= IFF_UP | IFF_RUNNING;
        // ifr.ifr_flags &= ~selector;  // unset something

        ioctl(sockfd, SIOCSIFFLAGS, &ifr);
        close(sockfd);
        return 0;
}
