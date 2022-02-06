package puzzle15

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func getXYFromIndex(i int) (int, int) {
	return i % numberOfRows, i / numberOfRows
}

func getIndexFromXY(x, y int) int {
	return x + y*numberOfRows
}
