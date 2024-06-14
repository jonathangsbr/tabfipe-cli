/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jonathangsbr/tabfipe-api-gateway/v2/pkg/fipeCarro"
	fipeVeiculo "github.com/jonathangsbr/tabfipe-cli/internal"
)

var fipecarroCmd = &cobra.Command{
	Use:   "fipecarro",
	Short: "Comando para buscar na fipe de carro",
	Long:  `###`,
	Run: func(cmd *cobra.Command, args []string) {

		cli := &fipeVeiculo.Cli{}
		cli.SetFipeVeiculo(fipeCarro.NewCarro())

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
	rootCmd.AddCommand(fipecarroCmd)

	fipecarroCmd.Flags().BoolVarP(&fipeVeiculo.PrintA, "marcas", "A", false, "Printar todas as marcas de carros disponíveis")
	fipecarroCmd.Flags().StringVarP(&fipeVeiculo.A.Marca, "marca", "a", "", "flag de seleção da marca para busca")
	fipecarroCmd.Flags().StringVarP(&fipeVeiculo.PrintO, "modelos", "O", "", "Printar os modelos de carro da marca selecionada, ")
	fipecarroCmd.Flags().StringVarP(&fipeVeiculo.O, "modelo", "o", "", "flag de seleção do modelo do carro")
	fipecarroCmd.Flags().BoolVarP(&fipeVeiculo.PrintN, "anos", "N", false, "Printar todas os anos do modelo do carro selecionado")
	fipecarroCmd.Flags().StringVarP(&fipeVeiculo.N, "ano", "n", "", "flag de seleção do ano do carro")
}
