wfc tileset for path on grass assuming square and 4 way rotation

straight through  rotate 90
+ no rotate
L	rotate 90, 180, 270
End-one path in rotate 90, 180, 270
No path no rotate
(T junction too? -- can be added later -- 4x rotation)


Logic holds for rivers as well as paths -- but would need to be different tiles (and potentially allow bridges and fords...)

we can build a set of tiles from our input tiles, and have the tile objects create multiples of themselves rotated. Then use the simple wfc method to place them in the world.

How do we quantify what they can go next to... The other model simply had a list of other tiles, but that wouldn't be known at the point of multiplying (and rotating) the tiles.

Instead of having a slice of allowable tiles, we can have a set of allowable connectors. And each tile can list its own connectors. This would allow for future enhancement where we could add rivers as well as roads.


Just Grass: {connectors: [g,g,g,g]}
+: {connectors: [r, r, r, r]} 
I : {connectors: [r,g,r,g]}
- : {connectors: [g,r,g,r]}
L : {connectors: [r,r,g,g]}
end (top): {connectors[g,g,r,g]}


======
To build the board, we can have a Board made of [][]Tiles

and in addition to the we'll have a BuildBoard with [][]BuildCell 

Build Cell can have {connectors   []Connector} and the connector for each quadrant can be a sum of all of the possible connectors. In the current case just grass and roads, so we start with connectors [3,3,3,3]
Build loop: {
	calc entropy-board [][][]int ---- [x][y][index of possible cards]
	
	loop through the entropy-board and find the x/ys with the lowest number of possible ids.
	
	randomly choose one of the returned lowest entropy cells

	randomly choose one of the possible ids that is returned.

	update the build board and the board
}

How do we know its finished...

	
