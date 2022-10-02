package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	//"math"
	"os"
	"strconv"

	union "github.com/MoonSHRD/IKY-telegram-bot/artifacts"

	//passport "IKY-telegram-bot/artifacts/TGPassport"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Verify personal wallet")),
)
var nullAddress common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")

//to operate the bot, put a text file containing key for your bot acquired from telegram "botfather" to the same directory with this file
var tgApiKey, err = os.ReadFile(".secret")
var bot, error1 = tgbotapi.NewBotAPI(string(tgApiKey))

type user struct {
	tgid          int64
	tg_username   string
	dialog_status int64
}

//main database for dialogs, key (int64) is telegram user id
var userDatabase = make(map[int64]user) // consider to change in persistend data storage?

var msgTemplates = make(map[string]string)

var myenv map[string]string

// file with settings for enviroment
const envLoc = ".env"

func main() {

	_ = godotenv.Load()
	ctx := context.Background()
	pk := os.Getenv("PK") // load private key from env
	gateway := os.Getenv("GATEWAY_RINKEBY_WS")

	bot, err = tgbotapi.NewBotAPI(string(tgApiKey))
	if err != nil {
		log.Panic(err)
	}

	// Connecting to blockchain network
	//  client, err := ethclient.Dial(os.Getenv("GATEWAY"))	// for global env config
	client, err := ethclient.Dial(gateway) // load from local .env file
	if err != nil {
		log.Fatalf("could not connect to Ethereum gateway: %v\n", err)
	}
	defer client.Close()

	// setting up private key in proper format
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	// Creating an auth transactor
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(4))
	//auth2:= bind.NewKeyedTransactorWithChainID(privateKey,big.NewInt(4))
	//NewKeyedTransactorWithChainID

	accountAddress := common.HexToAddress("0xc905803BbC804fECDc36850281fEd6520A346AC5")
	balance, _ := client.BalanceAt(ctx, accountAddress, nil) //our balance
	fmt.Printf("Balance of the validator bot: %d\n", balance)

	// Setting up Union
	union, err := union.NewUnion(common.HexToAddress("0x9024cF0a889233Af1fd4afaF949d5aF8C633D7fc"), client)
	if err != nil {
		log.Fatalf("Failed to instantiate a Union contract: %v", err)
	}

	log.Printf("session with union initialized")

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//whenever bot gets a new message, check for user id in the database happens, if it's a new user, the entry in the database is created.
	for update := range updates {

		if update.Message != nil {
			if _, ok := userDatabase[update.Message.From.ID]; !ok {
				userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, update.Message.Chat.UserName, 0}
				isRegistered := checkDao(auth, union, update.Message.Chat.ID)
				if isRegistered {
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {

						updateDb.dialog_status = 1
						userDatabase[update.Message.From.ID] = updateDb
					}
				}

			} else {

			}
		}
	}
}

func checkDao(auth *bind.TransactOpts, pc *union.Union, tgid int64) bool {

	str := strconv.FormatInt(tgid, 10)
	registration, err := pc.DaoAddresses(nil, str)

	if err != nil {
		log.Println("Can't check dao")
		log.Print(err)
	}
	if registration == nullAddress {
		return false
	} else {
		return true
	}

}
