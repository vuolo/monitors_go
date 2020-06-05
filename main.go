package main

import (
  "go.mongodb.org/mongo-driver/mongo" // MongoDB
	"go.mongodb.org/mongo-driver/mongo/options" // MongoDB Options

  "github.com/pusher/pusher-http-go" // Pusher

  . "github.com/logrusorgru/aurora" // colors

  "sync"
  "log"
  "context"
  "time"
  // "math"

  // ## local shits ##
  "./src/tasks"
  "./src/structs"
)

var useSSHTunnel = true // if running local, always keep this on and be sure to enable ssh tunnel
var shopifyEnabled = false
var funkoEnabled = false
var cpfmEnabled = false
var snkrsEnabled = false
var supremeEnabled = true
var stockxEnabled = false
var twitterEnabled = false
var instagramEnabled = false
var socialPlusEnabled = false

// var useSSHTunnel = false
// var shopifyEnabled = true
// var funkoEnabled = true
// var cpfmEnabled = true
// var snkrsEnabled = true
// var supremeEnabled = false
// var stockxEnabled = false
// var twitterEnabled = false
// var instagramEnabled = false
// var socialPlusEnabled = true

var successful_connections = 0
var failed_connections = 0

var SNKRSRegions = []string{
  "us",
  "gb",
  "jp",
  "cn",
  "sg",
  "ca",
  "de",
  "it",
  "ru",
}

