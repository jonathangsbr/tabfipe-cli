/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/jonathangsbr/tabfipe-api-gateway/v2/pkg/fipeCaminhao"
	fipeVeiculo "github.com/jonathangsbr/tabfipe-cli/internal"
	"github.com/spf13/cobra"
)

// fipecaminhaoCmd represents the fipecaminhao command
var fipecaminhaoCmd = &cobra.Command{
	Use:   "fipecaminhao",
	Short: "Comando para buscar na fipe de caminhão",
	Long:  `###`,
	Run: func(cmd *cobra.Command, args []string) {

		cli := &fipeVeiculo.Cli{}
		cli.SetFipeVeiculo(fipeCaminhao.NewCaminhao())

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
	rootCmd.AddCommand(fipecaminhaoCmd)

	fipecaminhaoCmd.Flags().BoolVarP(&fipeVeiculo.PrintA, "marcas", "A", false, "Printar todas as marcas de caminhões disponíveis")
	fipecaminhaoCmd.Flags().StringVarP(&fipeVeiculo.A.Marca, "marca", "a", "", "flag de seleção da marca para busca")
	fipecaminhaoCmd.Flags().StringVarP(&fipeVeiculo.PrintO, "modelos", "O", "", "Printar os modelos de caminhão da marca selecionada, ")
	fipecaminhaoCmd.Flags().StringVarP(&fipeVeiculo.O, "modelo", "o", "", "flag de seleção do modelo do caminhão")
	fipecaminhaoCmd.Flags().BoolVarP(&fipeVeiculo.PrintN, "anos", "N", false, "Printar todas os anos do modelo do caminhão selecionado")
	fipecaminhaoCmd.Flags().StringVarP(&fipeVeiculo.N, "ano", "n", "", "flag de seleção do ano do caminhão")
}
