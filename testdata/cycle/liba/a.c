#include "a.h"
#include "b.h"

int func_a(int x) {
    if (x <= 0) {
        return 0;
    }
    return func_b(x - 1);
}
