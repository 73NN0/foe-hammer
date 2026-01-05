#include <core.h>
#include <platform.h>

static int initialized = 0;

void core_init(void) {
    platform_init();
    initialized = 1;
}

void core_run(void) {
    if (!initialized) {
        return;
    }
    platform_print("core: running\n");
}

void core_shutdown(void) {
    platform_print("core: shutting down\n");
    platform_shutdown();
    initialized = 0;
}
