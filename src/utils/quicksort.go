package utils

import (
  "math/rand"

  // ## local shits ##
  "../structs"
)

// ########################################### UTILITY FUNCTIONS
func QuicksortProducts(products structs.Products) structs.Products {
  if len(products) < 2 {
      return products
  }

  left, right := 0, len(products)-1

  pivot := rand.Int() % len(products)

  products[pivot], products[right] = products[right], products[pivot]

  for i, _ := range products {
      if products[i].URL < products[right].URL {
          products[left], products[i] = products[i], products[left]
          left++
      }
  }

  products[left], products[right] = products[right], products[left]

  QuicksortProducts(products[:left])
  QuicksortProducts(products[left+1:])

  return products
}
