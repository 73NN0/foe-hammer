#include "b.h"
#include "a.h"

int sum_of_squares(int x, int y) {
    return add(multiply(x, x), multiply(y, y));
}

int square_of_sum(int x, int y) {
    int sum = add(x, y);
    return multiply(sum, sum);
}
