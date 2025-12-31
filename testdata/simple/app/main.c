#include <stdio.h>
#include "b.h"

int main(void) {
    int x = 3;
    int y = 4;
    
    printf("x = %d, y = %d\n", x, y);
    printf("sum_of_squares(%d, %d) = %d\n", x, y, sum_of_squares(x, y));
    printf("square_of_sum(%d, %d) = %d\n", x, y, square_of_sum(x, y));
    
    return 0;
}
