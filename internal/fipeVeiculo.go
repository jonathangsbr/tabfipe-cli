package fipeV

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/jonathangsbr/tabfipe-api-gateway/v2/entity/veiculo"
	"github.com/jonathangsbr/tabfipe-api-gateway/v2/entity/veiculoHelper"
	"github.com/jonathangsbr/tabfipe-api-gateway/v2/pkg/fipeVeiculo"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var fv fipeVeiculo.FipeVeiculo

var mapp map[string]string

var PrintA = false
var PrintO = ""
var PrintN = false

var A = veiculoHelper.VeiculoMarca{}

var O = ""
var N = ""

var Modelos = []veiculoHelper.Modelo{}

type Veiculo interface {
	SetFipeVeiculo(v fipeVeiculo.FipeVeiculo)
	GetVeiculo() fipeVeiculo.FipeVeiculo
	MatchMarca()
	PrintMarcas()
	PrintModelos()
	PrintVeiculo()
	PrintAnos()
}

type Cli struct{}

func init() {
	mapp = make(map[string]string)
	mapp["AUTOMATICO"] = "AUT."
	mapp["MANUAL"] = "MEC."
}

func (c *Cli) SetFipeVeiculo(v fipeVeiculo.FipeVeiculo) {
	fv = v
}

func (c *Cli) GetVeiculo() fipeVeiculo.FipeVeiculo {
	return fv
}

func getMarcas() []veiculoHelper.VeiculoMarca {
	m, err := fv.GetMarcasObj()
	if err != nil {
		fmt.Printf("%s", err)
		return m
	}
	return m
}

func (c *Cli) MatchMarca() {
	m := getMarcas()
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	for _, v := range m {
		ra, _, _ := transform.String(t, v.Marca)
		rb, _, _ := transform.String(t, A.Marca)
		if strings.Contains(strings.ToUpper(ra), strings.ToUpper(rb)) {
			A = v
			_, m, err := fv.GetModelos(A)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(0)
			}
			Modelos = m
			return
		}
	}
	fmt.Println("Marca não encontrada: " + A.Marca)
	os.Exit(0)
}

func (c *Cli) PrintMarcas() {
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

func (c *Cli) PrintModelos() {
	for _, v := range Modelos {
		if PrintO == "todos" || strings.Contains(strings.ToUpper(v.ModeloNome), strings.ToUpper(PrintO)) {
			fmt.Printf(A.Marca+" %s\t\t\t", v.ModeloNome)
			fmt.Println()
		}
	}
}

func getAnos() ([]veiculoHelper.Ano, veiculo.VeiculoModel) {
	ano := Modelos[0]
	args := strings.Split(O, " ")
	p := false
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	for _, v := range Modelos {
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
					O = v.ModeloNome
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
	v := fv.GetVeiculoModeloBase()
	v.SetMarca(A)
	a, err := fv.GetAnos(&v, ano)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return a, v
}

func (c *Cli) PrintAnos() {
	an, _ := getAnos()
	fmt.Println(A.Marca + " " + O + ":")
	for _, vl := range an {
		fmt.Printf("-%s\n", vl.AnoComb)
	}
}

func (c *Cli) PrintVeiculo() {
	a, v := getAnos()
	ano := veiculoHelper.Ano{}
	if strings.Contains(strings.ToUpper(N), "ZERO") || strings.Contains(strings.ToUpper(N), "KM") || strings.Contains(strings.ToUpper(N), "32000") || strings.ToUpper(N) == "0" {
		N = "Zero KM"
	}
	for _, vl := range a {
		if strings.Contains(vl.AnoComb, N) {
			ano = vl
			break
		}
	}
	if (ano == veiculoHelper.Ano{}) {
		fmt.Printf("Nenhum %s encontrado com o ano %s\n", O, N)
		return
	}

	vr, err := fv.GetVeiculo(&v, ano)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Println(vr.ToString())
}
