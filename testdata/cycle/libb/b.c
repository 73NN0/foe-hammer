#include "b.h"
#include "a.h"

int func_b(int x) {
    if (x <= 0) {
        return 1;
    }
    return func_a(x - 1);  // calls liba
}
