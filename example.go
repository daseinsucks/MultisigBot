package main

import (
	"log"

	//"math"
	"os"

	union "github.com/MoonSHRD/IKY-telegram-bot/artifacts"

	//passport "IKY-telegram-bot/artifacts/TGPassport"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var nullAddress common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")

//to operate the bot, put a text file containing key for your bot acquired from telegram "botfather" to the same directory with this file
var tgApiKey, err = os.ReadFile(".secret")
var bot, error1 = tgbotapi.NewBotAPI(string(tgApiKey))

//main database for dialogs, key (int64) is telegram user id
//var userDatabase = make(map[int64]user) // consider to change in persistend data storage?

var msgTemplates = make(map[string]string)

var myenv map[string]string

// file with settings for enviroment
const envLoc = ".env"

func main() {

	_ = godotenv.Load()
	//ctx := context.Background()
	//pk := os.Getenv("PK") // load private key from env
	gateway := os.Getenv("GATEWAY_RINKEBY_WS")

	// Connecting to blockchain network
	//  client, err := ethclient.Dial(os.Getenv("GATEWAY"))	// for global env config
	client, err := ethclient.Dial(gateway) // load from local .env file
	if err != nil {
		log.Fatalf("could not connect to Ethereum gateway: %v\n", err)
	}
	defer client.Close()

	union, err := union.NewUnionCaller(common.HexToAddress("0x9024cF0a889233Af1fd4afaF949d5aF8C633D7fc"), client)
	if err != nil {
		log.Fatalf("Failed to instantiate a Union contract: %v", err)
	}
	log.Printf("session with union initialized")

	//test := big.NewInt(17236213761)
	version, err := union.GetDaoAddressbyChatId(nil, 1238723645)
	if err != nil {
		log.Print(err)
	}

	log.Println(version)

}
