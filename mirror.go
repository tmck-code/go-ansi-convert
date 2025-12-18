package main

import "fmt"

func main() {
	// Print ASCII characters (0-127)
	fmt.Println("=== ASCII Characters (0-127) ===")
	fmt.Println("\nPrintable ASCII (32-126):")
	for i := 32; i <= 126; i++ {
		if (i-32)%16 == 0 && i != 32 {
			fmt.Println()
		}
		fmt.Printf("%c ", i)
	}
	fmt.Println("\n")

	fmt.Println("\n--- ASCII WITH UNICODE MIRRORS ---")
	asciiUnicodeMirrors := [][2]rune{
		{'!', '¡'}, {'?', '¿'}, // Inverted punctuation
		{'6', '9'}, // Rotational
		{'A', '∀'}, {'E', 'Ǝ'}, {'ɐ', 'a'}, {'ə', 'e'}, // Turned letters
	}
	for _, pair := range asciiUnicodeMirrors {
		fmt.Printf("%c ↔ %c\n", pair[0], pair[1])
	}

	// Box Drawing Characters organized by mirror pairs
	fmt.Println("=== Box Drawing Characters - Organized by Mirroring ===")

	// Horizontal mirror pairs (left-right)
	fmt.Println("\n--- HORIZONTAL MIRROR PAIRS (Left ↔ Right) ---")
	hMirrorPairs := [][2]rune{
		{'<', '>'}, {'(', ')'}, {'[', ']'}, {'{', '}'},
		{'/', '\\'}, {'b', 'd'}, {'p', 'q'},
		{'B', 'ᗺ'}, {'C', 'Ɔ'}, {'D', 'ᗡ'}, {'E', 'Ǝ'},
		{'F', 'ꟻ'}, {'G', 'ວ'},
		{'J', 'ᒐ'}, {'K', 'ꓘ'}, {'L', '⅃'},
		{'N', 'И'}, {'O', 'O'}, // O mirrors itself
		{'P', 'ᑫ'}, {'Q', 'Ϙ'}, {'R', 'Я'}, {'S', 'Ƨ'},
		{'a', 'ɒ'}, {'b', 'd'}, {'c', 'ɔ'}, {'d', 'b'},
		{'e', 'ɘ'}, {'f', 'ᆿ'}, {'g', 'ϱ'}, {'h', '⑁'},
		{'j', 'ᒑ'}, {'k', 'ʞ'},
		{'p', 'q'}, {'q', 'p'}, {'r', 'ɿ'}, {'s', 'ƨ'}, {'t', 'ɟ'}, {'y', 'γ'},
		{'┌', '┐'}, {'┍', '┑'}, {'┎', '┒'}, {'┏', '┓'},
		{'└', '┘'}, {'┕', '┙'}, {'┖', '┚'}, {'┗', '┛'},
		{'├', '┤'}, {'┝', '┥'}, {'┞', '┦'}, {'┟', '┧'},
		{'┠', '┨'}, {'┡', '┩'}, {'┢', '┪'}, {'┣', '┫'},
		{'╒', '╕'}, {'╓', '╖'}, {'╔', '╗'},
		{'╘', '╛'}, {'╙', '╜'}, {'╚', '╝'},
		{'╞', '╡'}, {'╟', '╢'}, {'╠', '╣'},
		{'╭', '╮'}, {'╰', '╯'},
		{'╴', '╶'}, {'╸', '╺'}, {'╼', '╾'},
		{'▖', '▗'}, {'▘', '▝'}, {'▌', '▐'},
		{'▙', '▜'}, {'▚', '▞'},
	}
	for _, pair := range hMirrorPairs {
		fmt.Printf("%c ↔ %c\n", pair[0], pair[1])
	}
	fmt.Println()

	// Vertical mirror pairs (top-bottom)
	fmt.Println("--- VERTICAL MIRROR PAIRS (Top ↔ Bottom) ---")
	vMirrorPairs := [][2]rune{
		{'A', '∀'}, // Also ꓯ
		{'B', 'ꓭ'}, {'C', 'ꓛ'}, {'D', 'ꓷ'}, {'E', 'ꓱ'}, {'F', 'ꓞ'},
		{'G', 'ꓨ'}, {'J', 'ꓩ'}, {'K', 'ꓘ'}, {'L', 'ꓶ'}, {'M', 'ꟽ'},
		{'N', 'И'}, {'P', 'Ԁ'}, {'Q', 'Ό'}, {'R', 'ꓤ'}, {'T', 'ꓕ'}, {'Y', '⅄'},
		{'U', 'ꓵ'}, {'V', 'ꓥ'}, {'W', 'M'}, // or ꟽ
		{'^', 'v'}, {'w', 'm'}, {'u', 'n'},
		{'a', 'ɐ'}, {'b', 'q'}, {'c', 'ɔ'}, {'d', 'p'}, {'e', 'ǝ'},
		{'f', 'ɟ'}, {'g', 'ƃ'}, {'h', 'ɥ'}, {'i', 'ᴉ'}, {'j', 'ɾ'},
		{'k', 'ʞ'}, {'m', 'ɯ'}, {'n', 'u'}, {'p', 'd'},
		{'q', 'b'}, {'r', 'ɹ'}, {'t', 'ʇ'}, {'u', 'n'},
		{'v', 'ʌ'}, {'w', 'ʍ'}, {'y', 'ʎ'},
		{'┌', '└'}, {'┍', '┕'}, {'┎', '┖'}, {'┏', '┗'},
		{'┐', '┘'}, {'┑', '┙'}, {'┒', '┚'}, {'┓', '┛'},
		{'┬', '┴'}, {'┭', '┵'}, {'┮', '┶'}, {'┯', '┷'},
		{'┰', '┸'}, {'┱', '┹'}, {'┲', '┺'}, {'┳', '┻'},
		{'╒', '╘'}, {'╓', '╙'}, {'╔', '╚'},
		{'╕', '╛'}, {'╖', '╜'}, {'╗', '╝'},
		{'╤', '╧'}, {'╥', '╨'}, {'╦', '╩'},
		{'╭', '╰'}, {'╮', '╯'},
		{'╵', '╷'}, {'╹', '╻'}, {'╽', '╿'},
		{'▀', '▄'}, {'▔', '▁'},
		{'▖', '▘'}, {'▗', '▝'},
		{'▙', '▛'}, {'▚', '▞'},
	}
	for _, pair := range vMirrorPairs {
		fmt.Printf("%c ↕ %c\n", pair[0], pair[1])
	}
	fmt.Println()

	// Symmetric characters (mirror themselves both ways)
	fmt.Println("--- SYMMETRIC (Mirror themselves H & V) ---")
	symmetric := []rune{
		'─', '━', '│', '┃', '═', '║',
		'┼', '┽', '┾', '┿', '╀', '╁', '╂', '╃', '╄', '╅', '╆', '╇', '╈', '╉', '╊', '╋',
		'╪', '╫', '╬', '╳',
		'█', '▬', '░', '▒', '▓',
	}
	for i, r := range symmetric {
		if i%12 == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("%c ", r)
	}
	fmt.Println("\n")

	// Horizontally symmetric only (mirror themselves left-right)
	fmt.Println("--- HORIZONTALLY SYMMETRIC ONLY (Mirror H, not V) ---")
	hSymmetric := []rune{
		'H', 'A', 'I', 'M', 'l',
		'i', 'T', 'U', 'V', 'W', 'X', 'Y',
		'm', 'n', 'o', 'u', 'v', 'w', 'x',
		'┬', '┭', '┮', '┯', '┰', '┱', '┲', '┳',
		'┴', '┵', '┶', '┷', '┸', '┹', '┺', '┻',
		'╤', '╥', '╦', '╧', '╨', '╩',
		'⎺', '⎻', '⎼', '⎽',
		'▀', '▄', '▔', '▁', '▂', '▃', '▅', '▆', '▇',
	}
	for i, r := range hSymmetric {
		if i%12 == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("%c ", r)
	}
	fmt.Println("\n")

	// Vertically symmetric only (mirror themselves top-bottom)
	fmt.Println("--- VERTICALLY SYMMETRIC ONLY (Mirror V, not H) ---")
	vSymmetric := []rune{
		'H','I','O', 'X','l','o','x',
		'├', '┝', '┞', '┟', '┠', '┡', '┢', '┣',
		'┤', '┥', '┦', '┧', '┨', '┩', '┪', '┫',
		'╞', '╟', '╠', '╡', '╢', '╣',
		'╴', '╵', '╶', '╷', '╸', '╹', '╺', '╻', '╼', '╽', '╾', '╿',
		'▌', '▍', '▎', '▏', '▐', '▕',
	}
	for i, r := range vSymmetric {
		if i%12 == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("%c ", r)
	}
	fmt.Println("\n")

	// No symmetry (unique in all directions)
	fmt.Println("--- NO SYMMETRY (Unique orientation) ---")
	noSymmetry := []rune{
		'┄', '┅', '┆', '┇', '┈', '┉', '┊', '┋',
		'╱', '╲',
		'▖', '▗', '▘', '▝', '▙', '▚', '▛', '▜', '▞', '▟',
		'▉', '▊', '▋',
	}
	for _, r := range noSymmetry {
		fmt.Printf("%c ", r)
	}
	fmt.Println("\n")

	// Diagonal pairs
	fmt.Println("--- DIAGONAL PAIRS ---")
	diagonalPairs := [][2]rune{
		{'╱', '╲'},
		{'▛', '▟'}, {'▜', '▙'},
	}
	for _, pair := range diagonalPairs {
		fmt.Printf("%c ↔ %c\n", pair[0], pair[1])
	}
	fmt.Println()
}
