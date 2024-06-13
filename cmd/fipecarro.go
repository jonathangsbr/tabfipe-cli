/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/spf13/cobra"

	"github.com/jonathangsbr/tabfipe-api-gateway/entity/veiculo"
	"github.com/jonathangsbr/tabfipe-api-gateway/entity/veiculoHelper"
	"github.com/jonathangsbr/tabfipe-api-gateway/pkg/fipeCarro"
)

var mapp map[string]string

var printA = false
var printO = ""
var printN = false

var a = veiculoHelper.VeiculoMarca{}

var o = ""
var n = ""

var modelos = []veiculoHelper.Modelo{}

// fipecarroCmd represents the fipecarro command
var fipecarroCmd = &cobra.Command{
	Use:   "fipecarro",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		mapp = make(map[string]string)
		mapp["AUTOMATICO"] = "AUT."
		mapp["MANUAL"] = "MEC."

		if printA {
			printMarcas()
		}
		if len(a.Marca) > 0 {
			matchMarca()
			if len(printO) > 0 {
				fmt.Println(a.Marca + ": ")
				printModelos()
				return
			}
			if len(n) > 0 {
				printCarro()
				return
			}
			if printN {
				printAnos()
				return
			}
		}
	},
}

func matchMarca() {
	m := getMarcas()
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	for _, v := range m {
		ra, _, _ := transform.String(t, v.Marca)
		rb, _, _ := transform.String(t, a.Marca)
		if strings.Contains(strings.ToUpper(ra), strings.ToUpper(rb)) {
			a = v
			_, m, err := fipeCarro.GetModelos(a)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(0)
			}
			modelos = m
			return
		}
	}
	fmt.Println("Marca não encontrada: " + a.Marca)
	os.Exit(0)
}

func getMarcas() []veiculoHelper.VeiculoMarca {
	m, err := fipeCarro.GetMarcasObj()
	if err != nil {
		fmt.Printf("%s", err)
		return m
	}
	return m
}

func printMarcas() {
	m := getMarcas()
	for i, v := range m {
		fmt.Printf("%s\t\t", v.Marca)
		if len(v.Marca) < 8 {
			fmt.Printf("\t")
		}
		if i%3 == 0 {
			fmt.Println()
		}
	}
}

func printModelos() {
	for _, v := range modelos {
		if printO == "todos" || strings.Contains(strings.ToUpper(v.ModeloNome), strings.ToUpper(printO)) {
			fmt.Printf(a.Marca+" %s\t\t\t", v.ModeloNome)
			fmt.Println()
		}
	}
}

func getAnos() ([]veiculoHelper.Ano, veiculo.VeiculoModel) {
	ano := modelos[0]
	args := strings.Split(o, " ")
	p := false
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	for _, v := range modelos {
		for i, a := range args {
			ra, _, _ := transform.String(t, v.ModeloNome)
			rb, _, _ := transform.String(t, a)
			ra = strings.ToUpper(ra)
			rb = strings.ToUpper(rb)
			if mapp[rb] != "" {
				rb = mapp[rb]
			}
			if strings.Contains(ra, rb) {
				if i == len(args)-1 {
					ano.CodigoModelo = v.CodigoModelo
					o = v.ModeloNome
					p = true
				}
				continue
			}
			break
		}
		if p {
			break
		}
	}
	v := fipeCarro.GetVeiculoModeloBase()
	v.SetMarca(a)
	a, err := fipeCarro.GetAnos(&v, ano)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return a, v
}

func printAnos() {
	an, _ := getAnos()
	fmt.Println(a.Marca + " " + o + ":")
	for _, vl := range an {
		fmt.Printf("-%s\n", vl.AnoComb)
	}
}

func printCarro() {
	a, v := getAnos()
	ano := veiculoHelper.Ano{}
	if strings.Contains(strings.ToUpper(n), "ZERO") || strings.Contains(strings.ToUpper(n), "KM") || strings.Contains(strings.ToUpper(n), "32000") || strings.ToUpper(n) == "0" {
		n = "Zero KM"
	}
	for _, vl := range a {
		if strings.Contains(vl.AnoComb, n) {
			ano = vl
			break
		}
	}
	if (ano == veiculoHelper.Ano{}) {
		fmt.Printf("Nenhum %s encontrado com o ano %s\n", o, n)
		return
	}

	vr, err := fipeCarro.GetCarro(&v, ano)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Println(vr.ToString())
}

func init() {
	rootCmd.AddCommand(fipecarroCmd)
	fipecarroCmd.Flags().BoolVarP(&printA, "marcas", "A", false, "Printar todas as marcas de carros disponíveis")
	fipecarroCmd.Flags().StringVarP(&a.Marca, "marca", "a", "", "flag de seleção da marca para busca")
	fipecarroCmd.Flags().StringVarP(&printO, "modelos", "O", "", "Printar os modelos de carro da marca selecionada, ")
	fipecarroCmd.Flags().StringVarP(&o, "modelo", "o", "", "flag de seleção do modelo do carro")
	fipecarroCmd.Flags().BoolVarP(&printN, "anos", "N", false, "Printar todas os anos do modelo de carro selecionado")
	fipecarroCmd.Flags().StringVarP(&n, "ano", "n", "", "flag de seleção do ano do carro")
}