var ShopifyStores = []structs.Store{

  // ##################### MAIN STORES

  structs.Store{
    URL: "travis-scott-secure.myshopify.com",
  	Name: "TRAVIS SCOTT",
  	Currency: "$",
  },
  structs.Store{
    URL: "undefeated.com",
  	Name: "UNDEFEATED",
  	Currency: "$",
  },
  structs.Store{
    URL: "bdgastore.com",
  	Name: "Bodega",
  	Currency: "$",
  },
  structs.Store{
    URL: "us.bape.com",
  	Name: "BAPE US",
  	Currency: "$",
  },
  structs.Store{
    URL: "kith.com",
  	Name: "KITH",
  	Currency: "$",
  },
  structs.Store{
    URL: "shop-usa.palaceskateboards.com",
    Name: "Palace Skateboards USA",
    Currency: "$",
  },

  // ##################### DSM EFLASH STORES

  structs.Store{
    URL: "eflash-jp.doverstreetmarket.com",
    Name: "Dover Street Market JP",
    Currency: "¥",
  },
  structs.Store{
    URL: "eflash-sg.doverstreetmarket.com",
    Name: "Dover Street Market SG",
    Currency: "S$",
  },
  structs.Store{
    URL: "eflash-us.doverstreetmarket.com",
    Name: "Dover Street Market US",
    Currency: "$",
  },
  structs.Store{
    URL: "eflash.doverstreetmarket.com",
    Name: "Dover Street Market London",
    Currency: "£",
  },

  // ##################### MISC STORES

  structs.Store{
    URL: "50-50.com.au",
    Name: "50-50 Skate Shop",
    Currency: "AU$",
  },
  structs.Store{
    URL: "www.6pmseason.com",
    Name: "6PM",
    Currency: "€",
  },
  structs.Store{
    URL: "a-cold-wall.com",
    Name: "A-COLD-WALL*",
    Currency: "£",
  },
  structs.Store{
    URL: "a-ma-maniere.com",
    Name: "A Ma Maniere",
    Currency: "$",
  },
  // structs.Store{
  //   URL: "abominabletoys.com",
  //   Name: "Abominable Toys",
  //   Currency: "$",
  // },
  // structs.Store{
  //   URL: "www.abovethecloudsstore.com", // robots.txt error
  //   Name: "Above The Clouds",
  //   Currency: "AU$",
  // },
  structs.Store{
    URL: "alifenewyork.com",
    Name: "Alife®",
    Currency: "$",
  },
  structs.Store{
    URL: "alumniofny.com",
    Name: "Alumni of NY",
    Currency: "$",
  },
  structs.Store{
    URL: "amigoskateshop.com",
    Name: "Amigos Skate Shop",
    Currency: "€",
  },
  structs.Store{
    URL: "antisocialsocialclub.com",
    Name: "AntiSocialSocialClub",
    Currency: "$",
  },
  structs.Store{
    URL: "anwarcarrots.com",
    Name: "Carrots By Anwar Carrots",
    Currency: "$",
  },
  structs.Store{
    URL: "casablancaparis.com",
    Name: "CasablancaParis",
    Currency: "€",
  },
  structs.Store{
    URL: "area51store.co.nz",
    Name: "Area 51",
    Currency: "NZ$",
  },
  structs.Store{
    URL: "bbbranded.com",
    Name: "BB Branded Boutique",
    Currency: "CA$",
  },
  structs.Store{
    URL: "bbcicecream.com",
    Name: "Billionaire Boys Club",
    Currency: "$",
  },
  structs.Store{
    URL: "blacksheepskateshop.com",
    Name: "Black Sheep Skate Shop",
    Currency: "$",
  },
  structs.Store{
    URL: "blendsus.com",
    Name: "BLENDS",
    Currency: "$",
  },
  structs.Store{
    URL: "bowsandarrowsberkeley.com",
    Name: "Bows and Arrows",
    Currency: "$",
  },
  structs.Store{
    URL: "bt21club.com",
    Name: "BT21CLUB",
    Currency: "$",
  },
  structs.Store{
    URL: "burnrubbersneakers.com",
    Name: "Burn Rubber",
    Currency: "$",
  },
  structs.Store{
    URL: "canada.finalmouse.com",
    Name: "Finalmouse CA",
    Currency: "$",
  },
  structs.Store{
    URL: "canteen.theberrics.com",
    Name: "The Berrics Canteen",
    Currency: "$",
  },
  structs.Store{
    URL: "cncpts.ae",
    Name: "Concepts DXB",
    Currency: "AE$",
  },
  structs.Store{
    URL: "cncpts.com",
    Name: "Concepts INTL",
    Currency: "$",
  },
  structs.Store{
    URL: "commonwealth-ftgg.com",
    Name: "Commonwealth",
    Currency: "$",
  },
  structs.Store{
    URL: "concrete.nl",
    Name: "Concrete",
    Currency: "€",
  },
  structs.Store{
    URL: "creme321.com",
    Name: "Creme321",
    Currency: "$",
  },
  structs.Store{
    URL: "crossoverconceptstore.com",
    Name: "CROSSOVER ONLINE",
    Currency: "RM",
  },
  structs.Store{
    URL: "crusoeandsons.com",
    Name: "Crusoeandsons",
    Currency: "$",
  },
  // structs.Store{
  //   URL: "culturekings.co.nz",
  //   Name: "Culture Kings NZ",
  //   Currency: "NZ$",
  // },
  structs.Store{
    URL: "culturekings.com",
    Name: "Culture Kings US",
    Currency: "$",
  },
  structs.Store{
    URL: "culturekings.com.au",
    Name: "Culture Kings AU",
    Currency: "AU$",
  },
  structs.Store{
    URL: "deadstock.ca",
    Name: "Deadstock.ca",
    Currency: "CA$",
  },
  structs.Store{
    URL: "defcongroup.com",
    Name: "DEFCON GROUP",
    Currency: "$",
  },
  structs.Store{
    URL: "dingdongtakuhaibin.com",
    Name: "Ding Dong Takuhaibin",
    Currency: "HK$",
  },
  structs.Store{
    URL: "dope-factory.com",
    Name: "Dope Factory",
    Currency: "€",
  },
  structs.Store{
    URL: "dropoutmilano.com",
    Name: "dropout",
    Currency: "€",
  },
  structs.Store{
    URL: "dropsau.com",
    Name: "DROPS.AU",
    Currency: "AU$",
  },
  structs.Store{
    URL: "dtlr.com",
    Name: "DTLR",
    Currency: "$",
  },
  structs.Store{
    URL: "empireskate.co.nz",
    Name: "Empire Skate",
    Currency: "NZ$",
  },
  structs.Store{
    URL: "eu.oneblockdown.it",
    Name: "One Block Down EU",
    Currency: "€",
  },
  structs.Store{
    URL: "everoboticsinc.com", // bot
    Name: "Eve Robotics",
    Currency: "$",
  },
  structs.Store{
    URL: "exoticpop.com",
    Name: "Exotic Pop",
    Currency: "$",
  },
  structs.Store{
    URL: "f3ather.io", // bot
    Name: "F3ATHER I/O",
    Currency: "$",
  },
  structs.Store{
    URL: "fearofgod.com",
    Name: "Fear of God",
    Currency: "$",
  },
  structs.Store{
    URL: "feature.com",
    Name: "Feature",
    Currency: "$",
  },
  structs.Store{
    URL: "ficegallery.com",
    Name: "ficegallery",
    Currency: "$",
  },
  structs.Store{
    URL: "finalmouse.com",
    Name: "Finalmouse",
    Currency: "$",
  },
  structs.Store{
    URL: "freshragsfl.com",
    Name: "Fresh Rags FL",
    Currency: "$",
  },
  structs.Store{
    URL: "fullsend.com",
    Name: "Full Send by Nelk Boys",
    Currency: "$",
  },
  // structs.Store{
  //   URL: "gymshark.com",
  //   Name: "Gymshark US",
  //   Currency: "$",
  // },
  structs.Store{
    URL: "hannibalstore.it",
    Name: "HannibalStore",
    Currency: "€",
  },
  structs.Store{
    URL: "hanon-shop.com",
    Name: "Hanon",
    Currency: "£",
  },
  structs.Store{
    URL: "highsandlows.net.au",
    Name: "Highs and Low",
    Currency: "AU$",
  },
  structs.Store{
    URL: "hoopsheaven.com.au",
    Name: "Hoops Heaven",
    Currency: "$",
  },
  structs.Store{
    URL: "humanmade.jp",
    Name: "HUMAN MADE ONLINE STORE",
    Currency: "¥",
  },
  structs.Store{
    URL: "it.oneblockdown.it",
    Name: "One Block Down IT",
    Currency: "€",
  },
  structs.Store{
    URL: "johngeigerco.com",
    Name: "John Geiger",
    Currency: "$",
  },
  structs.Store{
    URL: "juicestore.com",
    Name: "JUICESTORE",
    Currency: "HK$",
  },
  structs.Store{
    URL: "justdon.com",
    Name: "Just Don",
    Currency: "$",
  },
  structs.Store{
    URL: "kawsngv.com",
    Name: "KAWSNGV",
    Currency: "AU$",
  },
  structs.Store{
    URL: "kawsone.com",
    Name: "KAWSONE",
    Currency: "$",
  },
  structs.Store{
    URL: "laced.com.au",
    Name: "Laced",
    Currency: "AU$",
  },
  structs.Store{
    URL: "laceupnyc.com",
    Name: "Lace Up NYC",
    Currency: "$",
  },
  structs.Store{
    URL: "lapstoneandhammer.com",
    Name: "Lapstone & Hammer",
    Currency: "$",
  },
  structs.Store{
    URL: "lessoneseven.com",
    Name: "LESS 17",
    Currency: "CA$",
  },
  structs.Store{
    URL: "likelihood.us",
    Name: "LIKELIHOOD",
    Currency: "$",
  },
  structs.Store{
    URL: "limitededt.com",
    Name: "Limited Edt",
    Currency: "S$",
  },
  structs.Store{
    URL: "menuskateshop.com",
    Name: "Menu Skateboard Shop",
    Currency: "CA$",
  },
  structs.Store{
    URL: "mondotees.com",
    Name: "Mondo",
    Currency: "$",
  },
  structs.Store{
    URL: "mrcnoir.com",
    Name: "mrcnoir",
    Currency: "£",
  },
  structs.Store{
    URL: "nohble.com",
    Name: "Nohble",
    Currency: "$",
  },
  structs.Store{
    URL: "notre-shop.com",
    Name: "Notre",
    Currency: "$",
  },
  structs.Store{
    URL: "nrml.ca",
    Name: "NRML",
    Currency: "CA$",
  },
  structs.Store{
    URL: "offthehook.ca",
    Name: "Off The Hook",
    Currency: "CA$",
  },
  structs.Store{
    URL: "onenessboutique.com",
    Name: "Oneness Boutique",
    Currency: "$",
  },
  structs.Store{
    URL: "oqium.com",
    Name: "Oqium",
    Currency: "€",
  },
  structs.Store{
    URL: "packershoes.com",
    Name: "PACKER SHOES",
    Currency: "$",
  },
  structs.Store{
    URL: "pampamlondon.com",
    Name: "Pam Pam",
    Currency: "£",
  },
  structs.Store{
    URL: "par5-milano-yeezy.com",
    Name: "WALLSNEAKERS",
    Currency: "€",
  },
  structs.Store{
    URL: "pesoclo.com",
    Name: "PESOCLO",
    Currency: "€",
  },
  structs.Store{
    URL: "pleasuresnow.com",
    Name: "PLEASURES",
    Currency: "$",
  },
  structs.Store{
    URL: "primeonline.com.au",
    Name: "Prime",
    Currency: "AU$",
  },
  structs.Store{
    URL: "properlbc.com",
    Name: "Proper",
    Currency: "$",
  },
  structs.Store{
    URL: "prosperlosangeles.com",
    Name: "prosper",
    Currency: "$",
  },
  structs.Store{
    URL: "purchase.spectrebots.com", // bot
    Name: "Spectre Bots",
    Currency: "$",
  },
  structs.Store{
    URL: "renarts.com",
    Name: "Renarts",
    Currency: "$",
  },
  structs.Store{
    URL: "rh-ude.com",
    Name: "R H U D E",
    Currency: "$",
  },
  structs.Store{
    URL: "rockcitykicks.com",
    Name: "Rock City Kicks",
    Currency: "$",
  },
  structs.Store{
    URL: "row.oneblockdown.it",
    Name: "One Block Down ROW",
    Currency: "€",
  },
  structs.Store{
    URL: "rsvpgallery.com",
    Name: "RSVP Gallery",
    Currency: "$",
  },
  structs.Store{
    URL: "saintalfred.com",
    Name: "Saint Alfred",
    Currency: "$",
  },
  // structs.Store{
  //   URL: "shanedawsonmerch.com",
  //   Name: "Shane Dawson Merch",
  //   Currency: "$",
  // },
  structs.Store{
    URL: "shoegallerymiami.com",
    Name: "Shoe Gallery Inc",
    Currency: "$",
  },
  structs.Store{
    URL: "shop.100thieves.com",
    Name: "100 Thieves",
    Currency: "$",
  },
  structs.Store{
    URL: "shop.atlasskateboarding.com",
    Name: "Atlas",
    Currency: "$",
  },
  structs.Store{
    URL: "shop.balkobot.com", // bot
    Name: "balkobot",
    Currency: "$",
  },
  structs.Store{
    URL: "shop.destroyerbots.com", // bot
    Name: "Destroyer Shop",
    Currency: "$",
  },
  structs.Store{
    URL: "shop.exclucitylife.com",
    Name: "EXCLUCITYLIFE",
    Currency: "CA$",
  },
  structs.Store{
    URL: "shop.extrabutterny.com",
    Name: "Extra Butter New York",
    Currency: "$",
  },
  structs.Store{
    URL: "shop.ghostaio.com", // bot
    Name: "Ghost AIO",
    Currency: "$",
  },
  structs.Store{
    URL: "shop.goodasgoldshop.com",
    Name: "Good As Gold",
    Currency: "NZ$",
  },
  structs.Store{
    URL: "shop.havenshop.com",
    Name: "HAVEN",
    Currency: "CA$",
  },
  structs.Store{
    URL: "shop.marathonsports.com",
    Name: "Marathon Sports",
    Currency: "$",
  },
  structs.Store{
    URL: "shop.reigningchamp.com",
    Name: "Reigning Champ US",
    Currency: "$",
  },
  structs.Store{
    URL: "shopnicekicks.com",
    Name: "ShopNiceKicks.com",
    Currency: "$",
  },
  structs.Store{
    URL: "shopvlone.com",
    Name: "VLONE",
    Currency: "$",
  },
  structs.Store{
    URL: "sneakerpolitics.com",
    Name: "Sneaker Politics",
    Currency: "$",
  },
  structs.Store{
    URL: "sneakerworldshop.com",
    Name: "Sneakerworld",
    Currency: "€",
  },
  structs.Store{
    URL: "socialstatuspgh.com",
    Name: "Social Status",
    Currency: "$",
  },
  structs.Store{
    URL: "soleaio.com", // bot
    Name: "Sole AIO",
    Currency: "£",
  },
  structs.Store{
    URL: "soleclassics.com",
    Name: "Sole Classics",
    Currency: "$",
  },
  structs.Store{
    URL: "solefiness.com",
    Name: "SOLE FINESS",
    Currency: "AU$",
  },
  structs.Store{
    URL: "solefly.com",
    Name: "SoleFly",
    Currency: "$",
  },
  structs.Store{
    URL: "solestop.com",
    Name: "Solestop",
    Currency: "CA$",
  },
  structs.Store{
    URL: "soulland.com",
    Name: "Soulland",
    Currency: "£",
  },
  structs.Store{
    URL: "srarchive.myshopify.com",
    Name: "Stray Rats Archive",
    Currency: "$",
  },
  structs.Store{
    URL: "stashedsf.com",
    Name: "STASHED",
    Currency: "$",
  },
  structs.Store{
    URL: "stay-rooted.com",
    Name: "ROOTED",
    Currency: "$",
  },
  structs.Store{
    URL: "store.strayrats.com",
    Name: "Stray Rats",
    Currency: "$",
  },
  structs.Store{
    URL: "store.tomsachs.org",
    Name: "Tom Sachs Store",
    Currency: "$",
  },
  structs.Store{
    URL: "store.unionlosangeles.com",
    Name: "UNION LOS ANGELES",
    Currency: "$",
  },
  structs.Store{
    URL: "suede-store.com",
    Name: "SUEDE Store",
    Currency: "€",
  },
  structs.Store{
    URL: "thechinatownmarket.com",
    Name: "Chinatown Market",
    Currency: "$",
  },
  structs.Store{
    URL: "theclosetinc.com",
    Name: "The Closet Inc.",
    Currency: "CA$",
  },
  structs.Store{
    URL: "thedarksideinitiative.com",
    Name: "The Darkside Initiative",
    Currency: "$",
  },
  structs.Store{
    URL: "themobilebot.com",
    Name: "MBot",
    Currency: "$",
  },
  structs.Store{
    URL: "thenextdoor.fr",
    Name: "THE NEXT DOOR",
    Currency: "€",
  },
  structs.Store{
    URL: "thepremierstore.com",
    Name: "Premier",
    Currency: "$",
  },
  structs.Store{
    URL: "trainers-store.com.au",
    Name: "trainers",
    Currency: "AU$",
  },
  structs.Store{
    URL: "trophyroomstore.com",
    Name: "TROPHY ROOM STORE",
    Currency: "$",
  },
  structs.Store{
    URL: "unheardofbrand.com",
    Name: "Unheardof Brand",
    Currency: "$",
  },
  structs.Store{
    URL: "upsskateshop.com.au",
    Name: "U.P.S.",
    Currency: "AU$",
  },
  structs.Store{
    URL: "urbanindustry.co.uk",
    Name: "Urban Industry",
    Currency: "£",
  },
  structs.Store{
    URL: "us.octobersveryown.com",
    Name: "October's Very Own",
    Currency: "$",
  },
  structs.Store{
    URL: "usgstore.com.au",
    Name: "Urban Street Gear",
    Currency: "AU$",
  },
  structs.Store{
    URL: "warrenlotas.com",
    Name: "WARREN LOTAS",
    Currency: "$",
  },
  structs.Store{
    URL: "wearebraindead.com",
    Name: "Brain Dead",
    Currency: "$",
  },
  structs.Store{
    URL: "welcomeleeds.com",
    Name: "Welcome Skate Store",
    Currency: "£",
  },
  structs.Store{
    URL: "wishatl.com",
    Name: "Wish Atlanta",
    Currency: "$",
  },
  structs.Store{
    URL: "www.addictmiami.com",
    Name: "ADDICT Miami",
    Currency: "$",
  },
  structs.Store{
    URL: "www.aimeleondore.com",
    Name: "Aimé Leon Dore",
    Currency: "$",
  },
  structs.Store{
    URL: "www.amongstfew.com",
    Name: "amongst few",
    Currency: "AE$",
  },
  structs.Store{
    URL: "www.apbstore.com",
    Name: "APB Store",
    Currency: "$",
  },
  structs.Store{
    URL: "www.capsuletoronto.com",
    Name: "Capsule Online",
    Currency: "CA$",
  },
  structs.Store{
    URL: "www.centretx.com",
    Name: "Centre",
    Currency: "$",
  },
  structs.Store{
    URL: "www.cityblueshop.com",
    Name: "City Blue",
    Currency: "$",
  },
  structs.Store{
    URL: "www.diamondsupplyco.com",
    Name: "Diamond Supply Co.",
    Currency: "$",
  },
  structs.Store{
    URL: "www.flatspot.com",
    Name: "Flatspot",
    Currency: "£",
  },
  structs.Store{
    URL: "www.hombreofficial.com",
    Name: "Hombre Official",
    Currency: "€",
  },
  structs.Store{
    URL: "www.invincible.id",
    Name: "Invincible Jakarta Indonesia",
    Currency: "Rp ",
  },
  structs.Store{
    URL: "www.jimmyjazz.com",
    Name: "Jimmy Jazz",
    Currency: "$",
  },
  structs.Store{
    URL: "www.kongonline.co.uk",
    Name: "Kong Online",
    Currency: "£",
  },
  structs.Store{
    URL: "www.laces.mx",
    Name: "LACES ONLINE",
    Currency: "MEX$",
  },
  structs.Store{
    URL: "www.ldrs1354.com",
    Name: "Leaders 1354",
    Currency: "$",
  },
  structs.Store{
    URL: "www.machusonline.com",
    Name: "MACHUS",
    Currency: "$",
  },
  structs.Store{
    URL: "www.manorphx.com",
    Name: "Manor.",
    Currency: "$",
  },
  structs.Store{
    URL: "www.neighborhood.jp",
    Name: "NEIGHBORHOOD ONLINE STORE",
    Currency: "¥",
  },
  structs.Store{
    URL: "www.nocturnalskateshop.com",
    Name: "Nocturnal Skate Shop",
    Currency: "$",
  },
  structs.Store{
    URL: "www.noirfonce.eu",
    Name: "NOIRFONCE",
    Currency: "€",
  },
  structs.Store{
    URL: "www.oipolloi.com",
    Name: "Oi Polloi",
    Currency: "£",
  },
  structs.Store{
    URL: "www.patta.nl",
    Name: "Patta",
    Currency: "€",
  },
  structs.Store{
    URL: "www.philipbrownemenswear.co.uk",
    Name: "Philip Browne Menswear",
    Currency: "£",
  },
  structs.Store{
    URL: "www.rimenyc.com",
    Name: "RIME",
    Currency: "$",
  },
  structs.Store{
    URL: "www.rooneyshop.com",
    Name: "Rooney",
    Currency: "$",
  },
  structs.Store{
    URL: "www.slamcity.com",
    Name: "Slam City Skates",
    Currency: "£",
  },
  structs.Store{
    URL: "www.space23.it",
    Name: "Space23",
    Currency: "€",
  },
  structs.Store{
    URL: "www.staplepigeon.com",
    Name: "Staple Pigeon",
    Currency: "$",
  },
  structs.Store{
    URL: "www.stoneisland.co.uk",
    Name: "Stone Island UK",
    Currency: "£",
  },
  structs.Store{
    URL: "www.streetart.fr",
    Name: "STREETART.FR",
    Currency: "€",
  },
  structs.Store{
    URL: "www.stussy.com",
    Name: "Stussy US",
    Currency: "$",
  },
  structs.Store{
    URL: "www.stussy.co.uk",
    Name: "Stussy UK",
    Currency: "€",
  },
  structs.Store{
    URL: "www.sukamii.com",
    Name: "Sukamii",
    Currency: "$",
  },
  structs.Store{
    URL: "www.unknwn.com",
    Name: "UNKNWN",
    Currency: "$",
  },
  structs.Store{
    URL: "www.victorbraunstudios.com",
    Name: "Victor Braun Studios",
    Currency: "€",
  },
  structs.Store{
    URL: "www.westnyc.com",
    Name: "West NYC",
    Currency: "$",
  },
  structs.Store{
    URL: "xhibition.co",
    Name: "XHIBITION",
    Currency: "$",
  },
}

