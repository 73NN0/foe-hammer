#include "platform.h"
#include <stdio.h>

void platform_init(void) {
    printf("platform-unix: initialized\n");
}

void platform_shutdown(void) {
    printf("platform-unix: shutdown\n");
}

void platform_print(const char *msg) {
    printf("%s", msg);
}
