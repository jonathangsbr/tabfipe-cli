/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/jonathangsbr/tabfipe-api-gateway/v2/pkg/fipeMoto"
	fipeVeiculo "github.com/jonathangsbr/tabfipe-cli/internal"
	"github.com/spf13/cobra"
)

// fipemotoCmd represents the fipemoto command
var fipemotoCmd = &cobra.Command{
	Use:   "fipemoto",
	Short: "Comando para buscar na fipe de moto",
	Long:  `###`,
	Run: func(cmd *cobra.Command, args []string) {

		cli := &fipeVeiculo.Cli{}
		cli.SetFipeVeiculo(fipeMoto.NewMoto())

		if fipeVeiculo.PrintA {
			cli.PrintMarcas()
		}
		if len(fipeVeiculo.A.Marca) > 0 {
			cli.MatchMarca()
			if len(fipeVeiculo.PrintO) > 0 {
				fmt.Println(fipeVeiculo.A.Marca + ": ")
				cli.PrintModelos()
				return
			}
			if len(fipeVeiculo.N) > 0 {
				cli.PrintVeiculo()
				return
			}
			if fipeVeiculo.PrintN {
				cli.PrintAnos()
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fipemotoCmd)

	fipemotoCmd.Flags().BoolVarP(&fipeVeiculo.PrintA, "marcas", "A", false, "Printar todas as marcas de motos disponíveis")
	fipemotoCmd.Flags().StringVarP(&fipeVeiculo.A.Marca, "marca", "a", "", "flag de seleção da marca para busca")
	fipemotoCmd.Flags().StringVarP(&fipeVeiculo.PrintO, "modelos", "O", "", "Printar os modelos de motos da marca selecionada, ")
	fipemotoCmd.Flags().StringVarP(&fipeVeiculo.O, "modelo", "o", "", "flag de seleção do modelo da moto")
	fipemotoCmd.Flags().BoolVarP(&fipeVeiculo.PrintN, "anos", "N", false, "Printar todas os anos do modelo da moto selecionado")
	fipemotoCmd.Flags().StringVarP(&fipeVeiculo.N, "ano", "n", "", "flag de seleção do ano da moto")
}
