package tasks

// ########################################### IMPORTS
import (
	"encoding/json"
	"log"
	"strings"
	"time"
	"context"
	"fmt"

	"github.com/parnurzeal/gorequest"

	"go.mongodb.org/mongo-driver/mongo" // MongoDB

	"github.com/pusher/pusher-http-go" // Pusher

	. "github.com/logrusorgru/aurora" // colors

	// ## local shits ##
	"../database"
	"../proxies"
	"../requests"
	"../structs"
	"../utils"
)

// ########################################### VARIABLES
func SNKRS(region string, mongoClient *mongo.Client, pusherClient *pusher.Client) {
	identifier := "snkrs"
	timeout := 5 * time.Second
	for {
		scrapeSNKRS(region, identifier, mongoClient, pusherClient)
		time.Sleep(timeout)
	}
}

func scrapeSNKRS(region string, identifier string, mongoClient *mongo.Client, pusherClient *pusher.Client) {

	// ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// ########################################### VARIABLES
	simple_url := "nike.com/" + region + "/launch"
	locale, currency, currencySymbol, cookies_string := "", "", "", ""
	if region == "us" {
		locale = "&country=US&language=en"
		currency = "USD"
		currencySymbol = "$"
		// cookies_string = "bm_sz=8DFBB70D4F871697FA2B264D79F6BAB9~YAAQzRbJFwkv9RlxAQAAzU26IgfV9K0t1Ha9blRbtuU8bJVFC9nDJMjPwUHN//JPlZ1MZn+g69ZwGve1UK+UKpuAuxi8u/ym4+fTcT7MlxWrK4p3bG+U143V/2MVumw2GuPyWCJX1H7GjLqCY06DSeLvI8FRzWKKOl56yc8zwFlGrR/3WpC7vtMSPzbsTQ==; _abck=E06BB66C2F3E5CB98B45B1CAAE171E74~-1~YAAQzRbJFwov9RlxAQAAzU26IgPA9G3fQ1J6Z3c61DAgvQ09asW6623TsT3u6H6NbztDUB4XkeXZFGLnyDeWO7Nc+TXtZXOsPiClQMVXW5ZfMTTpqsBT4k4BSEYcJFcs9xEmnfMe0YuMIees8ix4SC5mZ0UwAzv1c4s4/1haZt7tuY+66vaQFPMtx/DjG6OKR+w7aN4YvRK3ABrqeQE4oFBRsSbheVte6m5R5V/oQPSghuhKhBd4yl/npDdkKkg2CBpY7k5mzZAugRxbgFHvF32DSvKbivDRQMh5ZdUX9Wh4gN+gFmi2Pw==~-1~-1~-1"
		cookies_string = "s_ecid=MCMID%7C07334963115233906060070716914176151705; AMCV_F0935E09512D2C270A490D4D%40AdobeOrg=1994364360%7CMCMID%7C07334963115233906060070716914176151705%7CMCAID%7CNONE%7CMCOPTOUT-1576565792s%7CNONE%7CvVersion%7C3.4.0; AnalysisUserId=72.246.65.105.35761576681273401; geoloc=cc=US,rc=FL,tp=vhigh,tz=EST,la=26.1418,lo=-81.795; bm_sz=89531E6C7CEA4D7AEAA155A9FD803A81~YAAQaUH2SIoHiNhuAQAAOfiGGQacS1CMDLzxbewxi8KH2dtzcRMKzRVM0zg6ScpovV1NUVbPoRt5Jr7HI6Y9nPEq/ZBf12HNDj1Qe4Ix+Fmn2dI4ATDvEQV/9Q3sE0iSLPO0NQeHaej5ARvb0meNI2tFuTRb1NKOPBYfsgqnbjqfCv8vgEuUu5hs0Jk=; anonymousId=0EF7F9E542C8D6067CD51EB22D6CE94D; ak_bmsc=42BF752ABDF8FA5F83C028C805110E2A48F64169F80D0000393FFA5D911EE31A~pl2HDI9kllY6a9JeZ6fX3oqtJdvVG7sLwKmCZGsNVwdHULQLwLCcvHlB/qBpjYxq8Id/ns87jm/BhJ0O08Q85So2XHAB6nA4qMuLuyxwHXHys5jK1HDrtHwROXhbEbf7qtcBOziAci16DfdywC0xGlB/Xdpv8xbPzLDtJ2McAJfGKwd+vSUoFAHwuzZ9VRrZQQ0UKuNRqHyNSf5XFPFELnhtFHggaF0KmlxZlcFZ0lYvkLii8HEc+3C3e+GV48m7Rg; sq=3; NIKE_CART=b:c; siteCatalyst_sample=52; dreamcatcher_sample=40; neo_sample=10; guidU=87805623-8e43-4f2d-87aa-ff99faa9b76d; DAPROPS=\"sdevicePixelRatio:1|sdeviceAspectRatio:16/9|bcookieSupport:1\"; s_cc=true; AKA_A2=A; guidS=4e314ee1-c069-4bb8-d889-a84499415ad3; fs_uid=rs.fullstory.com`BM7A6`5820811341561856:6231879364739072/1608094591; guidSTimestamp=1576685825942|1576686224919; ppd=productgrid:cic:snkrs|nikecom>snkrs>discover feed>feed; bm_sv=D1F63919E30361D2B9CCE3738E8BA637~r4Ixp1Jnkqu1OGcHwQJ4PhuXXV8AkqlYAs/z1FkjJiUkoRvzM2zpW/fvMbCOygJLq8PFAxz7VvZ+V6xEMmTJGqSZ7NoYZmEVXkxbGB7ZgEsnoReLgkvE7/6rBpx+xzIxqa+igti0Q0I25mMqIxaZGI8CF4LT5YTu/9tI2WQCz58=; _abck=90C4D3BEC2D0FC5D3ABE9FFA64F44E53~-1~YAAQJ9VMaIQPjBBvAQAAxc/SGQP2g9/IukDDfwksZvI/ywXWFWaTfW0mdUMSXYkqKyubeAUiIzFwvvT2i3dWe59Y8atv6T4t1yxDXMMlnIlTf4GhR32j6swGeYx3UiQHGr/1CNdDg4plvSXv1BdpPAVw1lkpkbwsNm5e0fVKjm2w+ZBd//i9uhHSgT767PMs+DZUQfiNr630xUFiCGwa/moo90cgK/gZaD5f69wKgf9ZAlTKNf+fl1114iNELQmLJLAYkZOD4XT335elZHIX13+i0A+GjBgfDofcxJFLOQH/Xc5i8wiuy9Ow6FlURtxbXSW4QGiwmjYBjOVfPo0jY9UB~-1~-1~-1; CONSUMERCHOICE=us/en-us; NIKE_COMMERCE_COUNTRY=US; NIKE_COMMERCE_LANG_LOCALE=en_US; nike_locale=us/en_us; RT=\"dm=nike.com&si=1fb88bae-697e-4087-a9a4-a859da98b61d&ss=1576685821829&sl=8&tt=27768&obo=0&sh=1576686225091%3D8%3A0%3A27768%2C1576686207614%3D7%3A0%3A24966%2C1576686159135%3D6%3A0%3A22785%2C1576686118578%3D5%3A0%3A19700%2C1576686048066%3D4%3A0%3A15844&bcn=%2F%2F173c5b08.akstat.io%2F&ld=1576686225092&nu=https%3A%2F%2Fwww.nike.com%2Fus%2Flaunch%2F&cl=1576686249075&r=https%3A%2F%2Fwww.nike.com%2Fgb%2Flaunch%2F&ul=1576686249102\"; s_pers=%20c6%3Dno%2520value%7C1576688040738%3B%20c5%3Dno%2520value%7C1576688049155%3B"
	} else if region == "gb" {
		locale = "&country=GB&language=en-GB"
		currency = "GBP"
		currencySymbol = "€"
		cookies_string = "AnalysisUserId=2.18.66.134.10371576687780556; geoloc=cc=GB,rc=EN,tp=vhigh,tz=GMT,la=51.50,lo=-0.12; AKA_A2=A; bm_sz=F63DDCC69DD324F053333DFAC3DE53E0~YAAQhkISAskdXr5uAQAAzELqGQZdKNXTKjUyHCdMvKcHcgJ33+41nEoGWNFp7Z3mUorBuXF9c5h72j/c5ldNZuwQ/HSatsh3Fr383vbAgeCJyTLUTDLCePrjYcHFaGZsSRQZ1CcX63Dk0y8LEhjR5LYZVZxvONPT8lelzLygEptKUkvzQ2gDkD3DXk7FaA==; NIKE_COMMERCE_COUNTRY=GB; nike_locale=gb/en_gb; NIKE_COMMERCE_LANG_LOCALE=en_GB; s_ecid=MCMID%7C23182509664586892243582962843485349870; AMCVS_F0935E09512D2C270A490D4D%40AdobeOrg=1; AMCV_F0935E09512D2C270A490D4D%40AdobeOrg=1994364360%7CMCMID%7C23182509664586892243582962843485349870%7CMCAID%7CNONE%7CMCOPTOUT-1576694980s%7CNONE%7CvVersion%7C3.4.0; ak_bmsc=6593EBAE9565213BE09083B160FB68645C7A994C065E0000A958FA5D0F929828~plrXzhuFezJ+jbTe3G5QB1IfF5zUAepgbuqZ5a2xUe+sB9DxvwGbtan+AUBzxDxP8Y1kOqQ24vCWDEepUHa//QoX2Gn7TKxctWq3+h3NmilUVKgpfZCLXG3lxHAk0nON369xfiSxUgFwrNuItgliMJPs2KZWNgSiS/jrV78Spb7wkixKNwzRbDt7w/ap803cHbBUMNM76dUdYZBcn8o50FeLicL6ZMzPoWGB9onzt8VIjumn4TQ9Kq4fDAWiYIdHOA; CONSUMERCHOICE=gb/en_gb; anonymousId=749D761BD4C260B19F01522DACBE7683; siteCatalyst_sample=42; dreamcatcher_sample=93; neo_sample=30; NIKE_CART=b:c; guidS=dd3ad26e-62d5-411f-f171-6b2cae6ab3f0; guidSTimestamp=1576687785464|1576687785464; guidU=259d51d5-3b35-473b-8824-2c603e443bd5; s_pers=%20c6%3Dproductgrid%253Acic%253Asnkrs%7C1576689585592%3B%20c5%3Dnikecom%253Esnkrs%253Ediscover%2520feed%253Efeed%7C1576689585606%3B; s_cc=true; bm_sv=581CCA7505A8B15A818C426453A1606C~yokqQ5o5UN+IIdCnCyvp3i8ZNSdhd2o+XKAw9GT59O+rI9YVR3usdlEwMx85XjnypceBHRuIIBl0jS5ziVocOlWbc5/S8++dQMhTD4e73wksMn/d8CP0iavkSkvmkOqa+lHOgd7GOYXzJhbORqmSgw==; _abck=10D922113D5B1E29C0869E5B314789A2~-1~YAAQhkISAo0gXr5uAQAAO2vqGQPFuUxojJzG+xONBz4MABfZbLgYxh9kIzh+doS1l1jL4Q/KINk2wo4lNCi1ZWsDI890yeOMeXniwyN2KF7qutOmAoM2Ls/xw15JZTjMH5WsujDUJOxzL/3QWCImoAzuaCLbA1Iv3Y8chfPQI/sktbKetyaJlcgfoxh9PeUZ/XndZSdjj0wKf1P8NxAXSFQPN55xwxB1zPyv5Ab2+JdkY9J0gU6P8Z8xT6aJlUKN3UFO7RMZdksT/wOhssTtLvcQJk3kGVukrSOKTxI0SCtICfQKK4/QNKsaaemIW4QWCw2Vfg4vmGYxIPfy9+4UjXLf~-1~||-1~-1; ppd=productgrid:cic:snkrs|nikecom>snkrs>discover feed>feed; RT=\"sl=2&ss=1576687779041&tt=8772&obo=1&bcn=%2F%2F6852bd05.akstat.io%2F&sh=1576687791120%3D2%3A1%3A8772%2C1576687782264%3D1%3A1%3A0&dm=nike.com&si=eb726c69-d9c1-444c-8818-426ea081ea04&ld=1576687791124\""
	} else if region == "jp" {
		locale = "&country=JP&language=ja"
		currency = "JPY"
		currencySymbol = "¥"
		cookies_string = "AnalysisUserId=23.15.1.23.82181576687694664; geoloc=cc=JP,rc=13,tp=vhigh,tz=GMT+9,la=35.69,lo=139.75; AKA_A2=A; bm_sz=F81C17A7E58A57006FD13FD033867AA5~YAAQFwEPF3uSc39uAQAASPPoGQagKNjBFs/5zl0ssc5UNoD8nctjNB4qOREu/tKjL/rDCa62adunOcaO986YB87h/1vECNxkVmDWsaJt+MkVvxY0n1KIk3OeKs/iXNBpIaIMRCxtOF/K7MBlT7wH9FvRONbuTlu5u6BDnI/MKzyZkE3woRsF6i9rsayEnQ==; NIKE_COMMERCE_COUNTRY=JP; NIKE_COMMERCE_LANG_LOCALE=ja_JP; nike_locale=jp/ja_jp; anonymousId=05762C1E79045EF37AD35C42DD2FD83A; s_ecid=MCMID%7C55756562435475254138017663656487345888; AMCVS_F0935E09512D2C270A490D4D%40AdobeOrg=1; AMCV_F0935E09512D2C270A490D4D%40AdobeOrg=1994364360%7CMCMID%7C55756562435475254138017663656487345888%7CMCAID%7CNONE%7CMCOPTOUT-1576694895s%7CNONE%7CvVersion%7C3.4.0; CONSUMERCHOICE=jp/ja_jp; siteCatalyst_sample=9; dreamcatcher_sample=29; neo_sample=3; NIKE_CART=b:c; guidS=0e1b4609-7931-4ed7-f79f-a4436e0fb84e; guidU=222974c8-0fed-4299-9164-979f18fa6ac4; s_cc=true; bm_mi=5CB49B8F2D4DB11EF2A2CF9791A63ECD~tdGqzbVEWpXMoG/ofi7gi4IhN0442d0sFMDiGkUlC9PjuybSNPpfo5JCm1AXysR9Tb+983j/2JbNr2o5gX9pZVi+2r5VPIvRgX4s2aXbJS3gPrbI9wQU65xtiRNfM6XUa3HBJS6keijWpglWfZNOZaaS0dApbJkJWphpFB9eVIg94od9aN85PVCe433e/+MvFk/wwgPTvgvU7L5AfAgpVB8jUnOJM5GXC8/06Et9psGLhVi9zrd81u0TdaAR99uZr9LNOAF2zsfmBZnZ2IaI5g==; _abck=1AE7D37C5D650C953672AA13DC4785AC~-1~YAAQFwEPF/aSc39uAQAAPCbpGQN/Z/lO6XP6H3KDsvm5QOqe/+JFTyibgO0fpB/oa2ynV4+6JZOl/MWFGHwjJ/tIV1EveLqHil3txDS1cy3xces5bhpx1dOIxClUDmjc9kblST+sNhPZMiwLcAzYihFdUZbslkm8+Qu0mLpp4cfrFjyIxDtAXUfy5ztn6OJkt20NWU2DFbF9K9f1HgQdUm/CSo6rnF2ZRyF7MZgvIlWcCY71whrbvSR8uhMyRzTeratVK9S54JYJAUr2/OmANiyqxJRLRGkeoydFCxxHgTvS5keQSTcIAE2kDt84l+dWbKDR1YF2I1rh0+qvlmJkfY93~-1~-1~-1; ak_bmsc=2D14AC762C8FB0B927AEB4385CB05D787567B614D61800005258FA5D32FC3E06~pllGDN0Z4lIU9LeFE4aZaXjakuxe3uuv63IPu8JbJUF6sEqN8smp6Pd/h0fmYo0gCGZrJ9u6YoKdN69dnM1ZEmazOufnwLB5jMVJD2mBXQ7ZXsyehb7HefI4IPmkbi/QAiWC6Y7rtHzObjfvrHFtTGpNV90suFhy3LXoTDC8mtk6eUydgr74JDNk2JVcZdG6wKg++7Y/otRUKMpIM5ihDWDlL/+ohnlUA9bN7wGia89O+W3AozMz30Zav3Q88t6Vlj; guidSTimestamp=1576687699369|1576687705517; s_pers=%20c6%3Dproductgrid%253Acic%253Asnkrs%7C1576689505527%3B%20c5%3Dnikecom%253Esnkrs%253Ediscover%2520feed%253Efeed%7C1576689505532%3B; RT=\"dm=nike.com&si=e30a43f6-37b5-4b20-8186-9a6db0d765fa&ss=1576687693809&sl=3&tt=5964&obo=1&sh=1576687705802%3D3%3A1%3A5964%2C1576687701149%3D2%3A1%3A4437%2C1576687695686%3D1%3A1%3A0&bcn=%2F%2F684d0d36.akstat.io%2F&ld=1576687705803\"; bm_sv=2142647FCC5D9F65C13CD98D59ED7364~mtI7g/GX9lPLaiXgbRfLKlHHzCPZqIVneUpNP3xO2DS0qmC5n3fcKMMLcIMyf/29aWIFSZHm+QaX2hTtdZP8O0Yrq6QNh44etQxl+wyWN16sSEOUKKvklddB3OZpgisF2ZoQ+wTa1K5rhPGKCMbEAA==; ppd=productgrid:cic:snkrs|nikecom>snkrs>discover feed>feed"
	} else if region == "cn" {
		locale = "&country=CN&language=zh-Hans"
		currency = "CNY"
		currencySymbol = "¥"
		cookies_string = "s_ecid=MCMID%7C07334963115233906060070716914176151705; AMCV_F0935E09512D2C270A490D4D%40AdobeOrg=1994364360%7CMCMID%7C07334963115233906060070716914176151705%7CMCAID%7CNONE%7CMCOPTOUT-1576565792s%7CNONE%7CvVersion%7C3.4.0; fs_uid=rs.fullstory.com`BM7A6`5820811341561856:4861145312755712/1608094591; AnalysisUserId=72.246.65.105.35761576681273401; geoloc=cc=US,rc=FL,tp=vhigh,tz=EST,la=26.1418,lo=-81.795; bm_sz=89531E6C7CEA4D7AEAA155A9FD803A81~YAAQaUH2SIoHiNhuAQAAOfiGGQacS1CMDLzxbewxi8KH2dtzcRMKzRVM0zg6ScpovV1NUVbPoRt5Jr7HI6Y9nPEq/ZBf12HNDj1Qe4Ix+Fmn2dI4ATDvEQV/9Q3sE0iSLPO0NQeHaej5ARvb0meNI2tFuTRb1NKOPBYfsgqnbjqfCv8vgEuUu5hs0Jk=; anonymousId=0EF7F9E542C8D6067CD51EB22D6CE94D; ak_bmsc=42BF752ABDF8FA5F83C028C805110E2A48F64169F80D0000393FFA5D911EE31A~pl2HDI9kllY6a9JeZ6fX3oqtJdvVG7sLwKmCZGsNVwdHULQLwLCcvHlB/qBpjYxq8Id/ns87jm/BhJ0O08Q85So2XHAB6nA4qMuLuyxwHXHys5jK1HDrtHwROXhbEbf7qtcBOziAci16DfdywC0xGlB/Xdpv8xbPzLDtJ2McAJfGKwd+vSUoFAHwuzZ9VRrZQQ0UKuNRqHyNSf5XFPFELnhtFHggaF0KmlxZlcFZ0lYvkLii8HEc+3C3e+GV48m7Rg; sq=3; NIKE_CART=b:c; siteCatalyst_sample=52; dreamcatcher_sample=40; neo_sample=10; guidU=87805623-8e43-4f2d-87aa-ff99faa9b76d; DAPROPS=\"sdevicePixelRatio:1|sdeviceAspectRatio:16/9|bcookieSupport:1\"; s_cc=true; AKA_A2=A; guidS=4e314ee1-c069-4bb8-d889-a84499415ad3; guidSTimestamp=1576685825942|1576686048628; ppd=productgrid:cic:snkrs|nikecom>snkrs>discover feed>feed; bm_sv=D1F63919E30361D2B9CCE3738E8BA637~r4Ixp1Jnkqu1OGcHwQJ4PhuXXV8AkqlYAs/z1FkjJiUkoRvzM2zpW/fvMbCOygJLq8PFAxz7VvZ+V6xEMmTJGqSZ7NoYZmEVXkxbGB7ZgEtJDHCcSAxt7CFBnaIA/GuEV60bgNIy5qY+AJqc7rNsog==; _abck=90C4D3BEC2D0FC5D3ABE9FFA64F44E53~-1~YAAQJ9VMaKkFjBBvAQAAa9/QGQNd8h9AVeEz5IX6BqnL+E0UMDws4yjZNGl9j8N9s/Z//GD4n2R2ARZxS53kFYmjSxRshWvV7hZUtGmkSihL+GKP3n41JA1VoytD0ps0ekt4L2dkqVYi3lWxwCJct3//fQbHJzN/tuov85jQtYinFQdiC9eQiltusU5LAHho1SKW1y5Ij50bU5GOtl2a5rzKSxImu9PqvS85/NRM24gEhxOuBMf0a+hpwCDQFLFb3iU5XHrWz2rCrO5kdRDsFomqGqWV8agzCD4tcG/UqO2l1GO5osLHuC8Zn8+OHIqqGezCGZu4u1/f/MKv96iCd2lq~-1~-1~-1; CONSUMERCHOICE=cn/zh-cn; NIKE_COMMERCE_COUNTRY=CN; NIKE_COMMERCE_LANG_LOCALE=zh_CN; nike_locale=cn/zh_cn; RT=\"dm=nike.com&si=1fb88bae-697e-4087-a9a4-a859da98b61d&ss=1576685821829&sl=4&tt=15844&obo=0&sh=1576686048066%3D4%3A0%3A15844%2C1576685990429%3D3%3A0%3A11600%2C1576685836520%3D2%3A0%3A6466%2C1576685826741%3D1%3A0%3A3863&bcn=%2F%2F173c5b08.akstat.io%2F&ld=1576686048067&nu=https%3A%2F%2Fwww.nike.com%2Fcn%2Flaunch%2F&cl=1576686113608&r=https%3A%2F%2Fwww.nike.com%2Fsg%2Flaunch%2F&ul=1576686113636\"; s_pers=%20c5%3Dno%2520value%7C1576687913677%3B%20c6%3Dno%2520value%7C1576687913680%3B"
	} else if region == "sg" {
		locale = "&country=SG&language=en-GB"
		currency = "SGD"
		currencySymbol = "$"
		cookies_string = "AnalysisUserId=23.194.187.204.194671576688050611; geoloc=cc=SG,rc=,tp=vhigh,tz=GMT+8,la=1.29,lo=103.86; AKA_A2=A; bm_sz=47B977728F02AB8537151F8EB2FD9356~YAAQzLvCF+dymhBvAQAAs2HuGQboLwvoEAj77UtsKJUBfh/oEgblB8JFYwhMayIgDTySgHHGLBq/aGNMSusYEF/QEpI5YUiED4lNcyuaVw5hZaWHtnhC/SQ0s1NQSfeB8bGJSOlqOhgx+xgJN3oOKPzB++v2cRMmMkytEpkAreT60cbBTxVlmSD7hzCGng==; NIKE_COMMERCE_COUNTRY=SG; NIKE_COMMERCE_LANG_LOCALE=en_GB; nike_locale=sg/en_gb; anonymousId=6772F5081C241E3BAF474830669A47C2; ak_bmsc=704F5885B2263FBE352EFEC90B185A6F17C2BBCC0B4C0000B259FA5DF0F5780F~plGwKfKRXlbwf/tuFXgT9edbluTk2Fdu5NI0H0iHlhSgHlc9YjBOgDF7uVmCsW/tj8xngcCDGM7Tku11tN23cAKUdxmjT5N4lvvmxLkIojF5uUDZnEV6sEi9COKJoUBDpLKnTQn+GzumqEncttZchFy/lJqZOeLqX75kWU92T6EBNhqMHLFnxBBOmziDjnPhJskDG7neL78aq/zXlBpQ1+zKtbe8MBbhiBB3DzaFVJD3VqRK+FbrjhQdL2hRgpCwMe; s_ecid=MCMID%7C35342733638439236958331368682982867022; AMCVS_F0935E09512D2C270A490D4D%40AdobeOrg=1; AMCV_F0935E09512D2C270A490D4D%40AdobeOrg=1994364360%7CMCMID%7C35342733638439236958331368682982867022%7CMCAID%7CNONE%7CMCOPTOUT-1576695252s%7CNONE%7CvVersion%7C3.4.0; bounceClientVisit2441v=N4IgNgDiBcIBYBcEQM4FIDMBBNAmAYnvgO6kB0AdgJYDWApmQMYD2AtkSgOZEgA0IAJxggQAXyA; NIKE_CART=b:c; bc_nike_singapore_triggermail=%7B%22distinct_id%22%3A%20%2216f19ee81a22da-0eb4c845bb9de6-6701b35-1fa400-16f19ee81a3259%22%7D; _gcl_au=1.1.624533255.1576688060; _fbp=fb.1.1576688061193.391923856; CONSUMERCHOICE=sg/en_gb; sq=3; siteCatalyst_sample=61; dreamcatcher_sample=68; neo_sample=14; bm_sv=A46B2286252C05C3B88C1F584AC3801B~TaSpexdy1LYGGyzG5WVhm69grbXGa41roTZuN+tt4pg9tL4DEXpIGmFeXGT7VyrE25bwK+e+Xv3qtQtqjqQ7M4xwgHLQ/gBA7WUOOujtmhFeCZGpVDRnNHcyz/KyesqquKjLl27UuEYo0d5+ec1O0Q==; s_pers=%20c6%3Dproductgrid%253Acic%253Asnkrs%7C1576689899785%3B%20c5%3Dnikecom%253Esnkrs%253Ediscover%2520feed%253Efeed%7C1576689899791%3B; s_cc=true; guidS=6dd967b7-1f70-4357-dc49-d09912e63d90; guidSTimestamp=1576688099836|1576688099836; guidU=5d9aa9c5-5b95-446c-eec0-6b33e0029d05; ppd=productgrid:cic:snkrs|nikecom>snkrs>discover feed>feed; RT=\"sl=2&ss=1576688046087&tt=17657&obo=0&bcn=%2F%2F684fc53c.akstat.io%2F&sh=1576688100284%3D2%3A0%3A17657%2C1576688056633%3D1%3A0%3A10525&dm=nike.com&si=1edcdccf-8a29-40b3-b105-bc7ab642f793&ld=1576688100285\"; _abck=4CDCBE362B3F5AB91A4629EA636A54A9~-1~YAAQzLvCF1h6mhBvAQAAKjrvGQP4fXTGYQYdl0wjgo41X/T5bay1rGqdQBrCNbWzuqkLxN2EdIkMswHTT2WyGBlsEqcgE4hhMkORSbp/j9JlqFIyyfoYtSuGFoobkUZ6ivFr8XpMEh78s9DGHuZ2Mwu7WPQ9Pz0+RkiHaJEON8t5P1Z2kLGfgqKciawNhaOk9NQQEctNwipfViWAhfP5tyrZYNOGXCFtpvnCgvjEPiJy7M71i8KFZnevFgvXftAiHezaHnPfBRiaapvaWWUaAST6d7O1WqniSfDO8thXw4VfSo07ESjpFtSMkAh+lAOi4V/V0W6zVeQWgZQ3Of8CcfGP~-1~-1~-1"
	} else if region == "ca" {
		locale = "&country=CA&language=en-GB"
		currency = "CAD"
		currencySymbol = "$"
		cookies_string = "s_ecid=MCMID%7C07334963115233906060070716914176151705; AMCV_F0935E09512D2C270A490D4D%40AdobeOrg=1994364360%7CMCMID%7C07334963115233906060070716914176151705%7CMCAID%7CNONE%7CMCOPTOUT-1576565792s%7CNONE%7CvVersion%7C3.4.0; fs_uid=rs.fullstory.com`BM7A6`5820811341561856:4861145312755712/1608094591; AnalysisUserId=72.246.65.105.35761576681273401; geoloc=cc=US,rc=FL,tp=vhigh,tz=EST,la=26.1418,lo=-81.795; bm_sz=89531E6C7CEA4D7AEAA155A9FD803A81~YAAQaUH2SIoHiNhuAQAAOfiGGQacS1CMDLzxbewxi8KH2dtzcRMKzRVM0zg6ScpovV1NUVbPoRt5Jr7HI6Y9nPEq/ZBf12HNDj1Qe4Ix+Fmn2dI4ATDvEQV/9Q3sE0iSLPO0NQeHaej5ARvb0meNI2tFuTRb1NKOPBYfsgqnbjqfCv8vgEuUu5hs0Jk=; NIKE_COMMERCE_LANG_LOCALE=en_GB; anonymousId=0EF7F9E542C8D6067CD51EB22D6CE94D; ak_bmsc=42BF752ABDF8FA5F83C028C805110E2A48F64169F80D0000393FFA5D911EE31A~pl2HDI9kllY6a9JeZ6fX3oqtJdvVG7sLwKmCZGsNVwdHULQLwLCcvHlB/qBpjYxq8Id/ns87jm/BhJ0O08Q85So2XHAB6nA4qMuLuyxwHXHys5jK1HDrtHwROXhbEbf7qtcBOziAci16DfdywC0xGlB/Xdpv8xbPzLDtJ2McAJfGKwd+vSUoFAHwuzZ9VRrZQQ0UKuNRqHyNSf5XFPFELnhtFHggaF0KmlxZlcFZ0lYvkLii8HEc+3C3e+GV48m7Rg; sq=3; NIKE_CART=b:c; siteCatalyst_sample=52; dreamcatcher_sample=40; neo_sample=10; guidU=87805623-8e43-4f2d-87aa-ff99faa9b76d; DAPROPS=\"sdevicePixelRatio:1|sdeviceAspectRatio:16/9|bcookieSupport:1\"; s_cc=true; AKA_A2=A; guidS=4e314ee1-c069-4bb8-d889-a84499415ad3; guidSTimestamp=1576685825942|1576685836056; ppd=productgrid:cic:snkrs|nikecom>snkrs>discover feed>feed; CONSUMERCHOICE=ca/en-gb; NIKE_COMMERCE_COUNTRY=CA; nike_locale=ca/en_gb; s_pers=%20c6%3Dno%2520value%7C1576687764043%3B%20c5%3Dno%2520value%7C1576687764049%3B; _abck=90C4D3BEC2D0FC5D3ABE9FFA64F44E53~-1~YAAQJ9VMaIX7ixBvAQAA1svOGQOQXfpVd7t6CzqVDjAAsqqpgoRc5uciagb68iBY4AcfWrv4LiWQUmuBhCXcs17chPzrlbvtj2vev0f/8bo/wLixhKtjQbWTD6gw2tFEDbxB5lOMYeqQrd5C+JEkf6HXCuS9dZpXXhKUnz0Qq39Qio9CAHs29SyiJ19+3Tf70YAf/96NyqlOrRgtGTwR08YgtdEWs6S+IhnrWygcQJTb+PZEEIzpVhZA3nvfvgpNMg3L7uWSR5/jtwiE+U08bLRfFi0wcIEEXlbrsYNvCTx06ZF5aMjvPJLtql5t3h6+mRW+R1J2SNnTITdew3Z7HPKH~-1~-1~-1; bm_sv=D1F63919E30361D2B9CCE3738E8BA637~r4Ixp1Jnkqu1OGcHwQJ4PhuXXV8AkqlYAs/z1FkjJiUkoRvzM2zpW/fvMbCOygJLq8PFAxz7VvZ+V6xEMmTJGqSZ7NoYZmEVXkxbGB7ZgEuB7xZ8BryCTSg0qUWRXixR+tGDFjkCUWey5kSgMTF3G5xbEBVy/xNt2b4gNacFFwM=; RT=\"dm=nike.com&si=1fb88bae-697e-4087-a9a4-a859da98b61d&ss=1576685821829&sl=2&tt=6466&obo=0&sh=1576685836520%3D2%3A0%3A6466%2C1576685826741%3D1%3A0%3A3863&bcn=%2F%2F173c5b09.akstat.io%2F&r=https%3A%2F%2Fwww.nike.com%2Fca%2Flaunch&ul=1576685985214\""
	} else if region == "de" {
		locale = "&country=DE&language=de"
		currency = "EUR"
		currencySymbol = "€"
		cookies_string = "bm_sz=BAC579A6E0819C9819272CB0B50EDC8B~YAAQxz1raBDWmsluAQAAKBrcGQbJELYgdDeCWHkahnwvMxRqPfDoOnK/hVCGdlu7YGi1K+b0PIajXZAJaG3djW/HJ9OZKePffW7hhxZVC6CCE/QDX3DkSKcOfvFfMQlMIobjvDcx0+c6aGWMxDKOZqD5W3r4JybBqfPa+AGssrEotlxsnORrUs/1/jMvpw==; geoloc=cc=US,rc=FL,tp=vhigh,tz=EST,la=26.1418,lo=-81.795; anonymousId=SCRX491332C442FD123B124DBF973AD02E4D; AKA_A2=A; AnalysisUserId=72.246.65.105.35761576686967927; bm_mi=D1F54CDD9F429032641424DC44B82AE0~s1vKvFBNrJ5AWsr/kFrtXZwzyTlE53BZcXtMi6sEkA5LjbnhZyGa144ZN8EBfYIJST33djRlBTRQIYwMMOIXJ8CNz4UMW+gp8Mx6FFL+MaT9a3ZNQnjruilv9OK0w+ZLd9DDkKI/AzfNbTsgfWQmLP2ieU17M5nw8FUi+hLslguG1SBVwLHrQXAxc0iE793BKr2FLFhISQmeF8w2c+Eo4CeYBMx/ObZx1bgGGhmUV7RN/95BQOMStRGlCACxZIcupQODW4r3RsxHgBzEQ+N7C0bzhe6AAgHUn7P5ncDO/3dYkZB1LoylqZzDOzCmOC0H5TpApL9H+XpUTm4z3iuSHw==; s_ecid=MCMID%7C12349456556682883831851494551393794799; AMCVS_F0935E09512D2C270A490D4D%40AdobeOrg=1; AMCV_F0935E09512D2C270A490D4D%40AdobeOrg=1994364360%7CMCMID%7C12349456556682883831851494551393794799%7CMCAID%7CNONE%7CMCOPTOUT-1576694166s%7CNONE%7CvVersion%7C3.4.0; ak_bmsc=B9C16CC9C945B46237F638D7D545824F48F64169F80D00007755FA5D5AF9D171~plCK8VWFdxd7GcgPdFrCmObyzrEikQB/K5gQD9QTo/3Ch0LelCcnPY94MacRhOsqSsEoFfMN885h9zqyZoD4ewhpV4RWfjPSvbj6yKlvnpAt1CGORa0t9CNOHl5kluQ52oGNUcBL6k+uUQFa2mWYDXg56q0Z+uMA0Wo0HCiUbBYOYkYyenuDFrBfFrtNuQwGSsY29UTOAnPzM/FxejCzSQWBEq4LBFVvnxfkwRwq5J4OFVo5iI10tRmuXbgZTPdcVE; sq=3; NIKE_CART=b:c; siteCatalyst_sample=66; dreamcatcher_sample=91; neo_sample=43; s_cc=true; guidS=54c458e5-394e-439f-bd08-eac2fd5e158f; guidU=ab9d9b33-2f39-4656-841f-899b3f71a5e5; NIKE_COMMERCE_COUNTRY=DE; nike_locale=de/de_de; NIKE_COMMERCE_LANG_LOCALE=de_DE; _abck=7F1EDB9E9AACBBF8A6EF05E9D3FC9834~-1~YAAQOIQUAhmYdQRvAQAAcO3kGQMkYElTPrTylnWibBhHBMQAM9/x1NbUlozrWfslSZtcepwkeeaDm+OvbEor/P4a41DoAXDoTSfSKiCivBF7LYyrTl+jBAuiTJ/3xIDh2xoDVd9VxumTZPfuhMzDkzANuep+SMeXkqlOI/hG71bZgcHu+9x5CCqKMLJ1s51NXJrY1DQ9RWtDX3kyHf0fnmieYv50gBYdl4Uch6GT94klGz8mYa3FNJ+OIH326u/kn6hDdqXmMR6Yah45xlLySsZ5nxQUEo0tyS+V9+NQTthXU8ophgNc3js/nlzvZcm04UkrXnOcRTfnLiUBFYgu8GC2~-1~-1~-1; CONSUMERCHOICE=de/de_de; RT=\"dm=nike.com&si=1fb88bae-697e-4087-a9a4-a859da98b61d&ss=1576685821829&sl=21&tt=72533&obo=0&sh=1576687429708%3D21%3A0%3A72533%2C1576687426463%3D20%3A0%3A70756%2C1576687103008%3D19%3A0%3A60695%2C1576687081673%3D18%3A0%3A57290%2C1576687040189%3D17%3A0%3A55755&bcn=%2F%2F173c5b08.akstat.io%2F&ld=1576687429708\"; guidSTimestamp=1576686968048|1576687429834; s_pers=%20c6%3Dproductgrid%253Acic%253Asnkrs%7C1576689229970%3B%20c5%3Dnikecom%253Esnkrs%253Ediscover%2520feed%253Efeed%7C1576689229976%3B; bm_sv=53487DBD8BB3F35FBAC95B2FAF77D1FF~r4Ixp1Jnkqu1OGcHwQJ4PvZlP10TFCxI/LM+sAhwOpzYi5lYelSGoRS4wiE3eCkpatDcfTLfJxLMfFbu0I9aMuxolOdWae0iGSaTXZnF7APE2Xh1MMkEHymVNHgNnfhCZ3HwH225hMQscI+OjUO5bHWFw02RnR9e2C9Pxr7YBQY=; ppd=productgrid:cic:snkrs|nikecom>snkrs>discover feed>feed"
	} else if region == "ru" {
		locale = "&country=RU&language=ru"
		currency = "RUB"
		currencySymbol = "руб"
		cookies_string = "bm_sz=BAC579A6E0819C9819272CB0B50EDC8B~YAAQxz1raBDWmsluAQAAKBrcGQbJELYgdDeCWHkahnwvMxRqPfDoOnK/hVCGdlu7YGi1K+b0PIajXZAJaG3djW/HJ9OZKePffW7hhxZVC6CCE/QDX3DkSKcOfvFfMQlMIobjvDcx0+c6aGWMxDKOZqD5W3r4JybBqfPa+AGssrEotlxsnORrUs/1/jMvpw==; geoloc=cc=US,rc=FL,tp=vhigh,tz=EST,la=26.1418,lo=-81.795; DAPROPS=\"sdevicePixelRatio:1|sdeviceAspectRatio:16/9|bcookieSupport:1\"; anonymousId=SCRX491332C442FD123B124DBF973AD02E4D; NIKE_COMMERCE_COUNTRY=RU; nike_locale=ru/ru_ru; NIKE_COMMERCE_LANG_LOCALE=ru_RU; AKA_A2=A; AnalysisUserId=72.246.65.105.35761576686967927; bm_mi=D1F54CDD9F429032641424DC44B82AE0~s1vKvFBNrJ5AWsr/kFrtXZwzyTlE53BZcXtMi6sEkA5LjbnhZyGa144ZN8EBfYIJST33djRlBTRQIYwMMOIXJ8CNz4UMW+gp8Mx6FFL+MaT9a3ZNQnjruilv9OK0w+ZLd9DDkKI/AzfNbTsgfWQmLP2ieU17M5nw8FUi+hLslguG1SBVwLHrQXAxc0iE793BKr2FLFhISQmeF8w2c+Eo4CeYBMx/ObZx1bgGGhmUV7RN/95BQOMStRGlCACxZIcupQODW4r3RsxHgBzEQ+N7C0bzhe6AAgHUn7P5ncDO/3dYkZB1LoylqZzDOzCmOC0H5TpApL9H+XpUTm4z3iuSHw==; CONSUMERCHOICE=ru/ru_ru; s_ecid=MCMID%7C12349456556682883831851494551393794799; AMCV_F0935E09512D2C270A490D4D%40AdobeOrg=1994364360%7CMCMID%7C12349456556682883831851494551393794799%7CMCAID%7CNONE%7CMCOPTOUT-1576694166s%7CNONE%7CvVersion%7C3.4.0; AMCVS_F0935E09512D2C270A490D4D%40AdobeOrg=1; ak_bmsc=B9C16CC9C945B46237F638D7D545824F48F64169F80D00007755FA5D5AF9D171~plCK8VWFdxd7GcgPdFrCmObyzrEikQB/K5gQD9QTo/3Ch0LelCcnPY94MacRhOsqSsEoFfMN885h9zqyZoD4ewhpV4RWfjPSvbj6yKlvnpAt1CGORa0t9CNOHl5kluQ52oGNUcBL6k+uUQFa2mWYDXg56q0Z+uMA0Wo0HCiUbBYOYkYyenuDFrBfFrtNuQwGSsY29UTOAnPzM/FxejCzSQWBEq4LBFVvnxfkwRwq5J4OFVo5iI10tRmuXbgZTPdcVE; sq=3; NIKE_CART=b:c; siteCatalyst_sample=66; dreamcatcher_sample=91; neo_sample=43; s_pers=%20c6%3Dpdp%253Acic%253Asnkrs%7C1576688767959%3B%20c5%3Dnikecom%253Epdp%253Asnkrs%253ENike%2520Air%2520Max%252090%7C1576688767965%3B; s_cc=true; guidS=54c458e5-394e-439f-bd08-eac2fd5e158f; guidSTimestamp=1576686968048|1576686968048; guidU=ab9d9b33-2f39-4656-841f-899b3f71a5e5; _abck=7F1EDB9E9AACBBF8A6EF05E9D3FC9834~-1~YAAQaUH2SPfUidhuAQAAp+rdGQNefKONEt5HETsFKsRQTcqEHjwe2r3JkYODo5BCq7Subd+/E8/eaROxh4mbamQ+KakseLl5DF3+i4YL7ZWKFm6QdNTzax/+COD1wfxgBxhbTtkR6o2178qH5gK+8hCvvAtD7wl+tBRp/faNzGUW5lC+Kr41S8cAek3rJBHXnpbx3uc+gwly5CEiylv6ylAjzei8nvWVf1ie8gJk5G9OHP5hPaebp35Ixbs3/3rD5Ub7o07ktGaYvRcWrcEEG/byE84elDgh8UJwG2cEVH32Df0MGniANR8XAsTSeQQZE/dIK0gkcUWwYGnj39qD1rgs~-1~-1~-1; ppd=feed|snkrs>feed; RT=\"dm=nike.com&si=1fb88bae-697e-4087-a9a4-a859da98b61d&ss=1576685821829&sl=14&tt=43505&obo=0&sh=1576686970455%3D14%3A0%3A43505%2C1576686968262%3D13%3A0%3A42530%2C1576686943713%3D12%3A0%3A39709%2C1576686697864%3D11%3A0%3A38972%2C1576686385253%3D10%3A0%3A35213&bcn=%2F%2F173c5b08.akstat.io%2F&ld=1576686970456&nu=https%3A%2F%2Fwww.nike.com%2Fru%2Flaunch%2F&cl=1576686968439&r=https%3A%2F%2Fwww.nike.com%2Fru%2Flaunch%2F&ul=1576686971574\"; bm_sv=53487DBD8BB3F35FBAC95B2FAF77D1FF~r4Ixp1Jnkqu1OGcHwQJ4PvZlP10TFCxI/LM+sAhwOpzYi5lYelSGoRS4wiE3eCkpatDcfTLfJxLMfFbu0I9aMuxolOdWae0iGSaTXZnF7AMIW2ORko6cZwSyqAtlrEPlrRbkKTnMeHuStuj2j5eMmg=="
	} else if region == "it" {
		locale = "&country=IT&language=it"
		currency = "EUR"
		currencySymbol = "€"
		cookies_string = "AnalysisUserId=92.123.181.45.151671576688423411; geoloc=cc=IT,rc=,tp=vhigh,tz=GMT+1,la=45.47,lo=9.20; AKA_A2=A; bm_sz=480450096FB9BCB33C918C469DC0FACE~YAAQLbV7XD9wCQFvAQAA9BH0GQYp7cdFO8WzWYLNEMSUKuCatObzKH/9Yp2WUTDcyqFqkth3xTinkg9/H4Hnn/dCE/9F852kDHT0KJOdUil+p56zhFRIs+RhoMaYQOgjt4icSb68ufKOgP0TlbRLjV+3+IVQniTfW3NsGf/03R9OhZjM/1FJQpIjSBfyPA==; nike_locale=it/it_it; NIKE_COMMERCE_COUNTRY=IT; NIKE_COMMERCE_LANG_LOCALE=it_IT; s_ecid=MCMID%7C07401869167211151500553801574789218610; AMCVS_F0935E09512D2C270A490D4D%40AdobeOrg=1; AMCV_F0935E09512D2C270A490D4D%40AdobeOrg=1994364360%7CMCMID%7C07401869167211151500553801574789218610%7CMCAID%7CNONE%7CMCOPTOUT-1576695624s%7CNONE%7CvVersion%7C3.4.0; bm_mi=11194787B0690B56E560171BEBB6F86D~cF/UlRJajAhieUPd/eW7d1B/5HDAvzPdRy2CrCeM8Du3DzlDwGyjQNr0qkUcRXqzhxG5zffSzRKrqVmTaCltS6s5vEuvkA2sv6o4Rm6LGN4LmVwUgvKEhT0w58IanAV5eGDjtkxZgwtEUBghWneN3adWlKTlz/SRASs+SNQf+Y4zGdPo0N229s9BE0cfDm3SHUk8F9Z8sV7Q4EyqQNNxzI724D33caj1ehPcxGqRHjkX6Hz4sv6BkjHFkEkoSLJa4cwtNxLdTeWIu8Hp7FcJ5A==; CONSUMERCHOICE=it/it_it; anonymousId=1E28C663DF1D795197BD1BC1F8BCBD46; ak_bmsc=0F014727B909E1E37847FCA40F7645B0C316C857A63000002B5BFA5D00C0444B~pll48U3CDSjDyWRifeieTVKKV+vTXu5W5VbZCRDbuBDa9gruaHiumYyLnh1dBYHO1/p2ScFYV6hf9kgBnmIYi0jXwDRonhcieGwjsV+b48uOTsN7RQc9x8VOl2UivsrZboi1Pr89sSAQBbTEwBn/R7hSMJXXKA6TBrIEY9rtLlTTYhANuSCXSwCkx2kXuG2IV0UU3Ba1bf3gG+79vQhOwAaG2sZbR/i8L2nUwrHkHrvePSrG8JioOt8hkORleMMuSa; NIKE_CART=b:c; siteCatalyst_sample=67; dreamcatcher_sample=66; neo_sample=1; RT=\"sl=2&ss=1576688422033&tt=1629&obo=1&bcn=%2F%2F5f651e6f.akstat.io%2F&sh=1576688427824%3D2%3A1%3A1629%2C1576688425131%3D1%3A1%3A0&dm=nike.com&si=fdb6423a-1e01-4f9b-8e9d-89a200c8beca&ld=1576688427825\"; s_pers=%20c6%3Dproductgrid%253Acic%253Asnkrs%7C1576690228198%3B%20c5%3Dnikecom%253Esnkrs%253Ediscover%2520feed%253Efeed%7C1576690228203%3B; s_cc=true; ppd=productgrid:cic:snkrs|nikecom>snkrs>discover feed>feed; guidS=07a3f0c1-5994-4db2-c14d-fbb2e48c17b2; guidSTimestamp=1576688428229|1576688428229; guidU=b187dfec-740d-4eb3-e00b-ac6ae7cd9320; bm_sv=5FB576BBBEC083292DCC7200B95F10EC~VaxrELuEq5MVwoUY1G5PM6wacn1tZxIBrY1eZCdIwqlx/0t6mvy3nvcThtnP9tH9jhxEI/Lw3HTRz8VLLc/ezWuefJGSaxY30YMry5ipZ/8lFSQ1rrFpY3W/g18yHJbI9w1gxrHoC7WY9+uZMLAmzQ==; _abck=930BB1DBC221AD3F83C2E1AB29FAA80B~-1~YAAQLbV7XG1wCQFvAQAArTj0GQMK88NuxyiwxX7xuihNi5cesqOhTWmCAu6ab0U+ZZd5VdNFAEGOlQVuWMTO3/bPXCB/i1TDf5hlOttRxUhvhzH3JnJl+2qSeaHKxseabonMvTuHLxlkrEpwlLaFaBmxFMHx2fEo48Wptc+FAx/rE6WsE4UOEg6fgy8s08gg8lI2o3GWQHiXSrqq1KnjZMjEZgcIz3m4m9LT17eV8IUIIgVx+bq8U8c5p4PhW7JAdZh64I9mvMvat3RTfwn9zdAC7biLs+/4Cod6NkWX6Mh/5yK6/Qa9JW0Ct/4puaEuFAKdgarCunGKKUxkrTBNj5RG~-1~||-1~-1"
	}
	url := "https://api.nike.com/snkrs/content/v1/?offset=0&orderBy=published" + locale
	cookies := requests.ConvertCookies(cookies_string)

	// currency not used
	_ = currency;

	// ########################################### START REQUEST
	request := gorequest.New()
	resp, bodyBytes, request_err := request.Proxy(proxies.GrabProxy()).Get(url).
	// resp, bodyBytes, request_err := request.Get(url).
	Set("User-Agent", requests.RandomUserAgent()).
	Set("connection", "keep-alive").
	Set("authority", "api.nike.com").
	Set("pragma", "no-cache").
	Set("cache-control", "no-cache").
	Set("dnt", "1").
	Set("upgrade-insecure-requests", "1").
	Set("sec-fetch-user", "?1").
	Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
	Set("sec-fetch-site", "none").
	Set("sec-fetch-mode", "navigate").
	Set("accept-encoding", "deflate, br").
	Set("accept-language", "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7").
	AddCookies(cookies).
	EndBytes()
	if request_err != nil {
		log.Println(request_err)
		return
	}

	// ########################################### HANDLE ERRORS
	if requests.ParseHTTPErrors(resp, bodyBytes, identifier, simple_url, true) {
		failed_connections++
		return
	}

	// ########################################### INITIAL CHECK
	initialChecked := false
	if utils.StringInSlice(simple_url, initialCheckedURLs) {
		log.Println(Green("[" + "SNKRS " + strings.ToUpper(region) + "] " + "Successful connection."))
		initialChecked = true
	} else {
		initialCheckedURLs = append(initialCheckedURLs, simple_url)
		log.Println(Inverse("[" + "SNKRS " + strings.ToUpper(region) + "] " + "Initial Check Done."))
	}
	successful_connections++

	// ########################################### HANDLE RESPONSE
	data := &structs.SNKRSProduct{}
	err := json.Unmarshal(bodyBytes, data)
	if err != nil {
		log.Println(err)
		return
	}

	// fmt.Printf("%+v\n", data.Products[0]) // print first product

	var products structs.Products
	for i := 0;  i < len(data.Products); i++ {
		var product structs.Product
		// == important ==
		product.Store = simple_url
		product.StoreName = "SNKRS " + strings.ToUpper(region)

		// == main info ==
		product.Name = data.Products[i].Name
		product.URL = "https://" + simple_url + "/t/" + data.Products[i].Handle
		// product.URL = "https://" + "nike.com/launch" + "/t/" + data.Products[i].Handle
		if data.Products[i].Product.Price.Retail != -1 {
			// product.Price = strconv.Itoa(data.Products[i].Product.Price.Retail) + " " + currency
			product.Price = currencySymbol + fmt.Sprintf("%.2f", data.Products[i].Product.Price.Retail)
		} else {
			product.Price = ""
		}
		if data.Products[i].ImageURL == "" {
			product.ImageURL = "https://i.imgur.com/fip3nw5.png";
		} else {
			product.ImageURL = data.Products[i].ImageURL
		}

		// fetch description
		foundDescription := false
		for j := 0;  j < len(data.Products[i].Cards); j++ {
			if data.Products[i].Cards[j].Description != "" {
				product.Description = data.Products[i].Cards[j].Description
				foundDescription = true
				break
			}
		}
		if !foundDescription {
			product.Description = data.Products[i].SEODescription
		}
		// format description
		product.Description = strings.ReplaceAll(product.Description, "<p>", "")
		product.Description = strings.ReplaceAll(product.Description, "</p>", "")
		product.Description = strings.ReplaceAll(product.Description, "<br>", "")
		product.Description = strings.ReplaceAll(product.Description, "<strong>", "**")
		product.Description = strings.ReplaceAll(product.Description, "</strong>", "**")
		product.Description = strings.ReplaceAll(product.Description, "<em>", "_")
		product.Description = strings.ReplaceAll(product.Description, "</em>", "_")
		// shrink description to fit on a discord embed
		if len(product.Description) > 1024 {
			product.Description = product.Description[0:1024-3] + "..."
		}

		for j := 0;  j < len(data.Products[i].Product.SKUS); j++ {
			variant := structs.Variant{
				data.Products[i].Product.SKUS[j].Name,
				"",
				data.Products[i].Product.SKUS[j].Available,
				product.Price,
				-420,
			}
			product.Variants = append(product.Variants, variant)
		}
		product.Available = data.Products[i].Product.Available
		product.Identifier = identifier

		// == extra ==
		product.LaunchDate = data.Products[i].Product.LaunchDate
		product.Tags = data.Products[i].Tags
		convertedVariants, _ := json.Marshal(product.Variants)
		product.MD5 = utils.GetMD5Hash(string(convertedVariants))
		products = append(products, product)
	}

	database.SendToDatabase(products, identifier, simple_url, "SNKRS " + strings.ToUpper(region), initialChecked, mongoClient, pusherClient)

}
