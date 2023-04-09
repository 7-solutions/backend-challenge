package services

func MaxSumPath(matrix [][]int) int {
	memo := make(map[[2]int]int)
	return dfs(matrix, 0, 0, memo)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func dfs(triangle [][]int, row, col int, memo map[[2]int]int) int {
	if row == len(triangle) {
		return 0
	}

	if value, ok := memo[[2]int{row, col}]; ok {
		return value
	}

	left := dfs(triangle, row+1, col, memo)
	right := dfs(triangle, row+1, col+1, memo)

	sum := triangle[row][col] + max(left, right)
	memo[[2]int{row, col}] = sum

	return sum
}
