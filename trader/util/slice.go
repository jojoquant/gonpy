package util

// xs := []int{2, 4, 6, 8}
// ys := []string{"C", "B", "K", "A"}
// fmt.Println(
//     SliceIndex(len(xs), func(i int) bool { return xs[i] == 5 }),
//     SliceIndex(len(xs), func(i int) bool { return xs[i] == 6 }),
//     SliceIndex(len(ys), func(i int) bool { return ys[i] == "Z" }),
//     SliceIndex(len(ys), func(i int) bool { return ys[i] == "A" }))
func SliceIndex(limit int, predicate func(i int) bool) int {
    for i := 0; i < limit; i++ {
        if predicate(i) {
            return i
        }
    }
    return -1
}