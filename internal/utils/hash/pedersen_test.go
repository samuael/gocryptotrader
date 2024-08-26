package hash

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	path "github.com/thrasher-corp/gocryptotrader/internal/testing/utils"
)

func TestPedersen(t *testing.T) {
	rootPath, err := path.RootPathFromCWD()
	if err != nil {
		t.Fatal(err)
	}
	const defaultPedersenConfigsPath = "internal/utils/hash/pedersen_config/"
	loadConfig, err := LoadPedersenConfig(rootPath + "/" + defaultPedersenConfigsPath + strings.ToLower("apexpro") + ".json")
	if err != nil {
		t.Fatal(t, err)
	}

	tests := []struct {
		a, b string
		want string
	}{
		{
			"0x03d937c035c878245caf64531a5756109c53068da139362728feb561405371cb",
			"0x0208a0a10250e382e1e4bbe2880906c2791bf6275695e02fbbc6aeff9cd8b31a",
			"0x030e480bed5fe53fa909cc0f8c4d99b8f9f2c016be4c41e13a4848797979c662",
		},
		{
			"0x58f580910a6ca59b28927c08fe6c43e2e303ca384badc365795fc645d479d45",
			"0x78734f65a067be9bdb39de18434d71e79f7b6466a4b66bbd979ab9e7515fe0b",
			"0x68cc0b76cddd1dd4ed2301ada9b7c872b23875d5ff837b3a87993e0d9996b87",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestHash %d", i), func(t *testing.T) {
			a, ok := big.NewInt(0).SetString(tt.a, 0)
			if !ok {
				t.Errorf("expected no error but got ")
			}
			b, ok := big.NewInt(0).SetString(tt.b, 0)
			if !ok {
				t.Errorf("expected no error but got ")
			}

			want, ok := big.NewInt(0).SetString(tt.want, 0)
			if !ok {
				t.Errorf("expected no error but got ")
			}

			ans := loadConfig.PedersenHash(a.String(), b.String())
			if ans != want.Text(10) {
				t.Errorf("TestHash got %s, want %s", ans, tt.want)
			}
		})
	}
}