var FunkoStores = []structs.Store{

  structs.Store{
    URL: "shop.funko.com",
    Name: "Funko Shop",
    Currency: "$",
  },
  structs.Store{
    URL: "bigpopshop.com",
    Name: "Big Pop Shop",
    Currency: "$",
  },
  // structs.Store{
  //   URL: "bungiestore.com",
  //   Name: "Bungie Store",
  //   Currency: "$",
  // },
  structs.Store{
    URL: "galactictoys.com",
    Name: "Galactic Toys & Collectibles",
    Currency: "$",
  },
  structs.Store{
    URL: "www.fugitivetoys.com",
    Name: "Fugitive Toys",
    Currency: "$",
  },

}

var StockXHandles = []string{
  "supreme-the-north-face-paper-print-nuptse-jacket-paper-print",
  "nike-air-force-1-low-travis-scott-cactus-jack",
  "adidas-eqt-support-mid-adv-primeknit-dragon-ball-z-super-shenron",
  "adidas-human-race-nmd-pharrell-oreo",
  "adidas-nmd-hu-pharrell-inspiration-pack-powder-blue",
  "adidas-nmd-hu-pharrell-inspiration-pack-white",
  "supreme-x-louis-vuitton-monogram-scarf-brown",
  "adidas-nmd-hu-pharrell-solar-pack-orange",
  "adidas-nmd-hu-pharrell-mother",
  "supreme-x-louis-vuitton-monogram-bandana-red",
  "supreme-x-louis-vuitton-downtown-sunglasses-white",
  "supreme-world-famous-taped-seam-hooded-pullover-pullover-red",
  "adidas-nmd-hu-pharrell-x-billionaire-boys-club-multi-color",
  "adidas-ultra-boost-1-undefeated-white-red",
  "supreme-windstopper-zip-up-hooded-sweatshirt-fw19-bright-yellow",
  "supreme-undercoverpublic-enemy-work-jacket-dusty-teal",
  "supreme-tivoli-pal-bt-speaker-red",
  "supreme-the-north-face-statue-of-liberty-tee-white",
  "supreme-the-north-face-statue-of-liberty-tee-white",
  "supreme-the-north-face-statue-of-liberty-hoooded-sweatshirt-black",
  "supreme-the-north-face-snakeskin-taped-seam-coaches-jacket-green",
  "supreme-the-north-face-paper-print-nuptse-jacket-paper-print",
  "adidas-ultra-boost-4-game-of-thrones-nights-watch",
  "supreme-the-north-face-metallic-mountain-bib-pants-gold",
  "adidas-ultra-boost-4-game-of-thrones-targaryen",
  "supreme-tape-stripe-ls-pique-top-black",
  "adidas-ultra-boost-4-grey",
  "supreme-stripe-velour-l-s-polo-black",
  "adidas-ultra-boost-og-2018",
  "adidas-yeezy-boost-350-turtledove",
  "adidas-yeezy-boost-350-v2-blue-tint",
  "adidas-yeezy-boost-350-v2-core-black-green",
  "adidas-yeezy-boost-350-v2-white-core-black-red",
  "supreme-soccer-polo-white",
  "supreme-s-s-pocket-tee-red",
  "supreme-ripple-tank-top-blue",
  "supreme-polartec-short-black",
  "supreme-polartec-half-zip-pullover-black",
  "supreme-polartec-half-zip-pullover-black",
  "supreme-polartec-crewneck-white",
  "supreme-playboy-rayon-s-s-shirt-white",
  "supreme-ny-hooded-sweatshirt-ash-grey",
  "supreme-nike-hooded-sport-jacket-green",
  "supreme-new-shit-tee-black",
  "adidas-yeezy-boost-700-analog",
  "adidas-yeezy-boost-700-inertia",
  "adidas-yeezy-boost-700-salt",
  "adidas-yeezy-boost-700-v2-geode",
  "adidas-yeezy-boost-700-v2-tephra",
  "adidas-yeezy-boost-700-v2-vanta",
  "adidas-yeezy-powerphase-calabasas-core-black",
  "adidas-yeezy-wave-runner-700-solid-grey",
  "air-jordan-1-low-gold-toe",
  "air-jordan-1-mid-pine-green",
  "supreme-new-era-hq-beanie-red",
  "air-jordan-1-mid-shattered-backboard",
  "air-jordan-1-mid-top-3",
  "supreme-new-era-championship-beanie-royal",
  "supreme-mike-kelley-supreme-the-empire-state-building-tee-black",
  "supreme-life-tee-white",
  "supreme-lacoste-puffy-half-zip-pullover-red",
  "supreme-lacoste-pique-knit-camp-cap-lt-pink",
  "supreme-lacoste-logo-panel-sweatshort-red",
  "air-jordan-1-retro-high-bred-toe",
  "air-jordan-1-retro-high-court-purple",
  "supreme-lacoste-logo-panel-hooded-sweatshirt-light-blue",
  "supreme-know-thyself-hooded-sweatshirt-navy",
  "supreme-jean-paul-gaultier-fuck-racism-trucker-jacket-red",
  "supreme-jean-paul-gaultier-fuck-racism-trucker-jacket-gold",
  "supreme-jean-paul-gaultier-floral-print-hooded-sweatshirt-cardinal",
  "supreme-hq-waffle-thermal-black",
  "supreme-honda-fox-racing-v1-helmet-moss",
  "supreme-honda-fox-racing-moto-pant-red",
  "supreme-hanes-thermal-crew-1-pack-fw19-woodland-camo",
  "air-jordan-1-retro-high-homage-to-home-nn",
  "air-jordan-1-retro-high-neutral-grey-hyper-crimson",
  "air-jordan-1-retro-high-off-white-university-blue",
  "air-jordan-1-retro-high-phantom-gym-red",
  "air-jordan-1-retro-high-pine-greenb",
  "supreme-guts-tee-pink",
  "supreme-gore-tex-700-fill-down-parka-brown",
  "supreme-gilbert-george-supreme-life-tee-white",
  "air-jordan-1-retro-high-rookie-of-the-year",
  "air-jordan-1-retro-high-sports-illustrated",
  "air-jordan-1-retro-high-unc-leather",
  "air-jordan-12-retro-winter-black",
  "air-jordan-3-retro-katrina",
  "supreme-faces-l-s-tee-navy",
  "air-jordan-4-retro-cool-grey-2019",
  "supreme-embossed-logo-hooded-sweatshirt-ss18-black",
  "supreme-cutouts-tee-navy",
  "supreme-cuff-logo-beanie-brown",
  "air-jordan-6-retro-reflective-infrared",
  "nike-air-fear-of-god-1-sa-black",
  "nike-air-fear-of-god-1-sail-black",
  "nike-air-force-1-low-off-white-university-blue",
  "nike-air-force-1-mid-supreme-nba-black",
  "supreme-crown-air-freshener-red",
  "supreme-crossover-hockey-jersey-black",
  "nike-blazer-mid-77-vintage-white-black",
  "supreme-creeper-tee-white",
  "nike-kyrie-5-spongebob-squidward",
  "nike-daybreak-undercover-blue-jay",
  "nike-react-element-55-black-aurora-green",
  "supreme-comme-des-garcons-shirt-box-logo-tee-white",
  "supreme-chenille-stripe-beanie-light-blue",
  "supreme-cheese-tee-white",
  "nike-react-element-87-blue-chill-solar-red",
  "nike-react-element-87-black-neptune-green",
  "nike-react-element-55-triple-black",
  "nike-react-element-87-dusty-peach",
  "nike-react-element-87-knicks",
  "nike-react-element-87-light-orewood-brown",
  "nike-zoom-fly-off-white-pink",
  "supreme-20th-anniversary-box-logo-tee-white",
  "supreme-champion-hooded-satin-varsity-jacket-navy",
  "supreme-butterfly-knife-keychain-red",
  "supreme-box-logo-hooded-sweatshirt-black",
  "supreme-700-fill-down-taped-seam-parka-black",
  "supreme-aguila-tee-red",
  "supreme-apple-hooded-sweatshirt-brown",
  "supreme-basket-weave-beanie-black",
  "supreme-bedroom-tee-pink",
}

