#include <stdio.h>

int main(int argc, char const *argv[]) {

  printf("please enter the number of layers\n");
  int layers = 0;
  scanf("%i", &layers);

  // loop for each line
  for (int i = 0; i < layers; i++) {
    // insert spaces
    for (int k = layers - i; k > 0; k--) {
      printf("  ");
    }
    // print the star
    for (int j = 0; j < 2 * i + 1; j++) {
      printf("* ");
    }
    printf("\n");
  }

  return 0;
}