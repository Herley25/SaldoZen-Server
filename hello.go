package main

import (
	"fmt"
	"strings"
)

func reconstruirTexto(array []string, textoOriginal string) string {
	// Prepara palavras do texto original
	textoPreparado := strings.ReplaceAll(textoOriginal, ",", ",")
	palavrasOriginais := strings.Fields(strings.ToLower(textoPreparado))

	// Cópia do array de palavras disponíveis
	disponiveis := make([]string, len(array))
	copy(disponiveis, array)

	var resultado []string

	for _, palavra := range palavrasOriginais {
		encontrada := false
		for i, disp := range disponiveis {
			if strings.ToLower(disp) == palavra {
				resultado = append(resultado, disp)
				// Remove da lista de disponíveis
				disponiveis = append(disponiveis[:i], disponiveis[i+1:]...)
				encontrada = true
				break
			}
		}
		if !encontrada {
			resultado = append(resultado, palavra)
		}
	}

	return strings.Join(resultado, " ")
}

func main() {
	array := []string{"Seja", "vindo", "bem", ",", "Tudo", "você", "com", "bem"}
	textoOriginal := "Seja bem vindo, tudo bem com você"

	resultado := reconstruirTexto(array, textoOriginal)
	fmt.Println(resultado)
}