var curTwitterApp = -1
var twitterApps = []structs.TwitterApp{
  structs.TwitterApp{ // 1 incognito_test
    ConsumerKey: "OG1vddbvaKXbbROcAz2uB0w2d",
    ConsumerSecret: "rXm2gGuTq0xbvBn3Si4dcKUMHRpsa3gyUpOkEfnbwny98319aH",
    AccessTokenKey: "1165301596508622848-aA8SEVFvAd59LaQ3vO9W4ZAVgUsHpV",
    AccessTokenSecret: "UJYzGaxNgbGqX5MY9PBHvbnnVLi6pfdiFnUJxk9pUvGPn",
  },
  structs.TwitterApp{ // 2 incognito_test
    ConsumerKey: "hHFgalIj8uneWE15pw5zhJXex",
    ConsumerSecret: "w8IvhHBFpxHeROufIajy8DBrprj6kHp8vuxLLeXzLdzOMk02Sw",
    AccessTokenKey: "1165301596508622848-Av5h4IqmDtewRO6Bko3PT4wbuBDY37",
    AccessTokenSecret: "4wNKyh3IYe4ZajkkTKc1PDk5grTN9it8GSWJXj0Dwzok5",
  },
  structs.TwitterApp{ // 3 incognito_test
    ConsumerKey: "yfPbuVCKW3vTv1FBzT8J8hiEX",
    ConsumerSecret: "ZCncyKIBmVgjAQMnK4h6nAd72RI41R8m7yCoiNTNaQDGEzCcKp",
    AccessTokenKey: "1165301596508622848-97HpVZQOe2rS1AoPztzJX5YIQJBmE3",
    AccessTokenSecret: "qx2oh5p2YGgH3FA95dY5UejTV1g8FD3SoStJQDML3qDvY",
  },
  structs.TwitterApp{ // 4 incognito_test
    ConsumerKey: "jZQDeVJ51aenGSP8UFK5Txoz4",
    ConsumerSecret: "09M2GXWmVarACl5IOEA9m2tWD2NfnWXI9VwjgN2b9KGYMmYonR",
    AccessTokenKey: "1165301596508622848-PkLoHNuGCKqqqGjcvoM6gprQYdLKOe",
    AccessTokenSecret: "iOxQ6GtGSriIUNzCnydD0xNowMAj2zzdW2U0nrlfxxhfx",
  },
  structs.TwitterApp{ // 5 incognito_test
    ConsumerKey: "wFDgChB6LFTxNKt1GWw9DDpaM",
    ConsumerSecret: "7EAT2t3gIMhBIRSt7Nvk4k4I5ZJLnEascpc2qJCaQNeLj6tK8W",
    AccessTokenKey: "1165301596508622848-nJB8CmRPSovXpEHP8BQNNbMfRjQoPr",
    AccessTokenSecret: "UY0vFE7zTKG1OL2nwljOc2y1LeKElIrcJKGqttJyup1X7",
  },
  structs.TwitterApp{ // 6 incognito_test
    ConsumerKey: "DH6TfpqiqVLLzepQkbqeWCMbB",
    ConsumerSecret: "S0m6J2QPBnVp9FSob7vUZWF20yPDetMB6Wi5NidO4Tmd6ntueK",
    AccessTokenKey: "1165301596508622848-NqUwcRF0BUX7b134UVKJCc2dhji9mT",
    AccessTokenSecret: "o0bimy3ozYUwNvtr0OCLatKej22sMS71RhMgyP0r8bodZ",
  },
  structs.TwitterApp{ // 7 incognito_test
    ConsumerKey: "EjsTboznIFdsQuf8lUjPAwQde",
    ConsumerSecret: "Ucy14mnUcTx8IzqW9wVNDkevQn1MQDE9Y7X8f2NW0MHlTzmfSi",
    AccessTokenKey: "1165301596508622848-pa27OMMJdYDG0vSWhcZMVyNUoR0SRF",
    AccessTokenSecret: "PDHW8wi31LsuN0IwWV9CoWo7QU8P8D4wx6OfYykaYSWfg",
  },
  structs.TwitterApp{ // 8 incognito_test
    ConsumerKey: "eNY6GppacgcYZe1sJaLLS8VhR",
    ConsumerSecret: "OLRmeIrojsUwxSqxnG7zoAtXsuVsWOZkYpwgDKbH9pg6nfp6c8",
    AccessTokenKey: "1165301596508622848-goYyHNR8ND0wZQBzxZg35BNwasC7p5",
    AccessTokenSecret: "fq7NmqpgUwZdZdwm87y2YPdAACkBc47W8QeXT1wW6cqpL",
  },
  structs.TwitterApp{ // 9 incognito_test
    ConsumerKey: "Q5soe34eROjhSCLZxUSGqwN3S",
    ConsumerSecret: "Xk7glLmSaKpzL5BDVcQmsEfz3RPbwLKtfniZZNVFohz6drkfbn",
    AccessTokenKey: "1165301596508622848-AB9EJZlus9o3DtHyIQODXv0bSWVtHT",
    AccessTokenSecret: "GmM5LblJEhpaInMWAJ55Q1x7zw0oGxMT5gRleKKdhUdHW",
  },
  structs.TwitterApp{ // 10 incognito_test
    ConsumerKey: "pQVmOAbWCKpbkSdV6jnNEWDTt",
    ConsumerSecret: "MuxFs5qNlLjWiS2vbq8imDSgixjv5rz0XhT9UqYreoAXcDJhxp",
    AccessTokenKey: "1165301596508622848-SHqRAf2aPCPccWnXiOTlJkXh9HOqnA",
    AccessTokenSecret: "raDCdHo81XtHYCn6RbLTgo2oN32vMLgNuR64ArSvWOHzl",
  },
  structs.TwitterApp{ // 11 incognito_test2
    ConsumerKey: "YDYVcFaeqBaNuESmYI07qCFLq",
    ConsumerSecret: "5WKThL4QCnozZmE6ctsfjtEwWu8ZlsqrLiWCSpsApgcyWUpDwj",
    AccessTokenKey: "1211364372561321984-JIdGHuagzoR7xBs9kJKpkGiBSPBHHl",
    AccessTokenSecret: "MPezFeddWo9Br3kkRFCXihohJAUZss60uAsOWBQRqPYXS",
  },
  structs.TwitterApp{ // 12 incognito_test2
    ConsumerKey: "usud1TmuxHJguoFPFxuvBCIGK",
    ConsumerSecret: "2fdkMIszbm22Jq4ZQWVSIk6D2N5ryHFfXegWFduNPB2tO1n7Hm",
    AccessTokenKey: "1211364372561321984-3KzmYld1L4xbW0JSoHKtMoIPuZH2nZ",
    AccessTokenSecret: "VdvvwiRtgLnZQnYwt5kKJMR4bSCII2lIRJm31FHZ0cXgc",
  },
  structs.TwitterApp{ // 13 incognito_test2
    ConsumerKey: "UcV3SxdaHR0Yl3wfcH1lTPJGM",
    ConsumerSecret: "O2Nc0MYEHIChmOsvu4KnyfEBB8Uxe0cjqJPVcKZ3bixtTh5BES",
    AccessTokenKey: "1211364372561321984-abtWkqYIfjPCwElY31IhGrlcZxu2ua",
    AccessTokenSecret: "uf8qucNp9XTSkWeJNGoeH1s7RHaFAus7z45c6warLLe1R",
  },
  structs.TwitterApp{ // 14 incognito_test2
    ConsumerKey: "we0e1o25gJ9rlZvSUb69IiEAG",
    ConsumerSecret: "EnIWBWTbDC4FRSy2hVLHmnCr484hcp8UbFhBhzUSNlzksOAOTe",
    AccessTokenKey: "1211364372561321984-t2r7OPUkcpR5aqfaCqWiPey0m3AHPK",
    AccessTokenSecret: "nrXjYZ0Ie2sebgTaUZXTxSeKlyelw2at2eIkstY0DrczV",
  },
  structs.TwitterApp{ // 15 incognito_test2
    ConsumerKey: "AbKXQirGft2cL5pte7MFbNuOq",
    ConsumerSecret: "Eh9Lz3NjmqaZZx3WJ1MkJBRFM2TcZ4Sy7Iv4OccKgdDrSWqTD7",
    AccessTokenKey: "1211364372561321984-uilB5Lba3lm2B0SuwRT6DCyky7N2Nz",
    AccessTokenSecret: "02yPdW7PGnNPwUNn8G5cEZCiKjFxhgHNVmqZz83iZPkdG",
  },
  structs.TwitterApp{ // 16 incognito_test2
    ConsumerKey: "NEfLTda1OSEDjFPJKC1vzUOVd",
    ConsumerSecret: "gsyanJB2L26zDkOcNcDo0gsYhHJOEcqBtz78dy633wGghriE8V",
    AccessTokenKey: "1211364372561321984-yfBagmLsElol2YeD1GL6rsbB4M9IAI",
    AccessTokenSecret: "Z8GcUCQqq6NdoZBTFR8AjBJSsyhbK5HnvVo7i8E4l79N1",
  },
  structs.TwitterApp{ // 17 incognito_test2
    ConsumerKey: "UjaAx6ovNx1K33uUp4h5LFtLG",
    ConsumerSecret: "ZrQe45TyMMJkDOhZpKoBKikIFPyN1sBeudNGe9DfV0ypuRFWgX",
    AccessTokenKey: "1211364372561321984-v1uVVauLxgsurPRFtPVjIVQCvPhgTz",
    AccessTokenSecret: "HY8lZ60fY9Vmuieac3OPEEfMA0xwT3PY7nSVKxJbvUOU2",
  },
  structs.TwitterApp{ // 18 incognito_test2
    ConsumerKey: "q3tl3s4eL0TvZwb4REpyxbXUb",
    ConsumerSecret: "jtkW5BKix7jWdXVrdDtC4Blo6ap0bcEeFLkOcfMcdjNrEHUmqf",
    AccessTokenKey: "1211364372561321984-8AOP6wsZDH68ZH8n9HmzhNF44g0E6i",
    AccessTokenSecret: "wZVGR0mgzx4wCsXK3xmsBvzLp6rTGbQaZoGk7D6RjJPZn",
  },
  structs.TwitterApp{ // 19 incognito_test2
    ConsumerKey: "ABjdHcPH1EwXi1jaCuieJqsqM",
    ConsumerSecret: "McJLrDQoIPAkRaDlMzgtmLX5FLOeq5gngH0rz7NuQJSIoJuKiY",
    AccessTokenKey: "1211364372561321984-jkopj3Fn5NGo5vAEqCf2rJnZ0md5EO",
    AccessTokenSecret: "SynrRYzQ4dECYDzvyAT2bqwjmkOuywVJTWsuSq2lGdHx1",
  },
  structs.TwitterApp{ // 20 incognito_test3
    ConsumerKey: "74qayCLQRivvpL6OecS2JSMxS",
    ConsumerSecret: "W7Mcpz9MXkPgXIhmOoSVrk5YLma8nvPhIAq8f2LGzxKNCjBqEl",
    AccessTokenKey: "1211364372561321984-z8lAZ2TKyRPLufd2dVXMekXiR65Ia3",
    AccessTokenSecret: "8BfurCxJHOmAeJWjQ7nDckhlBKQp32ElqvNZtzwcOI1wK",
  },
  structs.TwitterApp{ // 21 incognito_test3
    ConsumerKey: "BKP9YOflwbJuQWt4CTuPAkGlN",
    ConsumerSecret: "CXF3GQVkfjucKhUNNTTPkfpSDag8C0YS8uYrpO6w4WL1TeIphS",
    AccessTokenKey: "1211368447654805505-ODWlsdPKf3ADxH49r1WLfx4wqMRGZZ",
    AccessTokenSecret: "UekyfT8TVzZSdoPi2AQ37uoqTZZU123TsEydoRBzO37Te",
  },
  structs.TwitterApp{ // 22 incognito_test3
    ConsumerKey: "azIVRZMVQATTgvfm1DuIjBMhF",
    ConsumerSecret: "QFZQgdsMw9z3TzgWWV5fsTxjLbphRKN0GUG7DTRiU7dces1CVQ",
    AccessTokenKey: "1211368447654805505-M4d4fda8PT4tafjeeOhPjkpsaHp5Pk",
    AccessTokenSecret: "LAVbDUjRoV2Qm4F3yGeSJtUflvytc17uPIoX9PEPE4Rpd",
  },
  structs.TwitterApp{ // 23 incognito_test3
    ConsumerKey: "OgvtbM4KLlp8mUaPAvIe6dvSv",
    ConsumerSecret: "QmCjMNCHNz9owNXqrG0IBZOiSvwGmuApv6eY2cX6PTxCTDtW2U",
    AccessTokenKey: "1211368447654805505-dyC4Qm8Rzt1e0DdPqVk3AnypppjAbg",
    AccessTokenSecret: "SdEiZpFL6vLBfTl4XgApKLLg9mmQLb3ECy1wiMkRdKidf",
  },
  structs.TwitterApp{ // 24 incognito_test3
    ConsumerKey: "OKaGpWNsZhT1JmKnCLZTApFj8",
    ConsumerSecret: "zwUPZFaFwoiHU7zuMKnUMn4dBCyo5NKI7msvIukDUhnj0Fo18U",
    AccessTokenKey: "1211368447654805505-fLaXoaT6QRO6JoM9HVNWSpmj1sfB4P",
    AccessTokenSecret: "TSo0pIvCpnunotZv8hs7t8h1Ks7LXdgtmeasnSUeRkFf2",
  },
  structs.TwitterApp{ // 25 incognito_test3
    ConsumerKey: "qRzurixWDNJZv158ba0ZpwrUc",
    ConsumerSecret: "58rU0699nJaCuZi165f2FY2WKgOkP0CeYluevuuMKpL13DEnU1",
    AccessTokenKey: "1211368447654805505-tLRyvS4OjJC0rDhQgYwXRRRW2kQRrj",
    AccessTokenSecret: "p68sX7wuox4inxUYVkFPfGOyeYLfK4kWxCwr0BUViSwKV",
  },
  structs.TwitterApp{ // 26 incognito_test3
    ConsumerKey: "zWAw7NvnS2z33k4D2a86gtcBL",
    ConsumerSecret: "BrIW0iMUGqUDKOII5fQVftEaeYCnKJjIWTDZD0ijvM6DPnhL1l",
    AccessTokenKey: "1211368447654805505-p5nETp2d3tYUxJiu3rScc4HVQvCT09",
    AccessTokenSecret: "qaHJqFBu5d3bgDcpSIGBpqyFHeEZe5PkNGSvNbCmOtVch",
  },
  structs.TwitterApp{ // 27 incognito_test3
    ConsumerKey: "mwtpKqk4fjZkRoFaz48ll1oz6",
    ConsumerSecret: "RkJgHVGLz0y0Ddr1eWBFdsNmbkcNfEc59KE7zINoUJIfeAjCCp",
    AccessTokenKey: "1211368447654805505-jsLRo32pitXNIoaMT4syUQGpq17hn4",
    AccessTokenSecret: "dfTv9XHnr19VglBs7WjeuFLKMLq1wel4m1sJxKJhErhyu",
  },
  structs.TwitterApp{ // 28 incognito_test3
    ConsumerKey: "lzAqTKtSkSFhqGv13NbUZ3Btc",
    ConsumerSecret: "NVOzWq1l5w7DZiKa5FDr4chaUashmDHABm4jE1gB1HmfN0yFSs",
    AccessTokenKey: "1211368447654805505-X7rnhJ1EWSvmFlhS8fwQ4IPAZJ34Fd",
    AccessTokenSecret: "Ye8twPTSl06qlImI36sk6O9uAOkYCXMQGqRxsBSSzduel",
  },
  structs.TwitterApp{ // 29 incognito_test3
    ConsumerKey: "q5f0hjPKInBqCmMFxTJ5ThC4I",
    ConsumerSecret: "QLZFkb2yO8imVC11nIPlAqxO1EodoIPcLSOTeBqIUBupTpBSuf",
    AccessTokenKey: "1211368447654805505-iU7oZqhLDLPrCly5yO52H0wXDznJuZ",
    AccessTokenSecret: "heV9AkgpFj3VB1pv5kkBRq8t3zfGBHjXJCjsWTUzP2vNI",
  },
  structs.TwitterApp{ // 30 incognito_test3
    ConsumerKey: "yhSx8tAsoBIvbnpCSSNtFOhqw",
    ConsumerSecret: "YNUgRJ3NFI8u0J4y15qad2f6GaKsfUjvEUMiHt6G51lfl2Rsvj",
    AccessTokenKey: "1211368447654805505-EBm2eTgqSSuBt5t6ddaVePBpzjp8dU",
    AccessTokenSecret: "1tfBwvkaAJPIjYL2bCCdq9qZGXaCAu0AOmtYjLw3nxPyV",
  },
}

