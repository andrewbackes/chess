package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*******************************************************************************

	Perft:

*******************************************************************************/

func perft(G *Game, depth int) (nodes, checks, castles, mates, captures, promotions, enpassant uint64) {
	var moveCount uint64

	if depth == 0 {
		return 1, 0, 0, 0, 0, 0, 0
	}

	toMove := G.PlayerToMove()
	notToMove := []Color{Black, White}[toMove]

	isChecked := G.isInCheck(toMove)
	ml := G.LegalMoves()

	for mv := range ml {
		temp := G
		temp.MakeMove(mv)

		if temp.isInCheck(toMove) == false {
			//Count it for mate:
			moveCount++
			n, c, cstl, m, cap, p, enp := perft(temp, depth-1)
			nodes += n
			checks += c + toInt(temp.isInCheck(notToMove))
			castles += cstl + toInt(isCastle(G, mv))
			mates += m
			captures += cap + toInt(isCapture(G, mv))
			promotions += p + toInt(isPromotion(G, mv))
			enpassant += enp + toInt(isEnPassant(G, mv))
		}
	}
	if moveCount == 0 && isChecked {
		mates++
	}

	return nodes, checks, castles, mates, captures, promotions, enpassant

}

/*******************************************************************************

	Perft Suite:

*******************************************************************************/

func PerftSuite(filename string, maxdepth int) {

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file")
	}
	scanner := bufio.NewScanner(f)

	type Test struct {
		fen   string
		nodes []uint64
	}
	var test []Test
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ";")

		var newTest Test
		newTest.fen = words[0]
		newTest.nodes = append(newTest.nodes, 1) // depth 0 = 1 node

		for i := 1; i < len(words); i++ {
			n, _ := strconv.ParseUint(strings.Split(words[i], " ")[1], 10, 0)
			newTest.nodes = append(newTest.nodes, n)
		}

		test = append(test, newTest)
	}
	f.Close()

	for i, t := range test {
		G, _ := FromFEN(t.fen)
		fmt.Print("FEN ", i+1, ": ")
		for depth, nodes := range t.nodes {
			/*
				if nodes >= 100000000 {
					break
				}
			*/
			if depth > maxdepth {
				break
			}
			fmt.Print("D", depth, ": ")
			perftNodes, _, _, _, _, _, _ := perft(G, depth)
			fmt.Print(map[bool]string{true: "pass. ", false: "FAIL. "}[(perftNodes == nodes)])
		}
		fmt.Print("\n")
	}

}

/*******************************************************************************

	Divide:

*******************************************************************************/

func divide(G *Game, depth int) {
	fmt.Println("Depth", depth)
	var nodes, moveCount uint64
	ml := G.LegalMoves()
	toMove := G.PlayerToMove()
	for mv := range ml {
		temp := *G
		temp.MakeMove(mv)

		if temp.isInCheck(toMove) == false {
			//Count it for mate:
			moveCount++
			n, _, _, _, _, _, _ := perft(&temp, depth-1)
			fmt.Println(mv, ":", n)
			nodes += n
		}
	}
	fmt.Println("Total: ", nodes, ". moves:", moveCount)
}

/*******************************************************************************

	Helpers:

*******************************************************************************/

func isCastle(G *Game, m Move) bool {
	from, _ := getSquares(m)
	p := G.board.OnSquare(from)
	if p.Type == King {
		if (m == "e1g1") || (m == "e1c1") || (m == "e8g8") || (m == "e8c8") {
			return true
		}
	}
	return false
}

func isCapture(G *Game, m Move) bool {
	_, to := getSquares(m)
	capPiece := G.board.OnSquare(to)
	return (capPiece.Type != None)
}

func isPromotion(G *Game, m Move) bool {
	// TODO: will not work when more notation is added
	return (len(m) > 4)
}

func isEnPassant(G *Game, m Move) bool {
	if G.history.enPassant == nil {
		return false
	}
	from, to := getSquares(m)
	p := G.board.OnSquare(from)
	return (p.Type == Pawn) && (to == *G.history.enPassant) && ((from-to)%8 != 0)
}

func toInt(b bool) uint64 {
	if b == true {
		return 1
	}
	return 0
}
