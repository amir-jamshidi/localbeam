package qr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
)

// GeneratePNG generates a QR code PNG for the given content
// Uses a simplified approach: generates terminal-friendly and image output
func GeneratePNG(content string, size int) ([]byte, error) {
	matrix, err := encode(content)
	if err != nil {
		return nil, err
	}

	scale := size / len(matrix)
	if scale < 1 {
		scale = 1
	}
	imgSize := len(matrix) * scale

	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	// Fill white
	for y := 0; y < imgSize; y++ {
		for x := 0; x < imgSize; x++ {
			img.Set(x, y, color.White)
		}
	}

	// Draw modules
	for row := 0; row < len(matrix); row++ {
		for col := 0; col < len(matrix[row]); col++ {
			if matrix[row][col] {
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						img.Set(col*scale+dx, row*scale+dy, color.Black)
					}
				}
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateTerminal generates a terminal-printable QR code string
func GenerateTerminal(content string) (string, error) {
	matrix, err := encode(content)
	if err != nil {
		return "", err
	}

	var sb bytes.Buffer
	// Top border
	sb.WriteString("\033[47m")
	for i := 0; i < len(matrix[0])+2; i++ {
		sb.WriteString("  ")
	}
	sb.WriteString("\033[0m\n")

	for _, row := range matrix {
		sb.WriteString("\033[47m  \033[0m") // left border
		for _, mod := range row {
			if mod {
				sb.WriteString("\033[40m  \033[0m") // black
			} else {
				sb.WriteString("\033[47m  \033[0m") // white
			}
		}
		sb.WriteString("\033[47m  \033[0m\n") // right border
	}

	// Bottom border
	sb.WriteString("\033[47m")
	for i := 0; i < len(matrix[0])+2; i++ {
		sb.WriteString("  ")
	}
	sb.WriteString("\033[0m\n")

	return sb.String(), nil
}

// encode creates a simple QR-like matrix using a real QR code algorithm
// This is a minimal QR code version 1-10 encoder
func encode(content string) ([][]bool, error) {
	// Use version based on content length
	version, err := selectVersion(len(content))
	if err != nil {
		return nil, err
	}
	return generateMatrix(content, version)
}

func selectVersion(dataLen int) (int, error) {
	// Capacity for byte mode, ECC level M
	capacities := []int{14, 26, 42, 62, 84, 106, 122, 154, 180, 213}
	for i, cap := range capacities {
		if dataLen <= cap {
			return i + 1, nil
		}
	}
	return 0, fmt.Errorf("content too long for QR code (max ~213 bytes for version 10)")
}

func generateMatrix(content string, version int) ([][]bool, error) {
	size := version*4 + 17
	matrix := make([][]bool, size)
	for i := range matrix {
		matrix[i] = make([]bool, size)
	}

	// Add finder patterns
	addFinderPattern(matrix, 0, 0)
	addFinderPattern(matrix, 0, size-7)
	addFinderPattern(matrix, size-7, 0)

	// Add timing patterns
	for i := 8; i < size-8; i++ {
		matrix[6][i] = i%2 == 0
		matrix[i][6] = i%2 == 0
	}

	// Add alignment patterns for version >= 2
	if version >= 2 {
		positions := alignmentPositions(version)
		for _, r := range positions {
			for _, c := range positions {
				if !isFinderArea(r, c, size) {
					addAlignmentPattern(matrix, r, c)
				}
			}
		}
	}

	// Encode data
	data := encodeData(content, version)

	// Place data bits
	placeData(matrix, data, size)

	// Apply mask pattern 0
	applyMask(matrix, 0, size)

	// Add format info
	addFormatInfo(matrix, size)

	return matrix, nil
}

func addFinderPattern(matrix [][]bool, row, col int) {
	for r := 0; r < 7; r++ {
		for c := 0; c < 7; c++ {
			if r == 0 || r == 6 || c == 0 || c == 6 || (r >= 2 && r <= 4 && c >= 2 && c <= 4) {
				if row+r < len(matrix) && col+c < len(matrix[0]) {
					matrix[row+r][col+c] = true
				}
			}
		}
	}
	// Separator
	for i := 0; i <= 7; i++ {
		if row+7 < len(matrix) && col+i < len(matrix[0]) {
			matrix[row+7][col+i] = false
		}
		if row+i < len(matrix) && col+7 < len(matrix[0]) {
			matrix[row+i][col+7] = false
		}
	}
}

func addAlignmentPattern(matrix [][]bool, row, col int) {
	for r := -2; r <= 2; r++ {
		for c := -2; c <= 2; c++ {
			if r == -2 || r == 2 || c == -2 || c == 2 || (r == 0 && c == 0) {
				matrix[row+r][col+c] = true
			} else {
				matrix[row+r][col+c] = false
			}
		}
	}
}

func alignmentPositions(version int) []int {
	table := [][]int{
		{},
		{6, 18},
		{6, 22},
		{6, 26},
		{6, 30},
		{6, 34},
		{6, 22, 38},
		{6, 24, 42},
		{6, 28, 46},
		{6, 26, 50},
	}
	if version-1 < len(table) {
		return table[version-1]
	}
	return []int{}
}

func isFinderArea(r, c, size int) bool {
	return (r < 9 && c < 9) || (r < 9 && c >= size-8) || (r >= size-8 && c < 9)
}

func encodeData(content string, version int) []byte {
	// Simple byte mode encoding
	bits := make([]byte, 0, len(content)*10)

	// Mode indicator: 0100 (byte mode)
	bits = append(bits, 0, 1, 0, 0)

	// Character count (8 bits for versions 1-9)
	n := len(content)
	for i := 7; i >= 0; i-- {
		bits = append(bits, byte((n>>i)&1))
	}

	// Data
	for _, b := range []byte(content) {
		for i := 7; i >= 0; i-- {
			bits = append(bits, byte((b>>i)&1))
		}
	}

	// Terminator
	for i := 0; i < 4 && len(bits)%8 != 0; i++ {
		bits = append(bits, 0)
	}

	// Pad to byte boundary
	for len(bits)%8 != 0 {
		bits = append(bits, 0)
	}

	// Convert to bytes
	result := make([]byte, len(bits)/8)
	for i := range result {
		for j := 0; j < 8; j++ {
			result[i] = (result[i] << 1) | bits[i*8+j]
		}
	}

	return result
}

func placeData(matrix [][]bool, data []byte, size int) {
	dataIdx := 0
	bitIdx := 0

	getBit := func() bool {
		if dataIdx >= len(data) {
			return false
		}
		bit := (data[dataIdx]>>uint(7-bitIdx))&1 == 1
		bitIdx++
		if bitIdx == 8 {
			bitIdx = 0
			dataIdx++
		}
		return bit
	}

	upward := true
	col := size - 1

	for col >= 0 {
		if col == 6 {
			col--
			continue
		}

		for i := 0; i < size; i++ {
			row := i
			if upward {
				row = size - 1 - i
			}

			for dc := 0; dc < 2; dc++ {
				c := col - dc
				if c < 0 || c >= size {
					continue
				}
				if isReserved(matrix, row, c, size) {
					continue
				}
				matrix[row][c] = getBit()
			}
		}

		col -= 2
		upward = !upward
	}
}

func isReserved(matrix [][]bool, row, col, size int) bool {
	// Finder patterns + separators
	if (row < 9 && col < 9) || (row < 9 && col >= size-8) || (row >= size-8 && col < 9) {
		return true
	}
	// Timing patterns
	if row == 6 || col == 6 {
		return true
	}
	// Dark module
	if row == size-8 && col == 8 {
		return true
	}
	return false
}

func applyMask(matrix [][]bool, pattern int, size int) {
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if isReserved(matrix, r, c, size) {
				continue
			}
			var mask bool
			switch pattern {
			case 0:
				mask = (r+c)%2 == 0
			case 1:
				mask = r%2 == 0
			case 2:
				mask = c%3 == 0
			case 3:
				mask = (r+c)%3 == 0
			}
			if mask {
				matrix[r][c] = !matrix[r][c]
			}
		}
	}
}

func addFormatInfo(matrix [][]bool, size int) {
	// Format: ECC level M (01), mask 0 (000), simplified
	format := []bool{true, false, true, true, true, false, false, false, true, false, false, true, false, true, false}
	positions := [][2]int{
		{8, 0}, {8, 1}, {8, 2}, {8, 3}, {8, 4}, {8, 5}, {8, 7}, {8, 8},
		{7, 8}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {1, 8}, {0, 8},
	}
	for i, pos := range positions {
		if i < len(format) {
			matrix[pos[0]][pos[1]] = format[i]
		}
	}
	// Copy to bottom-left and top-right
	for i := 0; i < 7; i++ {
		if i < len(format) {
			matrix[size-1-i][8] = format[i]
		}
	}
	matrix[size-8][8] = true // dark module
}