var twitterHandles = []string{

  // ##################### BOTS

  "adeptbots",
  "backdoor",
  "balkobot",
  "cybersole",
  "dreamaio_",
  "eve_robotics",
  "offline", // cybersole owner
  "dashe",
  "f3ather",
  "fleekframework",
  "ganeshbot",
  "ghostaio",
  "hasteyio",
  "hawkmesh",
  "kodaiaio",
  "nova_aio",
  "mekrobotics",
  "prismaio",
  "destroyerbots",
  "sneakercopter",
  "soleaio",
  "splashforcebot",
  "swftaio",
  "thekickstation",
  "veloxpreme",
  "wrathbots",
  "FattyeXtension",

  // ##################### GROUPS

  "351io",
  "amnofify",
  "bouncealerts",
  "calicosio",
  "flipsio",
  "guap",
  "meknotify",
  "peachypings",
  "restockworld",
  "impersonated", // restock world owner
  // "ryxnszn", // restock world owner
  "thesitesupply",
  "strikeshoeshq",
  "saucemonitor",
}

var instagramHandles = []string{

  // ##################### MAIN ACCOUNTS
  // "dashpings_monitor",
  // "resellmonster_",

  // ##################### MISC ACCOUNTS

  "amnotifyus",
  "polarisaio",
  "fleekframework",
  "f3ather.io",
  "adeptbots",
  "kodaiaio",
  "destroyerbots",
  "cybersole",
  "offspringhq",
  "cncpts",
  "sacaiofficial",
  "ruggaio",
  "samuraibots",
  "kith",
  "dreamaio_",
  "balkobot",
  "nova_aio",
  "restockworld",
  "off____white",
  "splashforceio",
}

