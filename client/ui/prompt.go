package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func printGreeting() {
	fmt.Println("Welcome to TCP Chat!")
	fmt.Println("         _nnnn_")
	fmt.Println("        dGGGGMMb")
	fmt.Println("       @p~qp~~qMb")
	fmt.Println("       M|@||@) M|")
	fmt.Println("       @,----.JM|")
	fmt.Println("      JS^\\__/  qKL")
	fmt.Println("     dZP        qKRb")
	fmt.Println("    dZP          qKKb")
	fmt.Println("   fZP            SMMb")
	fmt.Println("   HZM            MMMM")
	fmt.Println("   FqM            MMMM")
	fmt.Println(" __| \".        |\\dS\"qML")
	fmt.Println(" |    `.       | `' \\Zq")
	fmt.Println("_)      \\.___.,|     .'")
	fmt.Println("\\____   )MMMMMP|   .'")
	fmt.Println("     `-'       `--'")
}

func readName() {
	fmt.Printf("\nPlease enter your name: ")
	br := bufio.NewReader(os.Stdin)
	for {
		buf, err := br.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		if len(buf) > 1 {
			if len(buf) > maxNameLength {
				buf = buf[:maxNameLength]
			}
			name = fmt.Sprintf("%s: ", buf[:len(buf)-1])
			break
		}
		fmt.Printf("Name cannot be missing. Please enter your name: ")
		br.Reset(os.Stdin)
	}
}
