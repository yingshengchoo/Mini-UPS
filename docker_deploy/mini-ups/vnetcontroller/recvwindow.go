package vnetcontroller

type RecvWindow struct {
}

// merge inaterval for ack to world
// 1 2 3
// [[1, 3]]

// 5 6 7
// [[1, 3], [5, 7]]

// 4
// [[1, 7]]