// ########################################### SETUP MONITORS MONGODB
// Set mongoClient options
var mongoClientOptions = options.Client().ApplyURI(getMongoURI(useSSHTunnel))

// Connect to MongoDB
var mongoClient, mongoClient_err = mongo.Connect(context.TODO(), mongoClientOptions)

func getMongoURI(useSSHTunnel bool) string {
  var mongoURI = "mongodb://localhost:27017"
  if useSSHTunnel {
    mongoURI = "mongodb://localhost:9999"
  }
  return mongoURI
}

func main() {

  // ########################################### SETUP PUSHER
  // pusherClient := pusher.Client{
  //   AppID: "943740",
  //   Key: "c5e1355ab7540c8f291a",
  //   Secret: "23effd2b170f76fbc5c9",
  //
  //   // Host: "websocket.resell.monster:6001",
  //   // Secure: true,
  //   // Cluster: "mt1",
  //
  //   Cluster: "us2",
  //   Secure: false,
  // }

  pusherClient, pusher_err := pusher.ClientFromURL("https://c5e1355ab7540c8f291a:23effd2b170f76fbc5c9@websocket.resell.monster:6001/apps/943740");
  if pusher_err != nil {
    log.Fatal(pusher_err)
  }

  // ########################################### SETUP MONITORS MONGODB
  if mongoClient_err != nil {
    log.Fatal(mongoClient_err)
  }

  // Check the connection
  ping_err := mongoClient.Ping(context.TODO(), nil)

  if ping_err != nil {
    log.Fatal(ping_err)
  }

  log.Println(Green("Connected to MongoDB!"))

  // ########################################### START SCRAPES
  var wg sync.WaitGroup

  if shopifyEnabled {
    // ########################################### SHOPIFY
    for _, Store := range ShopifyStores {
      // Increment the wait group counter
      wg.Add(1)
      go func(store structs.Store) {
        // Decrement the counter when the go routine completes
        defer wg.Done()
        // Call the task
        tasks.Shopify(store, mongoClient, pusherClient)
        }(Store)
      }
  }
  if funkoEnabled {
    // ########################################### FUNKO
    for _, Store := range FunkoStores {
      // Increment the wait group counter
      wg.Add(1)
      go func(store structs.Store) {
        // Decrement the counter when the go routine completes
        defer wg.Done()
        // Call the task
        tasks.Shopify(store, mongoClient, pusherClient)
        }(Store)
      }
  }
  if cpfmEnabled {
    // ########################################### CPFM
    CPFMStore := structs.Store{
      URL: "cactusplantfleamarket.com",
      Name: "CACTUS PLANT FLEA MARKET",
      Currency: "$",
    }
    wg.Add(1)
    go func(store structs.Store) {
      // Decrement the counter when the go routine completes
      defer wg.Done()
      // Call the task
      tasks.CPFM(store, mongoClient, pusherClient)
      }(CPFMStore)
  }
  if snkrsEnabled {
    // ########################################### SNKRS
    for _, SNKRSRegion := range SNKRSRegions {
      // Increment the wait group counter
      wg.Add(1)
      go func(region string) {
        // Decrement the counter when the go routine completes
        defer wg.Done()
        // Call the task
        tasks.SNKRS(region, mongoClient, pusherClient)
        }(SNKRSRegion)
      }
  }
  if supremeEnabled {
    // ########################################### SUPREME
    SupremeStore := structs.Store{
      URL: "supremenewyork.com",
      Name: "Supreme US",
      Currency: "$",
    }
    wg.Add(1)
    go func(store structs.Store) {
      // Decrement the counter when the go routine completes
      defer wg.Done()
      // Call the task
      tasks.Supreme(store, mongoClient, pusherClient)
      }(SupremeStore)
  }
  if stockxEnabled {
    // ########################################### StockX
    for _, StockXHandle := range StockXHandles {
      // Increment the wait group counter
      wg.Add(1)
      go func(stockXHandle string) {
        // Decrement the counter when the go routine completes
        defer wg.Done()
        // Call the task
        tasks.StockX(stockXHandle, mongoClient, pusherClient)
        }(StockXHandle)
      }
  }
  // if twitterEnabled {
  //   // ########################################### Twitter
  //   for _, twitterHandle := range twitterHandles {
  //     if curTwitterApp + 1 >= len(twitterApps) {
  //       curTwitterApp = 0
  //     } else {
  //       curTwitterApp++
  //     }
  //     // log.Println(len(twitterApps), curTwitterApp, twitterHandle)
  //     // Increment the wait group counter
  //     wg.Add(1)
  //     go func(handle string, appIndex int) {
  //       // Decrement the counter when the go routine completes
  //       defer wg.Done()
  //       // Call the task
  //       // var timeoutMS = math.Abs(1 - (100000/86400) - 1) * 1000 * (len(twitterHandles)/len(twitterApps)) * time.Millisecond
  //       var timeoutMS = 2 * 864 * time.Millisecond
  //       tasks.Twitter(handle, twitterApps[appIndex], timeoutMS, mongoClient, pusherClient)
  //     }(twitterHandle, curTwitterApp)
  //   }
  // }
  if instagramEnabled {
    // ########################################### Instagram
    for _, instagramHandle := range instagramHandles {
      // Increment the wait group counter
      wg.Add(1)
      go func(instagramHandle string) {
        // Decrement the counter when the go routine completes
        defer wg.Done()
        // Call the task
        tasks.Instagram(instagramHandle, mongoClient, pusherClient)
      }(instagramHandle)
    }
  }
  if socialPlusEnabled {
    wg.Add(1)
    go func() {
      // Decrement the counter when the go routine completes
      defer wg.Done()
      // Call the task
      for {
        // TODO: validate handle is still active, if not, end function
        if false {
          break
        }
        tasks.SetupSocialPlusActiveHandles(mongoClient, pusherClient)
        // log.Println(database.DatabaseTwitterHandles)
        // log.Println(database.DatabaseInstagramHandles)
        time.Sleep(864 * time.Millisecond) // 864 ms (0.864 second delay)
      }
      }()
  }

  // ########################################### WAIT FOR SCRAPES (prevents go from stopping early)
  wg.Wait()

}
