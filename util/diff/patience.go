// Patience diff algorithm stolen from
// http://bazaar.launchpad.net/~bzr-pqm/bzr/bzr.dev/view/head:/bzrlib/_patiencediff_py.py
// This is pretty much a straight-up port so possibly
// should be considered © Bram Cohen / Canonical Ltd.
package diff

import "errors"

var ErrDiff = errors.New("there are diffs")

type commonLine struct {
	a, b int
}

type commonSeq struct {
	a, b, l int
}

type diffOp string

const (
	equal  diffOp = "equal"
	insert diffOp = "insert"
	remove diffOp = "remove"
)

type diff struct {
	op      diffOp
	a, b, l int
}

func uniqueLCS(a, b []string) []commonLine {
	// map line in a -> +ve index of line in a
	// iff line only occurs once in a
	aIndex := map[string]int{}
	for i, line := range a {
		if _, ok := aIndex[line]; ok {
			aIndex[line] = -1
		} else {
			aIndex[line] = i
		}
	}
	// map index of line in b -> +ve index of line in a
	// iff line only occurs once in a or b
	// maintain ordering of positions in a and b
	// using -1 as a marker for "duplicate" gets messy :-(
	bIndex := map[string]int{}
	bPosInA := make([]int, len(b))
	for i, line := range b {
		aPos, ok := aIndex[line]
		if !ok || aPos < 0 {
			bPosInA[i] = -1
			continue
		}
		if bPos, ok := bIndex[line]; ok {
			aIndex[line] = -1
			bIndex[line] = -1
			bPosInA[i] = -1
			if bPos > 0 {
				bPosInA[bPos] = -1
			}
		} else {
			bIndex[line] = i
			bPosInA[i] = aPos
		}
	}
	// Patience.
	stacks, lasts, ptrs := []int{}, []int{}, map[int]int{}
	k := 0
	for bPos, aPos := range bPosInA {
		if aPos < 0 {
			continue
		}
		if len(stacks) > 0 && stacks[len(stacks)-1] < aPos {
			// Next line goes at the end of the stacks.
			k = len(stacks)
		} else if len(stacks) > 0 && stacks[k] < aPos &&
			(k == len(stacks)-1 || stacks[k+1] > aPos) {
			// Next line goes after current line.
			k += 1
		} else {
			// Binary search stacks to find insertion point.
			lo, hi := 0, len(stacks)
			for lo < hi {
				mid := (lo + hi) / 2
				if aPos < stacks[mid] {
					hi = mid
				} else {
					lo = mid + 1
				}
			}
			k = lo
		}
		if k > 0 {
			ptrs[bPos] = lasts[k-1]
		}
		if k < len(stacks) {
			stacks[k] = aPos
			lasts[k] = bPos
		} else {
			stacks = append(stacks, aPos)
			lasts = append(lasts, bPos)
		}
	}
	result := []commonLine{}
	if len(lasts) == 0 {
		return result
	}
	k, ok := lasts[len(lasts)-1], true
	for ok {
		result = append(result, commonLine{bPosInA[k], k})
		k, ok = ptrs[k]
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func recurseMatches(a, b []string, matches []commonLine, alo, ahi, blo, bhi int) []commonLine {
	lastLen := len(matches)
	if alo == ahi || blo == bhi {
		return matches
	}
	lastA := alo - 1
	lastB := blo - 1
	for _, pos := range uniqueLCS(a[alo:ahi], b[blo:bhi]) {
		aPos, bPos := pos.a+alo, pos.b+blo
		if lastA+1 != aPos || lastB+1 != bPos {
			matches = recurseMatches(a, b, matches, lastA+1, aPos, lastB+1, bPos)
		}
		lastA, lastB = aPos, bPos
		matches = append(matches, commonLine{aPos, bPos})
	}
	if len(matches) > lastLen {
		matches = recurseMatches(a, b, matches, lastA+1, ahi, lastB+1, bhi)
	} else if a[alo] == b[blo] {
		for ; alo < ahi && blo < bhi && a[alo] == b[blo]; alo, blo = alo+1, blo+1 {
			matches = append(matches, commonLine{alo, blo})
		}
		matches = recurseMatches(a, b, matches, alo, ahi, blo, bhi)
	} else if a[ahi-1] == b[bhi-1] {
		atop := ahi
		for ahi, bhi = ahi-1, bhi-1; ahi > alo && bhi > blo && a[ahi-1] == b[bhi-1]; ahi, bhi = ahi-1, bhi-1 {
		}
		matches = recurseMatches(a, b, matches, lastA+1, ahi, lastB+1, bhi)
		for i := 0; i < (atop - ahi); i++ {
			matches = append(matches, commonLine{ahi + i, bhi + i})
		}
	}
	return matches
}

func collapseSequences(matches []commonLine) []commonSeq {
	seqs := []commonSeq{}
	startA, startB, length := -1, -1, 0
	for _, match := range matches {
		if startA != -1 && match.a == startA+length && match.b == startB+length {
			length++
			continue
		}
		if startA != -1 {
			seqs = append(seqs, commonSeq{startA, startB, length})
		}
		startA, startB, length = match.a, match.b, 1
	}
	if length > 0 {
		seqs = append(seqs, commonSeq{startA, startB, length})
	}
	return seqs
}

func patienceDiff(a, b []string) []diff {
	matches := recurseMatches(a, b, []commonLine{}, 0, len(a), 0, len(b))
	seqs := collapseSequences(matches)

	lastA, lastB := -1, -1
	diffs := []diff{}
	for _, seq := range seqs {
		if lastA != -1 {
			l := seq.a - lastA
			if l > 0 {
				diffs = append(diffs, diff{remove, lastA, 0, l})
			}
		} else if seq.a > 0 && seq.b == 0 {
			diffs = append(diffs, diff{remove, 0, 0, seq.l})
		}

		if lastB != -1 {
			l := seq.b - lastB
			if l > 0 {
				diffs = append(diffs, diff{insert, 0, lastB, l})
			}
		} else if seq.b > 0 && seq.a == 0 {
			diffs = append(diffs, diff{insert, 0, 0, seq.l})
		}

		if seq.l > 0 {
			diffs = append(diffs, diff{equal, seq.a, seq.b, seq.l})
		}
		lastA, lastB = seq.a+seq.l, seq.b+seq.l
	}
	if lastA < len(a) {
		if lastA == -1 {
			lastA = 0
		}
		diffs = append(diffs, diff{remove, lastA, 0, len(a) - lastA})
	}
	if lastB < len(b) {
		if lastB == -1 {
			lastB = 0
		}
		diffs = append(diffs, diff{insert, 0, lastB, len(b) - lastB})
	}
	return diffs
}

func Unified(a, b []string) ([]string, error) {
	diffs := patienceDiff(a, b)
	err := ErrDiff
	if len(diffs) == 1 && diffs[0].op == equal {
		err = nil
	}
	unified := make([]string, 0, len(a)+len(b))
	for _, d := range diffs {
		switch d.op {
		case equal:
			for i := d.a; i < d.a+d.l; i++ {
				unified = append(unified, " "+a[i])
			}
		case insert:
			for i := d.b; i < d.b+d.l; i++ {
				unified = append(unified, "+"+b[i])
			}
		case remove:
			for i := d.a; i < d.a+d.l; i++ {
				unified = append(unified, "-"+a[i])
			}
		}
	}
	return unified, err
}
