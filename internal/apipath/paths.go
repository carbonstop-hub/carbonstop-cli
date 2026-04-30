package apipath

import "encoding/base64"

const (
	pingEncoded        = "L2Zvb3RwcmludDMvYWdlbnQvY2xpL3Bpbmc="
	whoamiEncoded      = "L2Zvb3RwcmludDMvYWdlbnQvY2xpL3dob2FtaQ=="
	echoEncoded        = "L2Zvb3RwcmludDMvYWdlbnQvY2xpL2VjaG8="
	productsEncoded    = "L2Zvb3RwcmludDMvcHJvZHVjdC9wcm9kdWN0TGlzdA=="
	productInfoEncoded = "L2Zvb3RwcmludDMvcHJvZHVjdC9pbmZv"
	accountsEncoded    = "L2Zvb3RwcmludDMvYWNjb3VudC9saXN0"
	accountViewEncoded = "L2Zvb3RwcmludDMvYWNjb3VudC9hY2NvdW50Vmlldw=="
	aiModelEncoded     = "L2Zvb3RwcmludDMvYWNjb3VudEVtaXNzaW9uL2FpTW9kZWxCeUNvbnRlbnQ="
	searchFactorEncoded = "L21hbmFnZW1lbnQvc3lzdGVtL3dlYnNpdGUvcXVlcnlGYWN0b3JMaXN0Q2xhdw=="
)

func decode(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(b)
}

func Ping() string        { return decode(pingEncoded) }
func Whoami() string      { return decode(whoamiEncoded) }
func Echo() string        { return decode(echoEncoded) }
func Products() string    { return decode(productsEncoded) }
func ProductInfo() string { return decode(productInfoEncoded) }
func Accounts() string    { return decode(accountsEncoded) }
func AccountView() string { return decode(accountViewEncoded) }
func AiModel() string     { return decode(aiModelEncoded) }
func SearchFactor() string { return decode(searchFactorEncoded) }
