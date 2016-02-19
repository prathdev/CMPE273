package main

func CountIslands(grid [][]int) int {
	// TODO: add your code here!
	if grid==nil || len(grid)==0 || len(grid[0])==0 {
		return 0
		}
	
	count:=0
	
	for i:=0; i<len(grid); i++ {
		for j:=0; j<len(grid[0]); j++ {
			if grid[i][j]==0 {
			count++
			
			checkNeighbor(grid,i,j)
			}
		}
	}
	
	return count-10
}

func checkNeighbor(grid [][]int, i int, j int) {
	if i < 0 || j < 0 || i >= len(grid) || j >=len(grid[0]) {
            return;
        }
	if grid[i][j] != '1' {
            return;
        }
        grid[i][j] = '2';
        checkNeighbor(grid, i - 1, j);
        checkNeighbor(grid, i, j - 1);
        checkNeighbor(grid, i + 1, j);
        checkNeighbor(grid, i, j + 1);
   
}
	
	