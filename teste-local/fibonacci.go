package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
    //O código abaixo só é chamado uma vez 
    fmt.Println("Fibonacci foi chamada")
    first, second := 0, 1
    fmt.Println("first antes de func",first)
    fmt.Println("second antes de func",second)
	
    return func() int {
        //o código que é chamado com frequência eh esse
        //os valores das variáveis não zeram a cada chamada
        //cada 
        fmt.Println("Função func foi chamada")
        fmt.Println("first",first)
        ret := first
        fmt.Println("ret",ret)
        fmt.Println("second",second)
        first, second = second, first+second


        return ret
    }
}

func main() {
	f := fibonacci()
    t := fibonacci()
	for i := 0; i < 10; i++ {
        fmt.Println("Laço rodou. Antes de f() ser chamada. Sequência número i:",i)
		fmt.Println(f())
        fmt.Println("Laço terminou. Depois de f() ser chamada. Sequência número i:",i)
	}
    fmt.Println(t())
    fmt.Println(t())
}
