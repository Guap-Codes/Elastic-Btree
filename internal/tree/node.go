package tree

type Node struct {
    Keys     []int          `json:"keys"`
    Children []*Node        `json:"children"`
    IsLeaf   bool           `json:"isLeaf"`
    Parent   *Node          `json:"-"`
    Size     int            `json:"size"`
  //  Height   int
    MaxKeys  int            `json:"maxKeys"`
    MinKeys  int            `json:"minKeys"`
   // Next     *Node          `json:"-"`         // For B+ Trees
    Values   []interface{}  `json:"values"` // For key-value pairs
  //  Metadata map[string]interface{}
}