#ifndef PLATFORM_H
#define PLATFORM_H

// Platform interface - implemented by platform-unix or platform-stub

void platform_init(void);
void platform_shutdown(void);
void platform_print(const char *msg);

#endif
