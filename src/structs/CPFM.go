package structs

import "github.com/shurcooL/graphql"

// ########################################### STRUCTS
type ProductQuery struct {
  Products struct {
    Edges []struct {
      Node struct {
        Title graphql.String
        Handle graphql.String
        Description graphql.String
        PublishedAt graphql.String
        Options []struct {
          Name graphql.String
          Values []graphql.String
        }
        Variants struct {
          Edges []struct {
            Node struct {
              ID graphql.String
              AvailableForSale graphql.Boolean
              PriceV2 struct {
                CurrencyCode graphql.String
                Amount graphql.String
              }
            }
          }
        } `graphql:"variants(first: 15)"`
        Images struct {
          Edges []struct {
            Node struct {
              OriginalSrc graphql.String
            }
          }
        } `graphql:"images(first: 1)"`
      }
    }
  } `graphql:"products(first: 20)"`
}
